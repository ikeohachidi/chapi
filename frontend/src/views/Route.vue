<template>
    <section>
        <config-test-overlay 
            v-if="showConfigResult" 
            :configResult="configTestResult"
            :isLoading="isLoadingTest"
            @redo-test="testRouteConfig"
            @close="showConfigResult = false"
        />

        <div class="content-padding bg-gray-100 flex items-center py-5">
            <p class="pb-0 font-mono text-sm">
                {{ serverURL }}
            </p>
            <button class="ml-auto" @click="testRouteConfig">Test</button>
        </div>

        <tabs position="center">
            <tab title="Request">
                <request :projectId="projectId" :route="route"/>
            </tab>
            <tab title="Metadata">
                <metadata :route="route"/>
            </tab>
            <tab title="Security">
                <security :route="route"/>
            </tab>
            <tab title="Analytics">
                <div class="section content-padding">
                    <p>Coming soon 😃</p>
                </div>
            </tab>
        </tabs>

        <div class="content-padding bg-gray-100 py-5">
            <h1 class="text-2xl text-red-600 font-bold">Delete Route</h1>
            <p class="mt-3">This operation is not reversible and will delete all the configuration set on this route</p>
            <button class="danger mt-5" @click="deleteRoute">Delete</button>
        </div>
    </section>
</template>

<script lang='ts'>
import {Vue, Component} from 'vue-property-decorator';
import { Route as VueRoute, NavigationGuardNext } from 'vue-router';

import ConfigTestOverlay from '@/components/ConfigTestOverlay/ConfigTestOverlay.vue';
import { Tab, Tabs } from '@/components/Tabs';
import { Metadata, Request, Security } from './route-tabs'

import Route from '@/types/Route'
import Project from '@/types/Project';
import { testRoute, getRoutes, deleteRoute, fetchProjectRoutes } from '@/store/modules/route';
import { getProjectById } from '@/store/modules/project';

@Component({
    components: {
        ConfigTestOverlay,
        Tab, Tabs,
        Metadata,
        Request,
        Security
    },
    beforeRouteEnter: (to: VueRoute, from: VueRoute, next: NavigationGuardNext) => {
        if (!to.query['project'] || !to.query['route']) {
            next({ path: 'Dashboard' })
        }

        if (from.name !== 'Routes List') {
            next({ name: 'Dashboard' })
        }

        next()
    }
})
export default class RouteView extends Vue {
    private showConfigResult = false;
    private configTestResult = {
        data: {},
        type: false,
        responseTime: 0,
    }
    private isLoadingTest = false;

    get route(): Route | null {
        const routes = getRoutes(this.$store);

        const index = routes.findIndex(route => route.id === this.routeIdQuery)

        if (index === -1) return null;

        return routes[index]
    }

    get projectId(): number {
        return Number(this.$route.query['project']);
    }

    get routeProject(): Project {
        return getProjectById(this.$store)(this.projectId) as Project;
    }

    get routeIdQuery(): number {
        return Number(this.$route.query['route']);
    }

    get siteURL(): string {
        return process.env.VUE_APP_SITE_URL
    }

    get serverURL(): string {
        if (this.route) {
            const scheme = process.env.NODE_ENV === 'development' ? 'http://' : 'https://';
            return `${scheme}${this.routeProject.name}.${this.siteURL}${this.route.path}`;
        }
        return '';
    }

    private deleteRoute() {
        if (!this.route) return;
        deleteRoute(this.$store, this.route.id)
            .then(() => {
                this.$toast.success('Route deleted successfully');
                this.$router.go(-1);
            })
            .catch(() => this.$toast.error('Error deleting route'))
    }

    private testRouteConfig() {
        this.showConfigResult = true;

        const start = new Date().getTime();

        this.isLoadingTest = true;
        testRoute(this.$store, this.serverURL)
            .then(response => {
                this.configTestResult = {
                    data: response,
                    type: true,
                    responseTime: (new Date().getTime() - start) / 1000
                }
            })
            .catch(error => {
                this.configTestResult = {
                    data: error,
                    type: false,
                    responseTime: (new Date().getTime() - start) / 1000
                }
            })
            .finally(() => {
                this.isLoadingTest = false;
            })
    }

    mounted(): void {
        if (this.route === null) {
            fetchProjectRoutes(this.$store, this.projectId)
        }
    }
}
</script>

<style scoped> 
section {
    @apply min-h-screen;
}
</style>