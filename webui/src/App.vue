<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { computed } from 'vue'
import { isLoggedIn } from './services/auth.js'

const route = useRoute()
const brandHref = computed(() => isLoggedIn() ? '#/' : '#/login') //totalnie nie potrzebne a le nawet fajne
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a :href="brandHref" class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6">WASAtext</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="app-container">
		<div class="row">
			<nav v-if="route.path !== '/login'" id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Menu</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/profile" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Profile
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/conversations" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#inbox"/></svg>
								Conversations
							</RouterLink>
						</li>
					</ul>

				</div>
			</nav>

			<main :class="route.path === '/login' ? 'col-12 px-md-4' : 'col-md-9 ms-sm-auto col-lg-10 px-md-4'">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
.app-container {
	max-width: 1100px; 
	margin: 0 auto;    
	padding-left: 1rem;
	padding-right: 1rem;
	font-size: 1.0rem; 
	line-height: 1.6;
}

@media (max-width: 767.98px) {
	.app-container {
		padding-left: 0.5rem;
		padding-right: 0.5rem;
		font-size: 1rem;
	}
}
</style>
