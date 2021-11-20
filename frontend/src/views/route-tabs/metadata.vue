<template>
    <section class="section content-padding">
        <div>
            <p class="section-name">Description</p>
            <p class="section-description">
                It would be best if your description explains what your route does clearly. Proxies created can get out of hand. 
            </p>
        </div>
        <div class="flex flex-col">
            <textarea class="w-full resize-none" rows="6" placeholder="Route Description" v-model="description"></textarea>
            <button 
                class="mt-4 ml-auto"
                :disabled="route.description === description"
                @click="updateMetadata"
            >
                Save
            </button>
        </div>
    </section>
</template>

<script lang='ts'>
import {Vue, Component, Prop, Watch} from 'vue-property-decorator';

import Route from '@/types/Route';
import { updateRoute } from '@/store/modules/route';

@Component
export default class Metadata extends Vue {
    @Prop({ default: new Route }) route!: Route;

    private description = '';

    @Watch('route', { deep: true, immediate: true })
    onRouteChange(): void {
        this.description = this.route.description.slice();
    }

    private updateMetadata(): void {
        updateRoute(this.$store, {
            ...this.route,
            description: this.description
        })
        .catch(error => { console.log(error) })
    }
}
</script>