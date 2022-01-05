<template>
  <div id="app">
		<p class="bg-green-600 text-white text-center font-bold py-3">Thanks for visiting, Please Note Chapi is still a work in progress.</p>
    <navbar/>
    <router-view/>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';

import Navbar from '@/components/Navbar/Navbar.vue';

import { fetchAuthUser } from '@/store/modules/user';

@Component({
  components: {
    Navbar
  }
})
export default class App extends Vue {
  created(): void {
    fetchAuthUser(this.$store); 
  }
}
</script>