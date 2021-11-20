<template>
    <section>
        <slot></slot>
    </section>
</template>

<script lang='ts'>

import {Vue, Component, Prop} from 'vue-property-decorator';

export type TabProps = {
    title: string
}

@Component
export default class Tab extends Vue {
    @Prop({ default: '', required: true }) title!: string;
}
</script>