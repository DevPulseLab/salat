import HomeView from '@/views/HomeView.vue'
import RegisterView from '@/views/RegisterView.vue'

export default [
    {
        name: 'home',
        path: '/',
        component: HomeView,
        meta: {
            label: 'Login',
            requiresAuth: false
        },
    },
    {
        name: 'register',
        path: '/register',
        component: RegisterView,
        meta: {
            label: 'Register',
            requiresAuth: false
        },
    },
]