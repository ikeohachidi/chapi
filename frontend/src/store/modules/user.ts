import { ActionContext } from 'vuex';
import { getStoreAccessors } from 'vuex-typescript';
import StoreState from '@/store/storeState';
import { Response } from '@/types/HTTP';
import User from '@/types/User';

const API = process.env.VUE_APP_SERVER;

export type UserContext = ActionContext<UserState, StoreState>

export type UserState = {
    user: User | null
}

const state: UserState = {
    user: null 
}

const userStore = {
    namespaced: true,
    state,
    getters: {
        authenticatedUser(state: UserState): User | null {
            return state.user;
        },
        isUserAuthenticated(state: UserState): boolean {
            if (!state.user) {
                return false;
            }

            return state.user.id !== 0;
        }
    },
    mutations: {
        setUser(state: UserState, user: User): void {
            state.user = user;
        },
        nullifyUser(state: UserState): void {
            state.user = null;
        }
    },
    actions: {
        fetchAuthUser(context: UserContext): Promise<User> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/auth/user`, {
                    credentials: 'include'
                })
                .then(response => response.json())
                .then((body: Response<User>) => {
                    context.commit('setUser', body.data) 
                    resolve(body.data)
                })
                .catch(error => reject(error))
            })
        },
        logoutUser(context: UserContext): Promise<void> {
            return new Promise<void>((resolve, reject) => {
                fetch(`${API}/auth/logout`, {
                    credentials: 'include',
                    mode: 'cors'
                })
                .then(response => response.json())
                .then(() => {
                    context.commit('nullifyUser')
                    resolve()
                })
                .catch(error => reject(error))
            })
        }
    }
}

const { dispatch, read } = getStoreAccessors<UserState, StoreState>('user');
const { actions, getters } = userStore;

export const authenticatedUser = read(getters.authenticatedUser);
export const isUserAuthenticated = read(getters.isUserAuthenticated);

export const fetchAuthUser = dispatch(actions.fetchAuthUser);
export const logoutUser = dispatch(actions.logoutUser);

export default userStore;