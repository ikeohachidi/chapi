import { ActionContext } from 'vuex';
import { getStoreAccessors } from 'vuex-typescript';
import Route, { CreateRouteRequest } from '@/types/Route';
import StoreState from '@/store/storeState';
import { Response, TestResponse } from '@/types/HTTP';

type RouteContext = ActionContext<RouteState, StoreState>;

const API = process.env.VUE_APP_SERVER;
const API_URL = process.env.VUE_APP_SITE_URL;

export type RouteState = {
    routes: Route[];
}

const state: RouteState = {
    routes: []
}

const route = {
    namespaced: true,
    state,
    getters: {
        getRoutes(state: RouteState): Route[] {
            return state.routes;
        }
    },
    mutations: {
        addRoute(state: RouteState, route: Route): void {
            state.routes.push(route)
        },
        updateRoute(state: RouteState, update: Route): void {
            const index = state.routes.findIndex(route => route.id === update.id)
            
            if (index === -1) return;

            Object.assign(state.routes[index], update)
        },
        removeRoute(state: RouteState, routeId: number): void {
            const index = state.routes.findIndex(route => route.id === routeId)
            
            if (index === -1) return;

            state.routes.splice(index, 1);
        }
    },
    actions: {
        createRoute(context: RouteContext, route: CreateRouteRequest): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/route`, {
                        method: 'POST',
                        body: JSON.stringify(route),
                        credentials: 'include'
                    })
                    .then((res) => res.json())
                    .then((body: Response<Route>) => {
                        if (body.successful) {
                            context.commit('addRoute', body.data);
                            resolve()
                        } else {
                            reject()
                        }
                    })
                    .catch((error) => {
                        reject(error)
                    })
            }) 
        },
        updateRoute(context: RouteContext, requestObject: Route): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/route`, {
                        method: 'PUT',
                        body: JSON.stringify(requestObject),
                        credentials: 'include'
                    })
                    .then((res) => res.json())
                    .then((body: Response<void>) => {
                        if (body.successful) {
                            context.commit('updateRoute', requestObject)
                            resolve()
                        } else {
                            reject(body.data)
                        }
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
        deleteRoute(context: RouteContext, routeId: number): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/route?id=${routeId}`, {
                        method: 'DELETE',
                        credentials: 'include'
                    })
                    .then((res) => res.json())
                    .then((body: Response<void>) => {
                        if (body.successful) {
                            context.commit('removeRoute', routeId)
                            resolve()
                        } else {
                            reject(body.data)
                        }
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
        fetchProjectRoutes(context: RouteContext, projectId: number): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/route?project=${projectId}`, {
                    credentials: 'include'
                })
                .then((res) => res.json())
                .then((body: Response<Route[]>) => {
                    if (body.successful) {
                        body.data.forEach(route => {
                            context.commit('addRoute', {
                                ...route,
                                projectId
                            })
                        })
                    }
                    resolve()
                })
                .catch((error) => {
                    reject(error)
                })
            }) 
        },
        testRoute(context: RouteContext, serverURL: string): Promise<TestResponse> {
            return new Promise((resolve, reject) => {
                let isResponseOk = true;
                const testResponse: TestResponse = {
                    headers: null,
                    status: 0,
                    responseType: "json",
                    statusText: '',
                    body: '',
                }
                fetch(serverURL, {
                    method: 'GET',
                    mode: 'cors',
                })
                .then(response => {
                    if (!response.ok) {
                        isResponseOk = false;
                    }

                    const contentType = response.headers.get('content-type');
                    testResponse.headers = response.headers;
                    testResponse.status = response.status;
                    testResponse.statusText = response.statusText;

                    if (contentType) {
                        if (contentType === 'application/json') {
                            return response.json();
                        }
    
                        if (contentType.includes('image')) {
                            testResponse.responseType = 'image';
                            return response.blob();
                        }
                    }
                    return response.text();
                })
                .then(body => {
                    if (isResponseOk) {
                        if (testResponse.responseType === 'image') {
                            testResponse.body = URL.createObjectURL(body)
                        } else {
                            testResponse.body = body;
                        }
                        resolve(testResponse)
                    } else {
                        reject(testResponse)
                    }
                })
                .catch(error => { 
                    reject(error)
                })
            })
        },
    }
}

const { read, dispatch } = getStoreAccessors<RouteState, StoreState>('route');

/**
 * key value pair with projectId as key and array 
 * of routes as value
 */
export const getRoutes = read(route.getters.getRoutes);

export const createRoute = dispatch(route.actions.createRoute);
export const updateRoute = dispatch(route.actions.updateRoute);
export const deleteRoute = dispatch(route.actions.deleteRoute);
export const fetchProjectRoutes = dispatch(route.actions.fetchProjectRoutes);

export const testRoute = dispatch(route.actions.testRoute);

export default route;