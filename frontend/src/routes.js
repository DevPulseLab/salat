import HomeView from '@/views/HomeView.vue'
import RegisterView from '@/views/RegisterView.vue'

export default [
    {
        name: 'home',
        path: '/',
        component: HomeView,
        meta: {
            requiresAuth: false
        },
    },
    {
        name: 'register',
        path: '/register',
        component: RegisterView,
        meta: {
            requiresAuth: false
        },
    },
]