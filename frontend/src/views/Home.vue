<template>
	<section>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4 w-11/12 max-w-6xl mx-auto pt-20">
			<div>
				<h1 class="text-2xl md:text-4xl font-bold mb-10">A simple and painless proxy for your JAMSTACK application</h1>
				<p class="mb-6">Chapi is a simple application aimed at the JAMSTACK or frontend developer who wants to access some api’s but doesn’t want to have to spin up an entire backend to hide API keys or doesn’t want to have to write complex cloud functions. It’s quick and easy to setup and get going!!</p>

				<div class="flex">
					<button @click="routeToDashboard">Get started</button>
					<button class="github-btn" @click="routeToGithub">
						<i class="ri-github-fill mr-2"></i>
						Github	
					</button>
				</div>
			</div>
			<div class="hidden md:block">
				<img src="@/assets/landing-page.svg" alt="Landing page" class="h-4/5 ml-auto">
			</div>
		</div>

		<div class="w-11/12 max-w-6xl mx-auto mt-20 md:mt-0">
			<div class="md:w-3/5 mb-9">
				<h1 class="text-3xl font-bold mb-3">Reasons to use</h1>
				<p>Chapi aims to be as lean and as simple as possible, maybe for a lack of imagination, laziness or just wanting it that way. At the time of writing it’s stable.</p>
			</div>

			<div class="reasons">
				<div>
					<p>Free to use</p>
					<p>Chapi is and always will be free to use. You can the hosted version or you can self host if you like. Noteworthy though is the fact i’m pressed for resources and so it may not be hosted on a $20,000/month machine.</p>
				</div>
				<div>
					<p>Open source</p>
					<p>It’s open source and publicly available on Github, with an MIT license, so that means you can feel free to do whatever you want it, tinker around, log issues and submit PR’s if you feel like it.</p>
				</div>
				<div>
					<p>Considerably secure</p>
					<p>😃 I took all the steps to the best of my knowledge to make Chapi secure. But i’m not well versed enough in web security to not make some mistakes. If you find any do well to inform me, create an Issue on Github or submit a PR.</p>
				</div>
			</div>
		</div>
	</section>
</template>

<script lang='ts'>
import {Vue, Component} from 'vue-property-decorator';

@Component
export default class Home extends Vue {
	private routeToDashboard() {
		this.$router.push({ name: 'Dashboard' });
	}

	private routeToGithub() {
		window.open('https://github.com/ikeohachidi/chapi', '_blank');
	}
}
</script>

<style lang="scss" scoped>
p {
	@apply text-gray-800 leading-7;
}

.github-btn {
	@apply ml-4 text-gray-800 bg-gray-300 flex items-center;
	@apply transition-all hover:bg-gray-400;
}

.reasons {
	@apply grid grid-cols-1 md:grid-cols-3 gap-9;

	div {
		p:first-of-type {
			@apply text-xl font-bold mb-2;
		}
	}
}
</style>