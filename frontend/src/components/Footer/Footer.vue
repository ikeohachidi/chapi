<template>
    <footer class="font-ibm">
        <p class="text-center">&copy; Ikeoha Chidi 21 til Infinity</p>
    </footer>
</template>

<script lang='ts'>
import {Vue, Component} from 'vue-property-decorator';

@Component
export default class Footer extends Vue {

}
</script>

<style lang="scss" scoped>
footer {
    @apply bg-gray-100 text-gray-700 py-6;
}
</style>