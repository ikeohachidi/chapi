import { ActionContext } from 'vuex';
import { getStoreAccessors } from "vuex-typescript";
import Header from '@/types/Header';
import StoreState from '@/store/storeState';
import { Response } from '@/types/HTTP';

type HeaderContext = ActionContext<HeaderState, StoreState>;

const API = process.env.VUE_APP_SERVER;

export type HeaderState = {
    headers: Header[];
}

const state: HeaderState = {
    headers: []
}

const header = {
    namespaced: true,
    state,
    getters: {
        getHeaders(state: HeaderState): Header[] {
            return state.headers;
        }
    },
    mutations: {
        addHeader(state: HeaderState, header: Header): void {
            state.headers.push(header);
        },
        updateHeader(state: HeaderState, update: Header): void {
            const index = state.headers.findIndex(header => header.id === update.id)

            if (index === -1) return;

            Object.assign(state.headers[index], update)
        },
        deleteHeader(state: HeaderState, headerId: number): void {
            const index = state.headers.findIndex(header => header.id === headerId);

            if (index !== -1) {
                state.headers.splice(index, 1)
            }
        },
    },
    actions: {
        fetchRouteHeaders(context: HeaderContext, routeId: number): Promise<Header[]> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/header?route=${routeId}`, {
                    credentials: 'include',
                    mode: 'cors'
                })
                .then(response => response.json())
                .then((body: Response<Header[]>) => {
                    body.data.forEach(header => {
                        context.commit('addHeader', {
                            ...header,
                            routeId
                        })
                    })

                    resolve(body.data)
                })
                .catch(error => reject(error))
            })
        },
        saveHeader(context: HeaderContext, header: Header): Promise<Header> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/header`, { 
                        method: 'POST',
                        credentials: 'include',
                        body: JSON.stringify(header)
                    })
                    .then((res) => res.json())
                    .then((body: Response<Header>) => {
                        context.commit('addHeader', body.data);

                        resolve(body.data);
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
        updateHeader(context: HeaderContext, header: Header): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/header?route=${header.routeId}`, { 
                        method: 'PUT',
                        credentials: 'include',
                        body: JSON.stringify(header)
                    })
                    .then((res) => res.json())
                    .then((body: Response<string>) => {
                        context.commit('updateHeader', {
                            ...header,
                            id: body.data
                        })
                        resolve()
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
        deleteHeader(context: HeaderContext, header: Header): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/header?id=${header.id}&route_id=${header.routeId}`, {
                        method: 'DELETE',
                        credentials: 'include'
                    })
                    .then((res) => res.json())
                    .then((body: Response<string>) => {
                        if (body.successful) {
                            context.commit('deleteHeader', header.id)
                        }
                        resolve()
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
    }
}

const { read, dispatch } = getStoreAccessors<HeaderState, StoreState>('header');

const {actions, getters} = header 

export const getHeaders = read(getters.getHeaders);

export const fetchRouteHeaders = dispatch(actions.fetchRouteHeaders);
export const saveHeader = dispatch(actions.saveHeader);
export const updateHeader = dispatch(actions.updateHeader);
export const deleteHeader = dispatch(actions.deleteHeader);

export default header;