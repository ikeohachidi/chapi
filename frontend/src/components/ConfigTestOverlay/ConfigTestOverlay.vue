<template>
    <section class="bg-black bg-opacity-80 fixed flex justify-end top-0 left-0 w-screen h-screen z-10">
        <div class="h-full bg-white w-2/5 flex flex-col">
            <p class="flex items-center py-5 px-4 bg-gray-100">
                <span class="text-gray-400 cursor-pointer" @click="close">
                    <i class="ri-close-line"></i>
                </span>
                <span class="ml-6">
                    Response 
                    <span class="ml-1">({{ configResult.responseTime }} seconds)</span>
                </span>
                <button @click="reRunTest" class="ml-auto">Re-run</button>
            </p>

            <div class="flex flex-col items-center justify-center my-auto" v-if="isLoading">
                <img src="@/assets/loading.gif" alt="loading image">
                <p class="mt-3">Fetching resource...</p>
            </div>
            <div class="font-mono px-4 py-5 text-sm overflow-y-auto flex-grow" v-else>
                <p>Status Code: <span :class="[configResult.type ? 'text-green-600' : 'text-red-600']">{{ configResult.data.status }}</span></p>
                <p>Status Text: <span :class="[configResult.type ? 'text-green-600' : 'text-red-600']">{{ configResult.data.statusText }}</span></p>

                <p class="mt-2">
                    <pre v-if="configResult.data.responseType !== 'image'" class="border bg-gray-50 rounded-md overflow-x-auto">
                        <code>
                            {{ configResult.data.body }}
                        </code>
                    </pre>
                    <img :src="configResult.data.body" v-else>
                </p>
            </div>
        </div>
    </section>
</template>

<script lang='ts'>
import {Vue, Component, Prop} from 'vue-property-decorator';

import { ProxyConfigResult } from '@/types/HTTP';

@Component
export default class className extends Vue {
    @Prop({ default: null }) configResult!: ProxyConfigResult;
    @Prop({ default: false }) isLoading!: boolean;

    private close() {
        this.$emit('close')
    }

    private reRunTest() {
        this.$emit('redo-test')
    }
}
</script>