import {ref} from 'vue'

const user = ref({
    token: null
})

const useUserStore = () => {

    const getUser = () => {
        return user.value
    }

    const setUserToken = (token) => {
        user.value = {
            token
        }
    }

    return {
        getUser,
        setUserToken
    }
}

export default useUserStore