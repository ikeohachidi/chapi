<template>
    <section>
        <div class="section content-padding">
            <div>
                <p class="section-name">Whitelisted URLs</p>
                <p class="section-description">
                    URLs listed here will be the only ones allowed to access the resource. Leaving it blank will mean all URLs are allowed to retrieve this resource.
                </p>
            </div>
            <div class="flex flex-col">
                <table class="w-full">
                    <thead>
                        <tr class="text-left text-gray-800">
                            <th class="py-2 pl-3 font-normal">URL</th>
                            <th class="px-3 font-normal"></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(origin, originIndex) in routeOrigins" :key="originIndex">
                            <td class="px-3 w-4/5 py-2">
                                <input class="w-full" type="text" v-model="origin.url" placeholder="https://">
                            </td>
                            <td class="w-1/5">
                                <div class="flex items-center justify-end">
                                    <button @click="updatePermOrigin(origin)" :disabled="!hasOriginChanged(originIndex)" v-if="origin.id > 0">
                                        Update
                                    </button>
                                    <button @click="savePermOrigin(origin)" v-else>
                                        Save 
                                    </button>
                                    <span class="text-red-600 cursor-pointer ml-4 hover:scale-50" @click="deletePermOrigin(origin, originIndex)">
                                        <i class="ri-delete-bin-line"></i>
                                    </span>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <button class="mt-4 ml-auto" @click="addOrigin">Add URL</button>
            </div>
        </div>
    </section>    
</template>

<script lang="ts">
import { Vue, Component, Watch, Prop } from 'vue-property-decorator';

import { createPermOrigin, deletePermOrigin, fetchPermOrigins, getPermOrigins, updatePermOrigin } from '@/store/modules/perm-origin';

import { PermOrigin } from '@/types/Security';
import Route from '@/types/Route';

@Component
export default class Security extends Vue {
    @Prop({ default: new Route }) route!: Route;
    
    private routeOrigins: PermOrigin[] = [];

    get _routeOrigins(): PermOrigin[] {
        return getPermOrigins(this.$store)[this.route.id!] || [];
    }
    @Watch('_routeOrigins', { deep: true })
    on_RouteOriginsChange(value: PermOrigin[]): void {
        this.routeOrigins = JSON.parse(JSON.stringify(value));
    }

    private addOrigin(): void {
        this.routeOrigins.push({ 
            id: 0,
            url: '',
            routeId: this.route.id!
        })
    }

    private hasOriginChanged(originIndex: number): boolean {
        const _routeOrign = this._routeOrigins[originIndex];
        const routeOrigin = this.routeOrigins[originIndex];

        if (_routeOrign.url === routeOrigin.url) {
            return false;
        }

        return true;
    }

    private savePermOrigin(permOrigin: PermOrigin): void {
        // TODO: handle promise
        createPermOrigin(this.$store, permOrigin)
    }

    private updatePermOrigin(permOrigin: PermOrigin): void {
        // TODO: handle promise
        updatePermOrigin(this.$store, permOrigin)
    }

    private deletePermOrigin(permOrigin: PermOrigin, index: number): void {
        if (permOrigin.id === 0) {
            this.routeOrigins.splice(index, 1);
            return;
        }
        // TODO: handle promise
        deletePermOrigin(this.$store, permOrigin)
    }
    
    mounted(): void {
        if (this.routeOrigins.length === 0 && this.route.id) {
            // TODO: handle promise
            fetchPermOrigins(this.$store, this.route.id)
        }
    }
}
</script>