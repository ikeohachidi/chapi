<template>
    <section class="py-4">
        <div class="flex flex-row justify-between w-11/12 max-w-6xl mx-auto">
            <span @click="navigate('/')" class="flex items-center">
                <img src="@/assets/logo.png" alt="Chapi logo">
            </span>

            <ul>
                <template v-if="isUserAuthenticated">
                    <li 
                        v-for="route in routes" 
                        :key="route.text" 
                        @click="navigate(route.link)" 
                        class="px-3 cursor-pointer"
                    >
                        {{ route.text }}
                    </li>
                    <li @click="logoutUser">
                        Sign out
                    </li>
                </template>
                <li v-else class="relative sign-in">
                    <span>Sign in</span>
                    <ul>
                        <li @click="gitubAuth">
                            <span>Github</span>
                        </li>
                    </ul>
                </li>
            </ul>
        </div>
    </section>
</template>

<script lang='ts'>
import {Vue, Component} from 'vue-property-decorator';

import { logoutUser, isUserAuthenticated } from '@/store/modules/user';

const API = process.env.VUE_APP_SERVER;

@Component
export default class Navbar extends Vue {
    private routes = [
        { link: '/dashboard', text: 'Dashboard' },
    ]

    get isUserAuthenticated(): boolean {
        return isUserAuthenticated(this.$store)
    }

    private logoutUser() {
        logoutUser(this.$store)
            .then(() => {
                this.navigate('/')
            })
    }

    private gitubAuth() {
        window.location.href = API + '/auth/github'
    }

    private navigate(path: string) {
        this.$router.push(path)
    }
}
</script>

<style lang="scss" scoped>
ul {
    display: flex;
    margin-bottom: 0;
    list-style-type: none;
}

.sign-in {
    @apply cursor-pointer;

    & ul {
        @apply absolute -right-1/4 flex flex-col border border-gray-100 rounded-md shadow-md invisible;

        li {
            @apply transition duration-300 ease-in-out hover:bg-gray-100;
            @apply whitespace-nowrap px-5 py-2;
        }
    }

    &:hover ul {
        @apply visible;
    }
}
</style>