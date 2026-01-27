import axios from 'axios';
import { Certificate, FacultyStatus, StudentStats } from '@/types';

// Create Axios Instance
export const api = axios.create({
    baseURL: 'http://localhost:8080',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request Interceptor: Attach Token
api.interceptors.request.use(
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
api.interceptors.response.use(
    (response) => response,
    (error) => {
        // Example: Redirect to login on 401/403?
        // For now, just reject so React Query/Components can handle it.
        if (error.response?.status === 401 || error.response?.status === 403) {
            console.warn('Unauthorized access:', error.config.url);
            // Optional: window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

// --- API Methods ---

// FACULTY
export const uploadCertificates = async (certificates: any[]) => {
    const { data } = await api.post('/certificates/upload', { certificates });
    return data;
};

export const getPendingReviewCertificates = async (limit = 50) => {
    const { data } = await api.get<{ data: Certificate[] }>(`/certificates/pending-review?limit=${limit}`);
    return data.data; // The controller returns { data: [...] }
};

export const submitReview = async (certificateId: string, status: FacultyStatus, isLegit: boolean) => {
    const { data } = await api.post('/certificates/review', {
        certificate_id: certificateId,
        status,
        is_legit: isLegit,
    });
    return data;
};

export const triggerMockVerification = async (certificateId: string) => {
    const { data } = await api.post('/faculty/certificate/verify', {
        certificate_id: certificateId,
    });
    return data;
};

// DASHBOARD
export const getDashboardOverview = async () => {
    const { data } = await api.get('/dashboard/overview');
    return data;
};

export const getDashboardSections = async () => {
    const { data } = await api.get('/dashboard/sections');
    return data;
};

// HOD
export const getStudentStats = async (facultyId: string) => {
    const { data } = await api.get<{ data: StudentStats[] }>(`/hod/faculty/students?faculty_id=${facultyId}`);
    return data.data;
};

export const getStudentCertificates = async (regNo: string) => {
    const { data } = await api.get<{ data: Certificate[] }>(`/hod/student/certificates?reg_no=${regNo}`);
    return data.data;
};

export const exportCertificatesBySection = async (section: string) => {
    // For file download, we might need responseType: 'blob'
    const response = await api.get(`/hod/export/certificates/section?section=${section}`, {
        responseType: 'blob',
    });
    return response;
};

export const exportCertificatesByStudent = async (regNo: string) => {
    const response = await api.get(`/hod/export/certificates/student?reg_no=${regNo}`, {
        responseType: 'blob',
    });
    return response;
};
