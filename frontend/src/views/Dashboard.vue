<template>
    <section class="dashboard">
        <div class="left-pane border-r border-gray-200">
            <ProjectNav />
        </div>
        <div class="right-pane overflow-y-auto bg-blue-100 bg-opacity-20">
            <router-view/>
        </div>
    </section>
</template>

<script lang='ts'>
import {Vue, Component} from 'vue-property-decorator';

import ProjectNav from '@/components/ProjectNav/ProjectNav.vue';

@Component({
    components: {
        ProjectNav
    }
})
export default class Dashboard extends Vue {}
</script>

<style lang="scss" scoped>
.dashboard {
    @apply border-t border-gray-200;
    height: calc(100vh - 63px);
    display: flex;
}

.active {
    background-color: rgba(0, 153, 255, 0.233);
}

.left-pane {
    width: 25%;
    max-width: 350px;

    ul {
        list-style-type: none;
        padding-left: 0;

        li {
            padding: 10px;
        }
    }
}

.right-pane {
    width: 75%;
    flex-grow: 1;
    padding: 0 0;
}
</style>