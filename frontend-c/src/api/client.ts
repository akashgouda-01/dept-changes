import axios from 'axios';

// Prefer env; fall back to local backend default to avoid blank UI when env is missing locally.
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// Create Axios Instance
export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request Interceptor: Attach Token
apiClient.interceptors.request.use(
    (config) => {
        // Read the user object from localStorage
        const storedUser = localStorage.getItem('eduvault_user');
        if (storedUser) {
            try {
                const user = JSON.parse(storedUser);
                // Header format: Bearer <email>|<role>
                if (user.email && user.role) {
                    config.headers.Authorization = `Bearer ${user.email}|${user.role}`;
                }
            } catch (error) {
                console.error('Failed to parse user from local storage:', error);
            }
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Response Interceptor: Error Handling
apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401 || error.response?.status === 403) {
            console.warn('Unauthorized access:', error.config.url);
            // in a real app, you might trigger a logout or redirect here
        }
        return Promise.reject(error);
    }
);
