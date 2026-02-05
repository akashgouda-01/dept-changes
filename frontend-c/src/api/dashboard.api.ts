import { apiClient } from './client';

export interface DashboardOverview {
    total_students: number;
    total_certificates: number;
    verified_count: number;
    rejected_count: number;
    pending_count: number;
}

export interface SectionStat {
    section: string;
    total_certificates: number;
    verified_count: number;
    rejected_count: number;
    pending_count: number;
    verification_rate: number;
}

export interface RecentActivity {
    id: string;
    student_name: string;
    reg_no: string;
    section: string;
    action: string;
    timestamp: string;

}

export interface Student {
    ID: number;
    RegisterNumber: string;
    Name: string;
    Email: string;
    Section: string;
    Semester: number;
    IsPresent: boolean;
    FacultyEmail: string;
}

// DASHBOARD METHODS

export const getDashboardOverview = async () => {
    const { data } = await apiClient.get<{ success: boolean; data: DashboardOverview }>('/dashboard/overview');
    return data; // Returns full response object, caller checks .success
};

export const getDashboardSections = async () => {
    const { data } = await apiClient.get<{ success: boolean; data: SectionStat[] }>('/dashboard/sections');
    return data;
};

export const getRecentActivity = async () => {
    const { data } = await apiClient.get<{ success: boolean; data: RecentActivity[] }>('/dashboard/recent-activity');
    return data;
};
export const getAssignedStudents = async () => {
    const { data } = await apiClient.get<{ success: boolean; data: Student[] }>('/dashboard/students');
    return data;
};

export const exportMyCertificates = async () => {
    const response = await apiClient.get('/dashboard/export/certificates', {
        responseType: 'blob',
    });
    return response;
};
