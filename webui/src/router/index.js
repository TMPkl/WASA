import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import { isLoggedIn } from '../services/auth.js'

const routes = [
	{ path: '/', component: HomeView },
	{ path: '/login', component: LoginView },
	{ path: '/profile', component: ProfileView },
]

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes,
})

router.beforeEach((to, from, next) => {
	const logged = isLoggedIn()
	if (to.path === '/') {
		return logged ? next('/profile') : next('/login')
	}
	if (to.path === '/login' && logged) {
		return next('/profile')
	}
	if (!logged && to.path !== '/login') {
		return next('/login')
	}
	next()
})

export default router