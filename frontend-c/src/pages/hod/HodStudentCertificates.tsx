import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useToast } from '@/contexts/ToastContext';
import { Download, Cpu, ShieldCheck, Eye, CheckCircle2, XCircle, Clock } from 'lucide-react';
import { getStudentCertificates, exportCertificatesByStudent } from '@/lib/api';
import { Certificate } from '@/types';

export default function HodStudentCertificates() {
    const { regNo } = useParams<{ regNo: string }>();
    const [certificates, setCertificates] = useState<Certificate[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const { toast } = useToast();

    useEffect(() => {
        if (regNo) {
            loadCertificates(regNo);
        }
    }, [regNo]);

    const loadCertificates = async (reg: string) => {
        try {
            setIsLoading(true);
            const data = await getStudentCertificates(reg);
            setCertificates(data || []);
        } catch (error) {
            console.error(error);
            toast({ title: 'Error', description: 'Failed to find student certificates.', variant: 'destructive' });
        } finally {
            setIsLoading(false);
        }
    };

    const handleExport = async () => {
        if (!regNo) return;
        try {
            const response = await exportCertificatesByStudent(regNo);
            // Create blob link to download
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            // Extract filename from header or use default
            const contentDisposition = response.headers['content-disposition'];
            let filename = `certificates_${regNo}.xlsx`;
            if (contentDisposition) {
                const matches = /filename="([^"]*)"/.exec(contentDisposition);
                if (matches != null && matches[1]) {
                    filename = matches[1];
                }
            }
            link.setAttribute('download', filename);
            document.body.appendChild(link);
            link.click();
            link.remove();
        } catch (error) {
            console.error(error);
            toast({ title: 'Export Failed', description: 'Could not download Excel file.', variant: 'destructive' });
        }
    };

    return (
        <DashboardLayout requiredRole="hod">
            <div className="space-y-8">
                <div className="dashboard-header">
                    <div>
                        <h1 className="dashboard-title">Student Certificates</h1>
                        {regNo && <p className="dashboard-subtitle font-mono">{regNo}</p>}
                    </div>
                    <button className="btn btn-outline" onClick={handleExport} disabled={certificates.length === 0}>
                        <Download /> Export Excel
                    </button>
                </div>

                {isLoading ? (
                    <div className="p-8 text-center text-muted-foreground">Loading...</div>
                ) : certificates.length === 0 ? (
                    <div className="empty-state">
                        <ShieldCheck />
                        <p>No certificates found for this student.</p>
                    </div>
                ) : (
                    <div className="space-y-4">
                        {certificates.map((cert) => (
                            <div key={cert.ID} className="certificate-card">
                                <div className="certificate-card-content">
                                    <div className="certificate-info">
                                        <div className="certificate-icon"><Cpu /></div>
                                        <div className="certificate-details">
                                            <h3 className="certificate-student-name">{cert.StudentName}</h3>
                                            <p className="certificate-register-number">{cert.DriveLink}</p>
                                            <div className="certificate-badges">
                                                {cert.MLScore !== undefined && <span className={`badge ${cert.MLScore >= 80 ? 'badge-success' : 'badge-warning'}`}>ML Score: {cert.MLScore}%</span>}
                                                <span className={`badge badge-outline ${cert.MLStatus === 'VERIFIED' ? 'text-success border-success' : ''}`}>ML: {cert.MLStatus}</span>
                                                <span className={`badge badge-outline ${cert.FacultyStatus === 'LEGIT' ? 'text-success border-success' : cert.FacultyStatus === 'NOT_LEGIT' ? 'text-destructive border-destructive' : ''}`}>Faculty: {cert.FacultyStatus}</span>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="certificate-actions">
                                        <a href={cert.DriveLink} target="_blank" rel="noopener noreferrer" className="btn btn-outline btn-sm"><Eye /> View</a>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </DashboardLayout>
    );
}
