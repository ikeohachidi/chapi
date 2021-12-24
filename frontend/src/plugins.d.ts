import Vue from 'vue';

declare module 'vue/types/vue'  {
    interface Vue {
        $toast: {
            success: (arg: string) => void,
            error: (arg: string) => void,
            info: (arg: string) => void,
        };
    }
}