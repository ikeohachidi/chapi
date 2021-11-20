import StoreState from '@/store/storeState';

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $store: StoreState;
  }
}
