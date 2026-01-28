import { useState } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useToast } from '@/contexts/ToastContext';
import { Search, Users, ChevronRight, Download } from 'lucide-react';
import { getStudentStats } from '@/api';
import { StudentStats } from '@/types';
import { useNavigate } from 'react-router-dom';

export default function HodFacultyStats() {
    const [facultyId, setFacultyId] = useState('');
    const [stats, setStats] = useState<StudentStats[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [hasSearched, setHasSearched] = useState(false);
    const { toast } = useToast();
    const navigate = useNavigate();

    const handleSearch = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!facultyId.trim()) return;

        setIsLoading(true);
        setHasSearched(true);
        try {
            const data = await getStudentStats(facultyId);
            setStats(data || []);
        } catch (error) {
            console.error(error);
            toast({ title: 'Error', description: 'Failed to find faculty stats or invalid ID.', variant: 'destructive' });
            setStats([]);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <DashboardLayout requiredRole="hod">
            <div className="space-y-8">
                <div className="dashboard-header">
                    <div>
                        <h1 className="dashboard-title">Faculty Statistics</h1>
                        <p className="dashboard-subtitle">View student progress by faculty mentor</p>
                    </div>
                </div>

                <div className="section-card">
                    <form onSubmit={handleSearch} className="flex gap-4 items-end">
                        <div className="flex-1 space-y-2">
                            <label htmlFor="facultyId" className="text-sm font-medium">Faculty ID</label>
                            <div className="input-with-icon-wrapper">
                                <Search className="w-4 h-4" />
                                <input
                                    id="facultyId"
                                    placeholder="e.g. FAC01"
                                    value={facultyId}
                                    onChange={(e) => setFacultyId(e.target.value)}
                                    className="input input-with-icon"
                                />
                            </div>
                        </div>
                        <button type="submit" className="btn btn-primary" disabled={isLoading}>
                            {isLoading ? 'Searching...' : 'Search'}
                        </button>
                    </form>
                </div>

                {hasSearched && (
                    <div className="section-card">
                        <div className="section-card-header mb-4">
                            <h2 className="section-card-title">Student List</h2>
                            {stats.length > 0 && <span className="badge badge-outline">{stats.length} Students</span>}
                        </div>

                        {stats.length === 0 ? (
                            <div className="empty-state">
                                <Users />
                                <p>No students found for this Faculty ID.</p>
                            </div>
                        ) : (
                            <div className="overflow-x-auto">
                                <table className="w-full text-sm text-left">
                                    <thead className="text-muted-foreground border-b border-border/50">
                                        <tr>
                                            <th className="py-3 px-4">Register No</th>
                                            <th className="py-3 px-4">Name</th>
                                            <th className="py-3 px-4">Section</th>
                                            <th className="py-3 px-4 text-center">Total</th>
                                            <th className="py-3 px-4 text-center text-success">Verified</th>
                                            <th className="py-3 px-4 text-center text-destructive">Rejected</th>
                                            <th className="py-3 px-4 text-center text-warning">Pending</th>
                                            <th className="py-3 px-4"></th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {stats.map((student) => (
                                            <tr key={student.register_number} className="border-b border-border/50 hover:bg-muted/50 transition-colors">
                                                <td className="py-3 px-4 font-mono">{student.register_number}</td>
                                                <td className="py-3 px-4 font-medium">{student.student_name}</td>
                                                <td className="py-3 px-4">{student.section}</td>
                                                <td className="py-3 px-4 text-center">{student.total_certificates}</td>
                                                <td className="py-3 px-4 text-center font-bold text-success">{student.verified_count}</td>
                                                <td className="py-3 px-4 text-center text-destructive">{student.rejected_count}</td>
                                                <td className="py-3 px-4 text-center text-warning">{student.pending_count}</td>
                                                <td className="py-3 px-4 text-right">
                                                    <button
                                                        className="btn btn-ghost btn-sm"
                                                        onClick={() => navigate(`/hod/student/${student.register_number}`)}
                                                    >
                                                        View <ChevronRight className="w-4 h-4 ml-1" />
                                                    </button>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>
                )}
            </div>
        </DashboardLayout>
    );
}
