<template>
    <section class="px-5 py-5 grid grid-cols-3 grid-rows-4 gap-5 min-h-full">
        <modal
            title="New Route"
            description="Please configure the base settings for this route"
            actionButtonText="Create Route"
            v-if="showNewRouteModal"
            @close="showNewRouteModal = false"
            @action="createNewRoute"
        >
            <div class="flex">
                <select v-model="method" class="rounded-r-none w-3/12" style="padding: 8px;">
                    <option v-for="option in HTTPMethodOptions" :key="option" :value="option">{{ option }}</option>
                </select>
                <p class="w-6/12 border mb-0 px-4 text-sm border-l-0 cursor-not-allowed flex items-center overflow-ellipsis border-gray-200" :title="serverURL">{{ serverURL }}</p>
                <input type="text" placeholder="/users" class="rounded-l-none border-l-0 w-5/12" v-model="path">
            </div>
            <small class="block text-right" v-if="correctedPath !== path">Your path will be: 
                <span class="text-green-600">{{ correctedPath }}</span>
            </small>
            <input type="text" placeholder="Destination URL" v-model="destination" class="w-full mt-5">
            <p class="text-sm error-text" v-if="validationErrors.destination">{{ validationErrors.destination }}</p>
            <textarea rows="8" class="mt-5 resize-none w-full" v-model="description" placeholder="Route Description"></textarea>
        </modal>
        <div class="w-full rounded-lg border-dashed border-4 flex items-center justify-center" @click="showNewRouteModal = true">
            <div class="transform scale-150 text-gray-400">
                <i class="ri-add-line"></i>
            </div>
        </div>
        <RouteCard 
            v-for="(route, routeIndex) in projectRoutes"
            :key="routeIndex"
            :route="route"
            @click.native="goToRoute(route)"
        />
    </section>
</template>

<script lang='ts'>
import {Vue, Component, Watch} from 'vue-property-decorator';

import RouteCard from '@/components/RouteCard/RouteCard.vue';
import Modal from '@/components/Modal/Modal.vue';

import Route, { CreateRouteRequest } from '@/types/Route';
import { HTTPMethod } from '@/types/HTTP';

import { createRoute, fetchProjectRoutes, getRoutes } from '@/store/modules/route';
import { getProjectById } from '@/store/modules/project';

type ErrorMessage = string;

type ValidatoinErrors = {
    errors: {[input: string]: ErrorMessage},
    errorsMessage: string;
    hasErrors: boolean;
}

@Component({
    components: {
        RouteCard,
        Modal
    }
})
export default class RouteList extends Vue {
    private path = '';
    get correctedPath(): string {
        return this.path.replace(/\s+/g, '-')
    }
    private description = '' ;
    private destination = '';
    private method: HTTPMethod = HTTPMethod.GET;
    private showNewRouteModal = false;

    get projectId(): number {
        return Number(this.$route.query['project']);
    }
    @Watch('projectId')
    onProjectIdChange(): void {
        if (this.projectRoutes.length === 0) this.fetchProjectRoutes();
    }

    get serverURL(): string {
        const project = getProjectById(this.$store)(this.projectId);

        return `${project ? project.name : ''}.${process.env.VUE_APP_SITE_URL}`;
    }

    get HTTPMethodOptions(): string[] {
        return Object.keys(HTTPMethod)
    }

    get projectRoutes(): Route[] {
        return getRoutes(this.$store).filter(route => route.projectId === this.projectId)
    }

    private resetInputs(): void {
        this.path = '';
        this.description = '';
        this.method = HTTPMethod.GET;
        this.destination = '';
    }

    private validationErrors(): ValidatoinErrors {
        const errors: {[input: string]: ErrorMessage} = {};
        if (this.destination === '') errors['destination'] = "A destination URL is required";

        const errorsMessage = Object.keys(errors).reduce((acc: string, errorInput) => {
            return acc + `${errors[errorInput]}\n`
        }, '')

        return {
            errors,
            errorsMessage,
            hasErrors: Object.keys(errors).length > 0
        };
    }

    private createNewRoute(): void {
        const { errorsMessage, hasErrors } = this.validationErrors();

        if (hasErrors) {
            this.$toast.error(errorsMessage)
            return
        }

        const route: CreateRouteRequest = {
            projectId: this.projectId,
            path: this.correctedPath.toLowerCase(),
            method: this.method,
            description: this.description,
            destination: this.destination.toLowerCase()
        }

        createRoute(this.$store, route)
            .then(() => {
                this.$toast.success('Route Created Successfully');
            })
            .catch(() => {
                this.$toast.error('Error creating route');
            })
        this.showNewRouteModal = false;
        this.resetInputs();
    }

    private goToRoute(route: Route): void {
        this.$router.push({
            name: 'Route' ,
            query: {
                project: String(this.projectId),
                route: route ? String(route.id) : ''
            }
        })
    }

    private fetchProjectRoutes(): void {
        fetchProjectRoutes(this.$store, this.projectId)
            .catch(error => {
                this.$toast.error(`Error getting project routes: ${error}`)
            })
    }

    mounted(): void {
        if (this.projectRoutes.length === 0) this.fetchProjectRoutes();
    }
}
</script>