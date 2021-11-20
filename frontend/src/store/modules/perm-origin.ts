import Vue from 'vue';
import { ActionContext } from 'vuex';
import { getStoreAccessors } from 'vuex-typescript';

import StoreState from '@/store/storeState';

import { Response } from '@/types/HTTP';
import { PermOrigin, RoutePermOrigin } from '@/types/Security';

type PermOriginContext = ActionContext<PermOriginState, StoreState>;

const API = process.env.VUE_APP_SERVER;

export interface PermOriginState {
    permOrigin: RoutePermOrigin;
}

const state: PermOriginState = {
    permOrigin: {} 
}

const store = {
    namespaced: true,
    state,
    getters: {
        getPermOrigins(state: PermOriginState): RoutePermOrigin {
            return state.permOrigin;
        }
    },
    mutations: {
        setPermOrigin(state: PermOriginState, permOrigins: PermOrigin[]): void {
            if (permOrigins.length === 0) return;

            const { routeId } = permOrigins[0];

            Vue.set(state.permOrigin, routeId, permOrigins);
        },
        savePermOrigin(state: PermOriginState, permOrigin: PermOrigin): void {
            const { routeId } = permOrigin;
            const routeOrigins = state.permOrigin[routeId];

            if (routeId in state.permOrigin) {
                const index = routeOrigins.findIndex(origin => origin.id === permOrigin.id);

                if (index !== -1) {
                    Object.assign(routeOrigins[index], permOrigin)
                } else {
                    routeOrigins.push(permOrigin);
                }

                return
            }

            // create the permission origin for the route
            // since it doesn't exist
            Vue.set(state.permOrigin, routeId, [permOrigin]);
        },
        removePermOrigin(state: PermOriginState, permOrigin: PermOrigin): void {
            const { routeId } = permOrigin;
            const routeOrigins = state.permOrigin[routeId];

            const index = routeOrigins.findIndex(origin => origin.id === permOrigin.id);
            if (index === -1) return;

            routeOrigins.splice(index, 1);
        }
    },
    actions: {
        fetchPermOrigins(context:  PermOriginContext, routeId: number): Promise<PermOrigin[]> {
            return new Promise<PermOrigin[]>((resolve, reject) => {
                fetch(`${API}/perm_origin?route_id=${routeId}`, {
                    method: "GET"
                })
                .then((res) => res.json())
                .then((body: Response<PermOrigin[]>) => {
                    if (body.successful) context.commit('setPermOrigin', body.data);
                    resolve(body.data);
                })
            })
        },
        createPermOrigin(context:  PermOriginContext, permOrigin: PermOrigin): Promise<PermOrigin> {
            return new Promise<PermOrigin>((resolve, reject) => {
                fetch(`${API}/perm_origin`, {
                    method: "POST",
                    body: JSON.stringify(permOrigin)
                })
                .then((res) => res.json())
                .then((body: Response<PermOrigin>) => {
                    if (body.successful) context.commit('savePermOrigin', body.data);
                    resolve(body.data);
                })
            })
        },
        updatePermOrigin(context:  PermOriginContext, permOrigin: PermOrigin): Promise<void> {
            return new Promise<void>((resolve, reject) => {
                fetch(`${API}/perm_origin`, {
                    method: "PUT",
                    body: JSON.stringify(permOrigin)
                })
                .then((res) => res.json())
                .then((body: Response<PermOrigin[]>) => {
                    if (body.successful) context.commit('savePermOrigin', body.data);
                    resolve();
                })
            })
        },
        deletePermOrigin(context:  PermOriginContext, permOrigin: PermOrigin): Promise<PermOrigin[]> {
            return new Promise<PermOrigin[]>((resolve, reject) => {
                fetch(`${API}/perm_origin?id=${permOrigin.id}&route_id=${permOrigin.routeId}`, {
                    method: "DELETE"
                })
                .then((res) => res.json())
                .then((body: Response<PermOrigin[]>) => {
                    if (body.successful) context.commit('removePermOrigin', permOrigin);
                })
            })
        }
    }
}

const { read, dispatch } = getStoreAccessors<PermOriginState, StoreState>("permOrigin");
const { getters, actions } = store;

export const getPermOrigins = read(getters.getPermOrigins);

export const fetchPermOrigins = dispatch(actions.fetchPermOrigins);
export const createPermOrigin = dispatch(actions.createPermOrigin);
export const updatePermOrigin = dispatch(actions.updatePermOrigin);
export const deletePermOrigin = dispatch(actions.deletePermOrigin);

export default store;