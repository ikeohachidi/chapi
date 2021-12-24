<template>
    <div class="route-card rounded-md bg-white border border-gray-200 p-4 hover:shadow-sm duration-75 relative">
        <p class="uppercase text-gray-400 absolute top-0 left-0 bg-gray-300 py-1 px-3 rounded-br-sm">{{ route.method }}</p>
        <div class="flex items-center mt-8">
            <span class="w-7 inline-block text-gray-500">
                <i class="ri-link"></i>
            </span>
            <p class="text-gray-900">{{ routeProject.name }}.{{ siteURL }}{{ route.path }}</p>
        </div>
        <div class="flex items-center mt-3">
            <span class="w-7 inline-block text-gray-500">
                <i class="ri-focus-3-line"></i>
            </span>
            <p class="text-gray-900">{{ route.destination }}</p>
        </div>
        <p class="font-bold text-sm mt-4">DESCRIPTION</p>
        <p class="text-gray-900 w-full">
            <template v-if="route.description">{{ route.description }}</template>
            <template v-else>No description provided for this route.</template>
        </p>
    </div>
</template>

<script lang='ts'>
import {Vue, Component, Prop} from 'vue-property-decorator';

import Route from '@/types/Route';
import Project from '@/types/Project';
import { getProjectById } from '@/store/modules/project';

@Component
export default class RouteCard extends Vue {
    @Prop({ default: () => { return new Route }}) route!: Route;

    get siteURL(): string {
        return process.env.VUE_APP_SITE_URL
    }

    get routeProject(): Project {
        if (!this.route.projectId) {
            return new Project();
        }

        return getProjectById(this.$store)(this.route.projectId) as Project
    }
}
</script>

<style scoped>
.route-card p {
    @apply overflow-ellipsis overflow-x-hidden;
}
</style>