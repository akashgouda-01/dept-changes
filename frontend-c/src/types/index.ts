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

export type MLStatus = 'PENDING' | 'VERIFIED' | 'DUPLICATE';
export type FacultyStatus = 'PENDING' | 'LEGIT' | 'NOT_LEGIT';

export interface Certificate {
  ID: string; // Changed to match Go JSON output (usually PascalCase if not tagged, but Gin/Gorm usually tags json? Wait, let's check backend tags)
  DriveLink: string;
  RegisterNumber: string;
  Section: string;
  StudentName: string;
  UploadedBy: string;
  UploadedAt: string; // Date comes as string in JSON
  MLStatus: MLStatus;
  FacultyStatus: FacultyStatus;
  IsLegit: boolean | null;
  MLScore?: number;
  Archived: boolean;
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

export interface StudentStats {
  register_number: string;
  student_name: string;
  section: string;
  total_certificates: number;
  verified_count: number;
  rejected_count: number;
  pending_count: number;
}

export interface UploadCertificatePayload {
  drive_link: string;
  register_number: string;
  section: string;
  student_name: string;
  uploaded_by: string;
  uploaded_at: string;
}
