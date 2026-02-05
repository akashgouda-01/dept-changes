import { useState, useEffect } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useAuth } from '@/contexts/AuthContext';
import { useToast } from '@/contexts/ToastContext';
import { User, Mail, BadgeCheck, Building2, Users, Plus, Pencil, Trash2, Search } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

// Hardcoded mock data for popup (Placeholder)
const MOCK_STUDENTS_PER_SECTION: Record<string, { regNo: string, name: string, email: string }[]> = {
  'A': [
    { regNo: '24CS0001', name: 'Student A1', email: 'a1@citchennai.net' },
    { regNo: '24CS0002', name: 'Student A2', email: 'a2@citchennai.net' },
  ],
  'B': [
    { regNo: '24CS0003', name: 'Student B1', email: 'b1@citchennai.net' },
    { regNo: '24CS0004', name: 'Student B2', email: 'b2@citchennai.net' },
  ],
  // Fallback for others
  'default': [
    { regNo: '24CSXXXX', name: 'Mock Student 1', email: 'mock1@citchennai.net' },
    { regNo: '24CSYYYY', name: 'Mock Student 2', email: 'mock2@citchennai.net' },
  ]
};

export default function FacultyProfile() {
  const { user } = useAuth();
  const { toast } = useToast();

  const [selectedSection, setSelectedSection] = useState<string | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const handleSectionClick = (section: string) => {
    setSelectedSection(section);
    setIsDialogOpen(true);
  };

  const getStudentsForSection = (section: string) => {
    return MOCK_STUDENTS_PER_SECTION[section] || MOCK_STUDENTS_PER_SECTION['default'];
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

        <div className="section-list-card">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-2 bg-primary/10 rounded-lg text-primary">
              <Users size={24} />
            </div>
            <div>
              <h2 className="text-xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text text-transparent">
                Assigned Sections
              </h2>
              <p className="text-sm text-muted-foreground">Click on a section to view student details</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {user?.assignedSections?.map((section) => (
              <button
                key={section}
                onClick={() => handleSectionClick(section)}
                className="group relative flex flex-col items-center justify-center p-8 bg-card hover:bg-gradient-to-br hover:from-primary/5 hover:to-transparent border border-border/60 rounded-xl transition-all duration-300 hover:shadow-lg hover:-translate-y-1"
              >
                <div className="absolute top-4 right-4 text-muted-foreground/20 group-hover:text-primary/20 transition-colors">
                  <BadgeCheck size={80} strokeWidth={1} />
                </div>

                <div className="relative z-10 text-center space-y-2">
                  <span className="text-4xl font-bold text-foreground group-hover:text-primary transition-colors">
                    {section}
                  </span>
                  <div className="h-1 w-12 bg-border group-hover:bg-primary/50 mx-auto rounded-full transition-colors" />
                  <p className="text-sm font-medium text-muted-foreground uppercase tracking-wider">
                    Section
                  </p>
                </div>
              </button>
            ))}
          </div>
        </div>

        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Students in Section {selectedSection}</DialogTitle>
              <DialogDescription>
                List of students assigned to this section.
              </DialogDescription>
            </DialogHeader>

            <div className="max-h-[60vh] overflow-y-auto mt-4 border rounded-md">
              <table className="w-full text-sm">
                <thead className="bg-muted/50 sticky top-0 text-left">
                  <tr>
                    <th className="p-3 font-medium text-muted-foreground">Reg No.</th>
                    <th className="p-3 font-medium text-muted-foreground">Name</th>
                    <th className="p-3 font-medium text-muted-foreground">Email</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  {selectedSection && getStudentsForSection(selectedSection).map((student, idx) => (
                    <tr key={idx} className="hover:bg-muted/20 transition-colors">
                      <td className="p-3 font-mono text-xs">{student.regNo}</td>
                      <td className="p-3 font-medium">{student.name}</td>
                      <td className="p-3 text-muted-foreground">{student.email}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  );
}
