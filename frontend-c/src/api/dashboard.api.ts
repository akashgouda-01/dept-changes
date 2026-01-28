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

// DASHBOARD METHODS

export const getDashboardOverview = async () => {
    const { data } = await apiClient.get<{ success: boolean; data: DashboardOverview }>('/dashboard/overview');
    return data; // Returns full response object, caller checks .success
};

export const getDashboardSections = async () => {
    const { data } = await apiClient.get<{ success: boolean; data: SectionStat[] }>('/dashboard/sections');
    return data;
};
