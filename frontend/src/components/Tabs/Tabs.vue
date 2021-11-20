<script lang='ts'>
import Vue, { CreateElement, VNode } from 'vue';
import { Component, Prop } from 'vue-property-decorator';
import { TabProps } from './Tab.vue';

@Component
export default class Tabs extends Vue {
    @Prop({ default: 'left' }) position!: string;

    private activeTabTitle = '';

    get defaultSlots(): VNode[] {
        return this.$slots.default || [];
    }

    private getTabByTitle(title: string): VNode {
        const index = this.defaultSlots.findIndex(slot => (slot.componentOptions?.propsData as TabProps).title === title)

        return this.defaultSlots[index];
    }

    created(): void {
        if (this.defaultSlots.length > 0) {
            this.activeTabTitle = (this.defaultSlots[0].componentOptions?.propsData as TabProps).title;
        }
    }

    render(createElement: CreateElement): VNode {
        return createElement(
            'div', 
            [
                createElement(
                    'div', 
                    {
                        class: {
                            'tab-headers': true,
                            'justify-center': this.position === 'center',
                            'justify-end': this.position === 'right'
                        }
                    },
                    this.defaultSlots.map(slot => {
                        const tabProps = (slot.componentOptions?.propsData as TabProps);
        
                        return createElement(
                                'p',
                                {
                                    class: {
                                        'tab-header': true,
                                        'active': this.activeTabTitle === tabProps.title
                                    },
                                    on: {
                                        click: (() => {
                                            this.activeTabTitle = tabProps.title;
                                        })
                                    }
                                },
                                [tabProps.title]
                            )
                    })
                ),
                createElement(
                    'div',
                    { 
                        class: ['tab-body'] 
                    },
                    [ this.getTabByTitle(this.activeTabTitle) ]
                )
            ]
        )
    }
}
</script>

<style lang="scss" scoped>
.tab-headers {
    @apply flex border-solid border-b border-gray-100;
}

.tab-body {
    @apply p-3;
}

.tab-header {
    @apply border-solid border-b-2 border-transparent cursor-pointer px-3 p-2 mb-0;

    &:hover {
        @apply bg-gray-100;
    }

    &.active {
        @apply border-indigo-700;
    }
}
</style>