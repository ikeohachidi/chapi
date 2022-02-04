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
        <div class="grid grid-cols-3 gap-7 content-padding py-7 border-b border-gray-200">
            <div>
                <p class="section-name flex">
                    Merge Request Body
                    <input-switch class="ml-2" v-model="mergeOptions.mergeBody" @change="updateMergeOptions($event)"/>
                </p>
                <p class="section-description">
                    Allowing this option means that the request body sent when fetching the projects provided endpoint gets merged with the configured <span class="inline-code">request body</span> set in the
                    <span class="inline-code">Request</span> tab
                </p>
            </div>
            <div>
                <p class="section-name flex">
                    Merge Request Headers
                    <input-switch class="ml-2" v-model="mergeOptions.mergeHeader" @change="updateMergeOptions($event)"/>
                </p>
                <p class="section-description">
                    Allowing this option means that the request headers sent when fetching the projects provided endpoint gets merged with the configured <span class="inline-code">headers</span> set in the
                    <span class="inline-code">Request</span> tab
                </p>
            </div>
            <div>
                <p class="section-name flex">
                    Merge Request Query 
                    <input-switch class="ml-2" v-model="mergeOptions.mergeQuery" @change="updateMergeOptions($event)"/>
                </p>
                <p class="section-description">
                    Allowing this option means that all url queries added to this Chapi URL when making a request will be appended to the set <span class="inline-code">destination</span> URL
                    <span class="inline-code">Request</span> tab
                </p>
            </div>
        </div>
    </section>    
</template>

<script lang="ts">
import { Vue, Component, Watch, Prop } from 'vue-property-decorator';
import InputSwitch from '@/components/InputSwitch/InputSwitch.vue';

import { 
    fetchPermOrigins,
    createPermOrigin,
    deletePermOrigin,
    getPermOrigins,
    updatePermOrigin
} from '@/store/modules/perm-origin';
import { 
    getMergeOptions, 
    fetchMergeOptions,
    updateMergeOption
} from '@/store/modules/merge-options';

import { PermOrigin, MergeOptions } from '@/types/Security';
import Route from '@/types/Route';

@Component({
    components: {
        InputSwitch
    }
})
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

    // private mergeOptions = new MergeOptions(); 
    get mergeOptions(): MergeOptions {
        return getMergeOptions(this.$store)[this.route.id!] || new MergeOptions();
    }

    private updateMergeOptions(): void {
        updateMergeOption(this.$store, this.mergeOptions)
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
        createPermOrigin(this.$store, permOrigin)
            .then(() => this.$toast.success('Origin saved successfully'))
            .catch(() => this.$toast.error('Error saving origin'));
    }

    private updatePermOrigin(permOrigin: PermOrigin): void {
        updatePermOrigin(this.$store, permOrigin)
            .then(() => this.$toast.success('Origin updated successfully'))
            .catch(() => this.$toast.error('Error updating origin'));
    }

    private deletePermOrigin(permOrigin: PermOrigin, index: number): void {
        if (permOrigin.id === 0) {
            this.routeOrigins.splice(index, 1);
            return;
        }
        deletePermOrigin(this.$store, permOrigin)
            .then(() => this.$toast.success('Origin deleted successfully'))
            .catch(() => this.$toast.error('Error deleting origin'));
    }
    
    mounted(): void {
        if (this.route.id) {
            const requests: Promise<unknown>[] = [];

            if (!this.mergeOptions.routeId) requests.push(fetchMergeOptions(this.$store, this.route.id));
            if (this.routeOrigins.length === 0) requests.push(fetchPermOrigins(this.$store, this.route.id))

            Promise.all(requests)
            .catch(error => {
                this.$toast.error(error)
            })
        }
    }
}
</script>