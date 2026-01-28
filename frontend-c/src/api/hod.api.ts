import { apiClient } from './client';
import { Certificate, StudentStats } from '@/types';

// HOD METHODS

export const getStudentStats = async (facultyId: string) => {
    const { data } = await apiClient.get<{ data: StudentStats[] }>(`/hod/faculty/students?faculty_id=${facultyId}`);
    return data.data;
};

export const getStudentCertificates = async (regNo: string) => {
    const { data } = await apiClient.get<{ data: Certificate[] }>(`/hod/student/certificates?reg_no=${regNo}`);
    return data.data;
};

export const exportCertificatesBySection = async (section: string) => {
    const response = await apiClient.get(`/hod/export/certificates/section?section=${section}`, {
        responseType: 'blob',
    });
    return response;
};

export const exportCertificatesByStudent = async (regNo: string) => {
    const response = await apiClient.get(`/hod/export/certificates/student?reg_no=${regNo}`, {
        responseType: 'blob',
    });
    return response;
};
