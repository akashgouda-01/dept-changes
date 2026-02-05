import axios from 'axios';
import { auth } from '@/config/firebase';

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
    async (config) => {
        // Get fresh token from Firebase Auth
        const user = auth.currentUser;
        if (user) {
            try {
                const token = await user.getIdToken();
                config.headers.Authorization = `Bearer ${token}`;
            } catch (error) {
                console.error("Error getting auth token", error);
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
