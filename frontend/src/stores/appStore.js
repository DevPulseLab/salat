import {ref} from 'vue'

const appMessage = ref(null)
let appMessageTimeout = null

const useAppStore = () => {
    const setAppMessage = (statusCode, message, hideInMs = 4000) => {
        let type = 'error'
        switch (statusCode) {
            case 200:
                type = 'success'
                break;
            case 400:
                type = 'warn'
                break;
        }

        appMessage.value = { type, message }
        if (hideInMs > 0) {
            window.clearTimeout(appMessageTimeout)
            appMessageTimeout = window.setTimeout(() => {
                appMessage.value = null
            }, hideInMs)
        }
    }

    return {
        setAppMessage,
        appMessage
    }
}

export default useAppStore