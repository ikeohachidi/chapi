import { Query } from "@/types/Route";
import { ActionContext } from 'vuex';
import { getStoreAccessors } from "vuex-typescript";
import StoreState from '@/store/storeState';
import { Response } from '@/types/HTTP';

type QueryContext = ActionContext<QueryState, StoreState>;

const API = process.env.VUE_APP_SERVER

export interface QueryState {
    queries: Query[]
}

const queriestate: QueryState = {
    queries: []
} 

const query = {
    namespaced: true,
    state: queriestate,
    getters: {
        getQueries(state: QueryState): Query[] {
            return state.queries
        }
    },
    mutations: {
        addQuery(state: QueryState, query: Query): void {
            state.queries.push(query);
        },
        updateQuery(state: QueryState, update: Query): void {
            const index = state.queries.findIndex(query => query.id === update.id)

            if (index === -1) return;

            Object.assign(state.queries[index], update)
        },
        deleteQuery(state: QueryState, queryId: number): void {
            const index = state.queries.findIndex(query => query.id === queryId);

            if (index !== -1) {
                state.queries.splice(index, 1)
            }
        },
    },
    actions: {
        fetchRouteQueries(context: QueryContext, routeId: number): Promise<Query[]> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/query?route=${routeId}`, {
                    credentials: 'include',
                    mode: 'cors'
                })
                .then(response => response.json())
                .then((body: Response<Query[]>) => {
                    body.data.forEach(query => {
                        context.commit('addQuery', {
                            ...query,
                            routeId
                        })
                    })

                    resolve(body.data)
                })
                .catch(error => reject(error))
            })
        },
        saveQuery(context: QueryContext, query: Query): Promise<Query> {
            return new Promise((resolve, reject) => {
                console.log(query)
                fetch(`${API}/query?route=${query.routeId}`, { 
                        method: 'POST',
                        credentials: 'include',
                        body: JSON.stringify(query)
                    })
                    .then((res) => res.json())
                    .then((body: Response<Query>) => {
                        context.commit('addQuery', body.data);

                        resolve(body.data)
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
        updateQuery(context: QueryContext, query: Query): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/query?route=${query.routeId}`, { 
                        method: 'PUT',
                        credentials: 'include',
                        body: JSON.stringify(query)
                    })
                    .then((res) => res.json())
                    .then((body: Response<string>) => {
                        context.commit('updateQuery', {
                            ...query,
                            id: body.data,
                        })
                        resolve()
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
        deleteQuery(context: QueryContext, query: Query): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/query?id=${query.id}&route_id=${query.routeId}`, {
                        method: 'DELETE',
                        credentials: 'include'
                    })
                    .then((res) => res.json())
                    .then((body: Response<string>) => {
                        if (body.successful) {
                            context.commit('deleteQuery', query.id)
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

const { read, dispatch } = getStoreAccessors<QueryState, StoreState>('query');

const {actions, getters} = query 

export const getQueries = read(getters.getQueries);

export const fetchRouteQueries = dispatch(actions.fetchRouteQueries);
export const saveQuery = dispatch(actions.saveQuery);
export const updateQuery = dispatch(actions.updateQuery);
export const deleteQuery = dispatch(actions.deleteQuery);

export default query;