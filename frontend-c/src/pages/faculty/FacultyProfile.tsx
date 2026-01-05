import { useState } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useAuth } from '@/contexts/AuthContext';
import { useToast } from '@/contexts/ToastContext';
import { User, Mail, BadgeCheck, Building2, Users, Plus, Pencil, Trash2, Search } from 'lucide-react';
import { Student } from '@/types';

const initialStudents: Student[] = [
  { id: '1', registerNumber: 'RA2211003010', name: 'John Doe', email: 'john@ctchennai.net', section: 'A', semester: 5, isPresent: true },
  { id: '2', registerNumber: 'RA2211003015', name: 'Jane Smith', email: 'jane@ctchennai.net', section: 'A', semester: 5, isPresent: true },
  { id: '3', registerNumber: 'RA2211003022', name: 'Mike Wilson', email: 'mike@ctchennai.net', section: 'B', semester: 5, isPresent: false },
];

export default function FacultyProfile() {
  const { user } = useAuth();
  const { toast } = useToast();
  const [students, setStudents] = useState<Student[]>(initialStudents);
  const [searchQuery, setSearchQuery] = useState('');
  const [sectionFilter, setSectionFilter] = useState('all');

  const filteredStudents = students.filter(student => {
    const matchesSearch = student.name.toLowerCase().includes(searchQuery.toLowerCase()) || student.registerNumber.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesSection = sectionFilter === 'all' || student.section === sectionFilter;
    return matchesSearch && matchesSection;
  });

  const handleDeleteStudent = (studentId: string) => {
    const student = students.find(s => s.id === studentId);
    setStudents(students.filter(s => s.id !== studentId));
    toast({ title: 'Student Removed', description: `${student?.name} has been removed.` });
  };

  return (
    <DashboardLayout requiredRole="faculty">
      <div className="space-y-8">
        <div><h1 className="dashboard-title">Profile</h1><p className="dashboard-subtitle">Manage your profile and student list</p></div>

        <div className="profile-card">
          <div className="profile-content">
            <div className="profile-avatar"><div className="profile-avatar-box"><User /></div></div>
            <div className="profile-info">
              <h2 className="profile-name">{user?.name}</h2>
              <p className="profile-position">{user?.position}</p>
              <div className="profile-details">
                <div className="profile-detail-item"><BadgeCheck /><div><p className="profile-detail-label">Staff ID</p><p className="profile-detail-value">{user?.staffId}</p></div></div>
                <div className="profile-detail-item"><Mail /><div><p className="profile-detail-label">Email</p><p className="profile-detail-value">{user?.email}</p></div></div>
                <div className="profile-detail-item"><Building2 /><div><p className="profile-detail-label">Department</p><p className="profile-detail-value">{user?.department}</p></div></div>
                <div className="profile-detail-item"><Users /><div><p className="profile-detail-label">Assigned Sections</p><div className="profile-sections">{user?.assignedSections?.map(s => <span key={s} className="badge badge-outline">Section {s}</span>)}</div></div></div>
              </div>
            </div>
          </div>
        </div>

        <div className="student-list-card">
          <div className="student-list-header">
            <div className="student-list-title-section"><div className="student-list-icon"><Users /></div><div><h2 className="student-list-title">Assigned Students</h2><p className="student-list-count">{students.length} students total</p></div></div>
                        </div>

          <div className="filters-row">
            <div className="search-input-wrapper"><Search /><input type="text" placeholder="Search by name or register number..." value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)} className="input input-with-icon" /></div>
            <select value={sectionFilter} onChange={(e) => setSectionFilter(e.target.value)} className="input filter-select"><option value="all">All Sections</option><option value="A">Section A</option><option value="B">Section B</option></select>
          </div>

          <div className="table-container">
            <table className="table">
              <thead><tr><th>Register No.</th><th>Name</th><th>Email</th><th>Section</th><th style={{textAlign:'right'}}>Actions</th></tr></thead>
              <tbody>
                {filteredStudents.map((student) => (
                  <tr key={student.id}>
                    <td className="font-mono">{student.registerNumber}</td>
                    <td>{student.name}</td>
                    <td style={{color:'var(--muted-foreground)'}}>{student.email}</td>
                    <td><span className="badge badge-outline">Section {student.section}</span></td>
                    <td style={{textAlign:'right'}}><button className="btn btn-ghost btn-icon"></button><button className="btn btn-ghost btn-icon" onClick={() => handleDeleteStudent(student.id)} style={{color:'var(--destructive)'}}><Trash2 style={{width:'1rem',height:'1rem'}} /></button></td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
