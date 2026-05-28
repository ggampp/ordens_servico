import axios from 'axios'

// Centralized axios instance. The JWT is attached on every request and a 401
// response clears the session and redirects to login.
const client = axios.create({
  baseURL: '/api/v1',
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(err)
  }
)

// apiError extracts a human-readable message from an axios error.
export function apiError(err) {
  return err?.response?.data?.error?.message || err.message || 'Erro inesperado'
}

export default client
