<template>
	<span 
		class="switch" 
		:class="{ 'checked': value }"
		@click="updateValue"
	>
	</span>
</template>

<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator';

@Component
export default class InputSwitch extends Vue {
	@Prop({ default: false }) value!: boolean;
	
	private updateValue() {
		const update = this.value ? false : true;
		this.$emit('input', update);
		this.$emit('change', update);
	}
}
</script>

<style scoped>
.switch {
	@apply inline-block relative rounded-full overflow-hidden;
	@apply bg-gray-100 border border-gray-200 cursor-pointer;
	height: 22px;
	width: 40px;
}

.switch.checked {
	@apply bg-indigo-700;
}

.switch:before {
	@apply absolute top-0 bottom-0 rounded-full;
	@apply bg-white shadow-md;
	@apply translate-x-0 transform transition-all;
	content: '';
	width: 20px;
}

.switch.checked:before {
	@apply transform translate-x-full;
}
</style>