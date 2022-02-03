import Vue from 'vue';
import { ActionContext } from 'vuex';
import { getStoreAccessors } from 'vuex-typescript';
import { MergeOptions } from '@/types/Security';

import StoreState from '@/store/storeState';

import { Response } from '@/types/HTTP';

type MergeOptionsContext = ActionContext<MergeOptions, StoreState>;

const API = process.env.VUE_APP_SERVER;

type MergeOptionsState = MergeOptions;

interface MutationPayload {
    property: keyof MergeOptionsState;
    state: boolean;
}

const state: MergeOptionsState = {
    mergeHeader: false,
    mergeQuery: false,
    mergeBody: false
}

const store = {
    namespaced: true,
    state,
    getters: {
        getMergeOptions(state: MergeOptionsState): MergeOptionsState {
            return state;
        }
    },
    mutations: {
        setState(state: MergeOptionsState, payload: MergeOptions): void {
            state = payload;
        },
        updateMergeOption(state: MergeOptionsState, payload: MutationPayload): void {
            Vue.set(state, payload.property, payload.state);
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
                    if (body.successful) context.commit('setState', body.data);
                    resolve(body.data);
                })
            })
        },
        updateMergeOption(context: MergeOptionsContext, payload: MutationPayload & { routeId: number }): Promise<void> {
            return new Promise<void>((resolve) => {
                fetch(`${API}/merge_options?route_id=${payload.routeId}`, {
                    method: "PUT",
                    body: JSON.stringify({ 
                        [payload.property]: payload.state
                    })
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