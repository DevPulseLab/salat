import {ref} from 'vue'

const user = ref({
    token: localStorage.getItem('token') ?? null
})

const useUserStore = () => {
    const getUser = () => {
        return user.value
    }

    const setUserToken = (token) => {
        user.value = {
            token
        }

        if (token === null) {
            localStorage.removeItem('token')
            return
        }

        localStorage.setItem('token', token)
    }

    const isAuthenticated = () => {
        return user.value.token !== null
    }

    return {
        isAuthenticated,
        getUser,
        setUserToken
    }
}

export default useUserStore