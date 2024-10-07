import axios from 'axios'
import useUserStore from '@/stores/userStore'
import useAppStore from '@/stores/appStore'

const appStore = useAppStore()
const userStore = useUserStore()

const http = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    headers: {
        'Content-Type': 'application/json'
    }
})

const handleError = (error) => {
    if (error.response === undefined) {
        console.error(error)
        appStore.setAppMessage(500, 'Unknown server error', 0)
        return null
    }

    if (error.response.status !== undefined) {
        appStore.setAppMessage(error.response.status, error.response.data.error)
    }

    return error.response
}

http.interceptors.request.use((config) => {
    const user = userStore.getUser()
    if (user.token !== null) {
        config.headers['Authorization'] = 'Bearer ' + user.token;
    }

    return config
}, (error) => {
    return handleError(error)
})

http.interceptors.response.use(function (response) {
    return response
}, function (error) {
    return handleError(error)
})

export default http