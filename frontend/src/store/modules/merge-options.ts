import { ActionContext } from 'vuex';
import { getStoreAccessors } from 'vuex-typescript';
import { MergeOptions, RouteMergeOption } from '@/types/Security';

import StoreState from '@/store/storeState';

import { Response } from '@/types/HTTP';
import Vue from 'vue';

type MergeOptionsContext = ActionContext<MergeOptionsState, StoreState>;

const API = process.env.VUE_APP_SERVER;

type MergeOptionsState = {
    options: RouteMergeOption
};

const state: MergeOptionsState = {
    options: {} 
}

const store = {
    namespaced: true,
    state,
    getters: {
        getMergeOptions(state: MergeOptionsState): RouteMergeOption {
            return state.options;
        },
    },
    mutations: {
        setState(state: MergeOptionsState, payload: MergeOptions[]): void {
            payload.forEach(item => {
                Vue.set(state.options, item.routeId, item)
            })
        },
        updateMergeOption(state: MergeOptionsState, payload: MergeOptions): void {
            state.options[payload.routeId] = payload;
        },
    },
    actions: {
        fetchMergeOptions(context: MergeOptionsContext, routeId: number): Promise<MergeOptions> {
            return new Promise<MergeOptions>((resolve) => {
                fetch(`${API}/merge_options?route_id=${routeId}`, {
                    method: "GET"
                })
                .then((res) => res.json())
                .then((body: Response<MergeOptions>) => {
                    if (body.successful) context.commit('setState', [{ ...body.data, routeId }]);
                    resolve(body.data);
                })
            })
        },
        updateMergeOption(context: MergeOptionsContext, payload: MergeOptions): Promise<void> {
            return new Promise<void>((resolve) => {
                fetch(`${API}/merge_options?route_id=${payload.routeId}`, {
                    method: "PUT",
                    body: JSON.stringify(payload)
                })
                .then((res) => res.json())
                .then((body: Response<void>) => {
                    if (body.successful) context.commit('updateMergeOption', payload);
                    resolve();
                })
            })
        }
    }
}

const { read, dispatch } = getStoreAccessors<MergeOptionsState, StoreState>("mergeOptions");
const { getters, actions } = store;

export const getMergeOptions = read(getters.getMergeOptions);

export const fetchMergeOptions = dispatch(actions.fetchMergeOptions);
export const updateMergeOption = dispatch(actions.updateMergeOption);

export default store;