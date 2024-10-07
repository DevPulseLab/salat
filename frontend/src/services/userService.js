import http from '@/http.js'
import useUserStore from '@/stores/userStore'

const useUserService = () => {
    const login = async (username, password) => {
        const response = await http.post('/api/login', {
            username,
            password
        })

        if (response.status !== 200) {
            return false
        }

        const { setUserToken } = useUserStore()
        setUserToken(response.data.token)

        return true
    }

    const logout = async () => {
        const response = await http.post('/api/logout')

        const { setUserToken } = useUserStore()
        setUserToken(null)

        return response.status === 200
    }

    const register = async (username, password) => {
        const response = await http.post('/api/register', {
            username,
            password
        })

        return {
            success: response.status === 200,
            message: response.data.message
        }
    }

    return {
        login,
        logout,
        register,
    }
}

export default useUserService