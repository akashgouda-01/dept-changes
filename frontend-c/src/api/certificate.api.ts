import { apiClient } from './client';
import { Certificate, FacultyStatus, UploadCertificatePayload } from '@/types';

// FACULTY METHODS

export const uploadCertificates = async (certificates: UploadCertificatePayload[]) => {
    const { data } = await apiClient.post('/certificates/upload', { certificates });
    return data;
};

export const getPendingReviewCertificates = async (limit = 50) => {
    const { data } = await apiClient.get<{ data: Certificate[] }>(`/certificates/pending-review?limit=${limit}`);
    return data.data;
};

export const submitReview = async (certificateId: string, status: FacultyStatus, isLegit: boolean) => {
    const { data } = await apiClient.post('/certificates/review', {
        certificate_id: certificateId,
        status,
        is_legit: isLegit,
    });
    return data;
};

export const triggerMockVerification = async (certificateId: string) => {
    const { data } = await apiClient.post('/faculty/certificate/verify', {
        certificate_id: certificateId,
    });
    return data;
};
