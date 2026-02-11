import { apiClient } from './client';

export interface Student {
  id: number;
  registerNumber: string;
  name: string;
  email: string;
  section: string;
  semester: number;
  isPresent: boolean;
  facultyEmail: string;
}

export const fetchStudentsBySection = async (section: string): Promise<Student[]> => {
  const response = await apiClient.get<Student[]>('/faculty/students', {
    params: { section },
  });
  return response.data;
};

export const fetchStudentByRegNo = async (regNo: string): Promise<Student> => {
  const response = await apiClient.get<Student>(`/faculty/student/${regNo}`);
  return response.data;
};
