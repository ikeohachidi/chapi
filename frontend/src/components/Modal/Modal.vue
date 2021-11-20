<template>
    <section class="bg-black bg-opacity-80 flex items-center justify-center fixed top-0 left-0 w-screen h-screen z-10">
        <div class="bg-white overflow-hidden rounded-lg w-2/6">
            <slot name="title"></slot>
            <p class="text-lg p-4 font-bold border-b border-gray-300 text-gray-700" v-if="title">{{ title }}</p>

            <div class="p-4">
                <p class="mb-4 text-gray-600" v-if="description">{{ description }}</p>
                <slot></slot>
            </div>

            <div class="bg-gray-200 py-4 flex justify-end px-4 mt-3">
                <p @click="close" class="mb-0 py-2 px-3 rounded-md mr-2 text-sm hover:bg-gray-300 cursor-pointer">Close</p>
                <button :disabled="!enableOK" @click="actionButtonClicked">{{ actionButtonText }}</button>
            </div>
        </div>
    </section>
</template>

<script lang='ts'>
import {Vue, Component, Prop} from 'vue-property-decorator';

@Component
export default class Modal extends Vue {
    @Prop({ default: '' }) title!: string;
    @Prop({ default: '' }) description!: string;
    @Prop({ default: 'Done' }) actionButtonText!: string;
    @Prop({ default: true }) enableOK!: boolean;

    private close() {
        this.$emit('close')
    }

    private actionButtonClicked() {
        this.$emit('action')
    }
}
</script>