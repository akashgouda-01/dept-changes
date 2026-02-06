import { useState, useEffect } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useAuth } from '@/contexts/AuthContext';
import { useToast } from '@/contexts/ToastContext';
import { User, Mail, BadgeCheck, Building2, Users, Loader2 } from 'lucide-react';
import { fetchStudentsBySection, Student } from '@/api/student.api';

export default function FacultyProfile() {
  const { user } = useAuth();
  const { toast } = useToast();
  const [students, setStudents] = useState<Student[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const loadAllStudents = async () => {
      // Wait for user data to be available
      if (!user || !user.assignedSections || user.assignedSections.length === 0) {
        return;
      }

      setIsLoading(true);
      try {
        console.log("Fetching students for sections:", user.assignedSections);

        // Fetch students for all sections in parallel
        const promises = user.assignedSections.map(section =>
          fetchStudentsBySection(section)
            .catch(err => {
              console.error(`Failed to fetch section ${section}:`, err);
              return []; // Return empty array on failure so other sections still load
            })
        );

        const results = await Promise.all(promises);

        // Flatten the array of arrays into a single list of students
        const allStudents = results.flat();

        // Sort by register number or name if desired (optional)
        allStudents.sort((a, b) => a.registerNumber.localeCompare(b.registerNumber));

        console.log(`Fetched total ${allStudents.length} students`);
        setStudents(allStudents);

      } catch (error) {
        console.error("Error fetching students:", error);
        toast({
          title: "Error fetching students",
          description: "Could not load student data.",
          variant: "destructive",
        });
      } finally {
        setIsLoading(false);
      }
    };

    loadAllStudents();
  }, [user, toast]);

  return (
    <DashboardLayout requiredRole="faculty">
      <div className="space-y-8">
        <div>
          <h1 className="dashboard-title">Profile</h1>
          <p className="dashboard-subtitle">Manage your profile and student list</p>
        </div>

        <div className="profile-card">
          <div className="profile-content">
            <div className="profile-avatar">
              <div className="profile-avatar-box"><User /></div>
            </div>
            <div className="profile-info">
              <h2 className="profile-name">{user?.name}</h2>
              <p className="profile-position">{user?.position}</p>
              <div className="profile-details">
                <div className="profile-detail-item">
                  <BadgeCheck />
                  <div>
                    <p className="profile-detail-label">Staff ID</p>
                    <p className="profile-detail-value">{user?.staffId}</p>
                  </div>
                </div>
                <div className="profile-detail-item">
                  <Mail />
                  <div>
                    <p className="profile-detail-label">Email</p>
                    <p className="profile-detail-value">{user?.email}</p>
                  </div>
                </div>
                <div className="profile-detail-item">
                  <Building2 />
                  <div>
                    <p className="profile-detail-label">Department</p>
                    <p className="profile-detail-value">{user?.department}</p>
                  </div>
                </div>
                <div className="profile-detail-item">
                  <Users />
                  <div>
                    <p className="profile-detail-label">Assigned Sections</p>
                    <div className="profile-sections">
                      {user?.assignedSections?.map(s => <span key={s} className="badge badge-outline">Section {s}</span>)}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-card border rounded-lg shadow-sm">
          <div className="p-6 border-b">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 rounded-lg text-primary">
                <Users size={24} />
              </div>
              <h2 className="text-xl font-bold">All Students</h2>
            </div>
            <p className="text-sm text-muted-foreground mt-2">
              Viewing {students.length} students from your assigned sections.
            </p>
          </div>

          <div className="overflow-x-auto">
            {isLoading ? (
              <div className="flex flex-col items-center justify-center py-10">
                <Loader2 className="h-8 w-8 animate-spin text-primary mb-2" />
                <p className="text-sm text-muted-foreground">Loading student data...</p>
              </div>
            ) : students.length > 0 ? (
              <table className="w-full text-sm">
                <thead className="bg-muted/50 text-left">
                  <tr>
                    <th className="p-4 font-medium text-muted-foreground">Reg No.</th>
                    <th className="p-4 font-medium text-muted-foreground">Name</th>
                    <th className="p-4 font-medium text-muted-foreground">Section</th>
                    <th className="p-4 font-medium text-muted-foreground">Email</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  {students.map((student) => {
                    // Debug log to check data integrity
                    if (!student.name || !student.email) console.warn("Missing data for student:", student);

                    return (
                      <tr key={student.registerNumber} className="hover:bg-muted/20 transition-colors">
                        <td className="p-4 font-mono text-xs text-primary/80">{student.registerNumber}</td>
                        <td className="p-4 font-medium">{student.name || '-'}</td>
                        <td className="p-4 font-medium">
                          <span className="inline-flex items-center px-2 py-1 rounded-md bg-secondary text-secondary-foreground text-xs">
                            Section {student.section}
                          </span>
                        </td>
                        <td className="p-4 text-muted-foreground">{student.email || '-'}</td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            ) : (
              <div className="p-8 text-center text-muted-foreground">
                No students found in your assigned sections.
              </div>
            )}
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
