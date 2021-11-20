<template>
    <section>
        <div class="section content-padding">
            <div>
                <p class="section-name">Destination URL</p>
                <p class="section-description">
                    Please enter the URL you'd like to make a request to. Make sure you have the right 
                    <span class="inline-block rounded-sm bg-gray-300 p-1 font-mono">HTTP Method</span> set
                </p>
            </div>
            <div class="flex flex-col">
                <div>
                    <select v-model="routeUpdate.method" class="rounded-r-none w-2/12" style="padding: 8px;" @change="updateRequest">
                        <option v-for="method in HTTPMethodOptions" :key="method" :value="method" class="uppercase">{{ method }}</option>
                    </select>
                    <input class="rounded-l-none border-l-0 w-10/12" v-model="routeUpdate.destination">
                </div>
                <button 
                    class="mt-4 ml-auto"
                    :disabled="(routeUpdate.destination === route.destination) && (routeUpdate.method === route.method)"
                    @click="updateRequest"
                >
                    Save
                </button>
            </div>
        </div>
        <div class="section content-padding">
            <div>
                <div class="section-name">URL Queries</div>
                <div class="section-description">
                    Enter the queries that will be placed on the request URL 
                </div>
            </div>
            <div>
                <table class="w-full">
                    <thead>
                        <tr class="text-left text-gray-800">
                            <th class="py-2 pl-3 font-normal">Name</th>
                            <th class="px-3 font-normal">Value</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(query, queryIndex) in routeQueries" :key="queryIndex">
                            <td class="py-2 w-2/5">
                                <input class="w-full" type="text" v-model="query.name">
                            </td>
                            <td class="px-3 w-2/5">
                                <input class="w-full" type="text" v-model="query.value">
                            </td>
                            <td class="w-1/5">
                                <div class="flex items-center justify-end">
                                    <button @click="updateQuery(query)" :disabled="hasQueryChanged(queryIndex)" v-if="query.id > 0">
                                        Update
                                    </button>
                                    <button @click="saveQuery(query)" v-else>
                                        Save 
                                    </button>
                                    <span class="text-red-600 cursor-pointer ml-4 hover:scale-50" @click="deleteQuery(query)">
                                        <i class="ri-delete-bin-line"></i>
                                    </span>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <button class="mt-4 ml-auto" @click="addQuery">Add Query</button>
            </div>
        </div>
        <div class="section content-padding">
            <div>
                <div class="section-name">Headers</div>
                <div class="section-description">
                    <p>Enter HTTP headers that will be added to your request during execution.</p>
                    <p class="font-bold text-indigo-700">Eg. Bearer-Token: xx-xx-xx</p>
                </div>
            </div>
            <div>
                <table class="w-full">
                    <thead>
                        <tr class="text-left text-gray-800">
                            <th class="py-2 pl-3 font-normal">Name</th>
                            <th class="px-3 font-normal">Value</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(header, headerIndex) in routeHeaders" :key="headerIndex">
                            <td class="py-2 w-2/5">
                                <input class="w-full" type="text" v-model="header.name">
                            </td>
                            <td class="px-3 w-2/5">
                                <input class="w-full" type="text" v-model="header.value">
                            </td>
                            <td class="w-1/5">
                                <div class="flex items-center justify-end">
                                    <button @click="updateHeader(header)" :disabled="hasHeaderChanged(headerIndex)" v-if="header.id > 0">
                                        Update
                                    </button>
                                    <button @click="saveHeader(header)" v-else>
                                        Save 
                                    </button>
                                    <span class="text-red-600 cursor-pointer ml-4 hover:scale-50" @click="deleteHeader(header)">
                                        <i class="ri-delete-bin-line"></i>
                                    </span>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <button class="mt-4 ml-auto" @click="addHeader">Add Header</button>
            </div>
        </div>
        <div class="section content-padding">
            <div>
                <div class="section-name">Request Body</div>
                <div class="section-description">
                    You can add a request body to your route. Please note that the body should be in a JSON format.
                </div>
            </div>
            <div>
                <textarea 
                    rows="10" 
                    class="w-full font-mono resize-none" 
                    onkeydown="if(event.keyCode===9){var v=this.value,s=this.selectionStart,e=this.selectionEnd;this.value=v.substring(0, s)+'\t'+v.substring(e);this.selectionStart=this.selectionEnd=s+1;return false;}"
                    placeholder="{}"
                    v-model="routeUpdate.body"
                >
                </textarea>
                <div class="flex">
                    <button 
                        class="mt-4 ml-auto" 
                        @click="updateRequest"
                        :disabled="routeUpdate.body === route.body"
                    >
                        Save
                    </button>
                </div>
            </div>
        </div>
    </section>
</template>

<script lang='ts'>
import {Vue, Component, Prop, Watch} from 'vue-property-decorator';

import Route, { Query } from '@/types/Route';
import { HTTPMethod } from '@/types/HTTP';
import Header from '@/types/Header';

import { saveQuery, updateQuery, deleteQuery } from '@/store/modules/query';
import { updateRoute } from '@/store/modules/route';
import { deleteHeader, fetchRouteHeaders, getHeaders, saveHeader, updateHeader } from '@/store/modules/header';
import { fetchRouteQueries, getQueries } from '@/store/modules/query';

@Component
export default class Request extends Vue {
    @Prop({ default: new Route }) route!: Route;
    @Prop({ default: 0 }) projectId!: number;

    @Watch('route', { deep: true, immediate: true })
    onRouteChange(value: Route): void {
        Object.assign(this.routeUpdate, { ...value })

        if (this.route.id) {
            if (this._routeHeaders.length === 0) {
                fetchRouteHeaders(this.$store, this.route.id)
                    .catch((error) => { console.log(error) })
            }

            if (this._routeQueries.length === 0) {
                fetchRouteQueries(this.$store, this.route.id)
                    .catch((error) => { console.log(error) })
            }
        }
    }

    // routeUpdate holds the input values bound on the page
    private routeUpdate = new Route;

    get HTTPMethodOptions(): string[] {
        return Object.keys(HTTPMethod)
    }

    private routeQueries: Query[] = [];
    get _routeQueries(): Query[] {
        return getQueries(this.$store).filter(query => query.routeId === this.route.id)
    }
    @Watch('_routeQueries', { deep: true, immediate: true })
    on_RouteQueriesChange(): void {
        this.routeQueries = JSON.parse(JSON.stringify(this._routeQueries));
    }

    private routeHeaders: Header[] = [];
    get _routeHeaders(): Header[] {
        return getHeaders(this.$store).filter(header => header.routeId === this.route.id)
    }
    @Watch('_routeHeaders', { deep: true, immediate: true})
    on_RouteHeadersChange(): void {
        this.routeHeaders = JSON.parse(JSON.stringify(this._routeHeaders))
    }

    private addQuery(): void {
        this.routeQueries.push({
            routeId: this.route.id,
            id: 0,
            name: '',
            value: '',
        })
    }

    private addHeader(): void {
        this.routeHeaders.push({
            routeId: this.route.id,
            name: '',
            value: ''
        })
    }

    private updateRequest() {
        updateRoute(this.$store, this.routeUpdate)
            .catch(error => { console.log(error) })
    }

    private hasQueryChanged(index: number): boolean {
        const _query = this._routeQueries[index];
        const query = this.routeQueries[index];

        if (query.name === _query.name && query.value === _query.value) return true;

        return false
    }

    private hasHeaderChanged(index: number): boolean {
        const _header = this._routeHeaders[index];
        const header = this.routeHeaders[index];

        if (header.name === _header.name && header.value === _header.value) return true;

        return false;
    }

    private saveQuery(query: Query): void {
        saveQuery(this.$store, query)
            .catch(error => { console.log(error) })
    }
    
    private updateQuery(query: Query): void {
        updateQuery(this.$store, query)
            .catch(error => { console.log(error) })
    }

    private deleteQuery(query: Query): void {
        deleteQuery(this.$store, query)
    }

    private saveHeader(header: Header) {
        saveHeader(this.$store, {
                ...header,
                routeId: this.route.id,
            })
            .catch(error => { console.error(error) })
    }

    private updateHeader(header: Header) {
        updateHeader(this.$store, {
                ...header,
                routeId: this.route.id,
            })
            .catch(error => { console.error(error) })
    }

    private deleteHeader(header: Header) {
        deleteHeader(this.$store, header)
            .catch(error => { console.error(error) })
    }
}
</script>