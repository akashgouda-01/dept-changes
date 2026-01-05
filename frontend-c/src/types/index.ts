export type UserRole = 'faculty' | 'hod';

export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  staffId: string;
  position: string;
  department: string;
  assignedSections?: string[];
}

export interface Student {
  id: string;
  registerNumber: string;
  name: string;
  email: string;
  section: string;
  semester: number;
  isPresent: boolean;
}

export interface Certificate {
  id: string;
  driveLink: string;
  registerNumber: string;
  section: string;
  studentName: string;
  uploadedBy: string;
  uploadedAt: Date;
  verificationStatus: 'pending' | 'ml_verified' | 'faculty_verified' | 'rejected';
  isLegit: boolean | null;
  mlScore?: number;
}

export interface Section {
  id: string;
  name: string;
  totalStudents: number;
  presentStudents: number;
  absentStudents: number;
  totalCertificates: number;
  legitCertificates: number;
}
