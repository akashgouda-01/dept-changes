import { useState, useEffect } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useToast } from '@/contexts/ToastContext';
import { useAuth } from '@/contexts/AuthContext';
import { Upload, Link as LinkIcon, Plus, Trash2, CheckCircle2, XCircle, Clock, Download, Eye, Cpu, ShieldCheck, User } from 'lucide-react';
import { Certificate, UploadCertificatePayload } from '@/types';
import { uploadCertificates, getPendingReviewCertificates, submitReview } from '@/lib/api';

interface CertificateEntry {
  id: string;
  driveLink: string;
  registerNumber: string;
  section: string;
  studentName: string;
}

export default function CertificateVerification() {
  const [entries, setEntries] = useState<CertificateEntry[]>([{ id: '1', driveLink: '', registerNumber: '', section: '', studentName: '' }]);
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [activeTab, setActiveTab] = useState<'pending' | 'legit' | 'not_legit'>('pending');
  const [isLoading, setIsLoading] = useState(false);
  const [isUploading, setIsUploading] = useState(false);

  const { toast } = useToast();
  const { user } = useAuth();

  const fetchPending = async () => {
    try {
      setIsLoading(true);
      const data = await getPendingReviewCertificates();
      setCertificates(data || []);
    } catch (error) {
      console.error(error);
      toast({ title: 'Error', description: 'Failed to load pending certificates.', variant: 'destructive' });
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchPending();
  }, []);

  const addEntry = () => {
    if (entries.length >= 10) return;
    setEntries([...entries, { id: Date.now().toString(), driveLink: '', registerNumber: '', section: '', studentName: '' }]);
  };

  const removeEntry = (id: string) => {
    if (entries.length === 1) return;
    setEntries(entries.filter(e => e.id !== id));
  };

  const updateEntry = (id: string, field: keyof CertificateEntry, value: string) => {
    setEntries(entries.map(e => e.id === id ? { ...e, [field]: value } : e));
  };

  const handleUpload = async () => {
    const validEntries = entries.filter(e => e.driveLink && e.registerNumber && e.section && e.studentName);
    if (validEntries.length === 0) {
      toast({ title: 'No Valid Entries', description: 'Please fill in all fields (including Student Name).', variant: 'destructive' });
      return;
    }

    try {
      setIsUploading(true);
      const payload: UploadCertificatePayload[] = validEntries.map(e => ({
        drive_link: e.driveLink,
        register_number: e.registerNumber,
        section: e.section,
        student_name: e.studentName,
        uploaded_by: user?.staffId || 'UNKNOWN',
        uploaded_at: new Date().toISOString()
      }));

      await uploadCertificates(payload);

      toast({ title: 'Success', description: `${validEntries.length} certificate(s) uploaded successfully.` });
      setEntries([{ id: Date.now().toString(), driveLink: '', registerNumber: '', section: '', studentName: '' }]);

      // Refresh list (though they might be in ML pending state, not Faculty pending yet, but good to refresh)
      fetchPending();
    } catch (error: any) {
      toast({ title: 'Upload Failed', description: error.response?.data?.message || 'Failed to upload certificates.', variant: 'destructive' });
    } finally {
      setIsUploading(false);
    }
  };

  const handleVerify = async (certId: string, isLegit: boolean) => {
    try {
      await submitReview(certId, isLegit ? 'LEGIT' : 'NOT_LEGIT', isLegit);

      // Optimistic update
      setCertificates(certificates.filter(c => c.ID !== certId));

      toast({ title: isLegit ? 'Certificate Verified' : 'Certificate Rejected', variant: isLegit ? 'default' : 'destructive' });
    } catch (error: any) {
      toast({ title: 'Action Failed', description: error.response?.data?.message || 'Failed to submit review.', variant: 'destructive' });
    }
  };

  // Filter local state? Currently we only fetch PENDING from API.
  // Unless we store history, the other tabs will remain empty or we need another API.
  // For now we only show pending.
  const displayCerts = certificates;

  return (
    <DashboardLayout requiredRole="faculty">
      <div className="space-y-8">
        <div className="dashboard-header">
          <div><h1 className="dashboard-title">Certificate Verification</h1><p className="dashboard-subtitle">Upload and verify student certificates</p></div>
          <button className="btn btn-outline" onClick={fetchPending} disabled={isLoading}><Clock /> Refresh</button>
        </div>

        <div className="upload-section">
          <div className="upload-header"><div className="upload-icon"><Upload /></div><div><h2 className="upload-title">Upload Certificates</h2><p className="upload-subtitle">Add up to 10 Google Drive links per batch</p></div></div>
          <div className="upload-entries">
            {entries.map((entry, index) => (
              <div key={entry.id} className="upload-entry">
                <div className="upload-entry-number">{index + 1}</div>

                <div className="upload-field">
                  <span className="upload-field-label">Student Name</span>
                  <div className="input-with-icon-wrapper"><User />
                    <input placeholder="John Doe" value={entry.studentName} onChange={(e) => updateEntry(entry.id, 'studentName', e.target.value)} className="input input-with-icon" />
                  </div>
                </div>

                <div className="upload-field">
                  <span className="upload-field-label">Register Number</span>
                  <input placeholder="24CS0001" value={entry.registerNumber} onChange={(e) => updateEntry(entry.id, 'registerNumber', e.target.value)} className="input font-mono" />
                </div>  

                <div className="upload-field">
                  <span className="upload-field-label">Section</span>
                  <select value={entry.section} onChange={(e) => updateEntry(entry.id, 'section', e.target.value)} className="input"><option value="">Select</option><option value="A">Section A</option><option value="B">Section B</option></select>
                </div>

                <div className="upload-field" style={{ flex: 2 }}>
                  <span className="upload-field-label">Google Drive Link</span>
                  <div className="input-with-icon-wrapper"><LinkIcon />
                    <input placeholder="https://drive.google.com/..." value={entry.driveLink} onChange={(e) => updateEntry(entry.id, 'driveLink', e.target.value)} className="input input-with-icon" />
                  </div>
                </div>

                <div className="upload-entry-actions"><button className="btn btn-ghost btn-icon" onClick={() => removeEntry(entry.id)} disabled={entries.length === 1} style={{ color: 'var(--destructive)' }}><Trash2 /></button></div>
              </div>
            ))}
          </div>
          <div className="upload-buttons"><button className="btn btn-outline" onClick={addEntry} disabled={entries.length >= 10}><Plus /> Add More ({entries.length}/10)</button><button className="btn btn-primary" onClick={handleUpload} disabled={isUploading}>{isUploading ? 'Uploading...' : <><Upload /> Upload & Verify</>}</button></div>
        </div>

        <div className="space-y-6">
          <div className="tabs-list">
            <button className={`tab-trigger ${activeTab === 'pending' ? 'active-warning' : ''}`} onClick={() => setActiveTab('pending')}><Clock /> Pending ({certificates.length})</button>
            {/* Disabled tabs for history as endpoint is missing */}
            <button className={`tab-trigger disabled opacity-50 cursor-not-allowed`}><CheckCircle2 /> Legitimate</button>
            <button className={`tab-trigger disabled opacity-50 cursor-not-allowed`}><XCircle /> Not Legitimate</button>
          </div>

          {(isLoading) ? (
            <div className="p-8 text-center text-muted-foreground">Loading...</div>
          ) : displayCerts.length === 0 ? (
            <div className="empty-state"><ShieldCheck /><p>No pending certificates found.</p></div>
          ) : (
            <div className="space-y-4">
              {displayCerts.map((cert) => (
                <div key={cert.ID} className="certificate-card">
                  <div className="certificate-card-content">
                    <div className="certificate-info">
                      <div className="certificate-icon"><Cpu /></div>
                      <div className="certificate-details">
                        <h3 className="certificate-student-name">{cert.StudentName}</h3>
                        <p className="certificate-register-number">{cert.RegisterNumber}</p>
                        <div className="certificate-badges"><span className="badge badge-outline">Section {cert.Section}</span>{cert.MLScore !== undefined && <span className={`badge ${cert.MLScore >= 80 ? 'badge-success' : 'badge-warning'}`}>ML Score: {cert.MLScore}%</span>}</div>
                      </div>
                    </div>
                    <div className="certificate-actions">
                      <a href={cert.DriveLink} target="_blank" rel="noopener noreferrer" className="btn btn-outline btn-sm"><Eye /> View</a>
                      <button className="btn btn-outline btn-sm" style={{ borderColor: 'var(--success)', color: 'var(--success)' }} onClick={() => handleVerify(cert.ID, true)}><CheckCircle2 /> Legitimate</button>
                      <button className="btn btn-outline btn-sm" style={{ borderColor: 'var(--destructive)', color: 'var(--destructive)' }} onClick={() => handleVerify(cert.ID, false)}><XCircle /> Not Legit</button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}
