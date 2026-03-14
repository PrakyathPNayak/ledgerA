import axios from 'axios'

import { getValidToken, signOut } from './firebase'

const api = axios.create({
    baseURL: '/api/v1',
})

api.interceptors.request.use(async (config) => {
    const token = await getValidToken()
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

api.interceptors.response.use(
    (response) => {
        if (response.data && typeof response.data === 'object' && 'data' in response.data) {
            return response.data.data
        }
        return response.data
    },
    async (error) => {
        if (error?.response?.status === 401) {
            await signOut()
            window.location.href = '/login'
        }
        return Promise.reject(error)
    },
)

export default api
