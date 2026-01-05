import { useState } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useToast } from '@/contexts/ToastContext';
import { Upload, Link as LinkIcon, Plus, Trash2, CheckCircle2, XCircle, Clock, Download, Eye, Cpu, ShieldCheck } from 'lucide-react';

interface CertificateEntry { id: string; driveLink: string; registerNumber: string; section: string; }
interface Certificate { id: string; studentName: string; registerNumber: string; section: string; driveLink: string; mlStatus: 'pending' | 'verified' | 'duplicate'; facultyStatus: 'pending' | 'legit' | 'not_legit'; mlScore?: number; }

export default function CertificateVerification() {
  const [entries, setEntries] = useState<CertificateEntry[]>([{ id: '1', driveLink: '', registerNumber: '', section: '' }]);
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [activeTab, setActiveTab] = useState<'pending' | 'legit' | 'not_legit'>('pending');
  const { toast } = useToast();

  const addEntry = () => { if (entries.length >= 10) return; setEntries([...entries, { id: Date.now().toString(), driveLink: '', registerNumber: '', section: '' }]); };
  const removeEntry = (id: string) => { if (entries.length === 1) return; setEntries(entries.filter(e => e.id !== id)); };
  const updateEntry = (id: string, field: keyof CertificateEntry, value: string) => { setEntries(entries.map(e => e.id === id ? { ...e, [field]: value } : e)); };

  const handleUpload = () => {
    const validEntries = entries.filter(e => e.driveLink && e.registerNumber && e.section);
    if (validEntries.length === 0) { toast({ title: 'No Valid Entries', description: 'Please fill in all fields.', variant: 'destructive' }); return; }
    toast({ title: 'Certificates Uploaded', description: `${validEntries.length} certificate(s) sent for verification.` });
    setEntries([{ id: '1', driveLink: '', registerNumber: '', section: '' }]);
  };

  const handleVerify = (certId: string, isLegit: boolean) => {
    setCertificates(certificates.map(c => c.id === certId ? { ...c, facultyStatus: isLegit ? 'legit' : 'not_legit' } : c));
    toast({ title: isLegit ? 'Certificate Verified' : 'Certificate Rejected' });
  };

  const pendingCerts = certificates.filter(c => c.mlStatus === 'verified' && c.facultyStatus === 'pending');
  const legitCerts = certificates.filter(c => c.facultyStatus === 'legit');
  const notLegitCerts = certificates.filter(c => c.facultyStatus === 'not_legit');
  const displayCerts = activeTab === 'pending' ? pendingCerts : activeTab === 'legit' ? legitCerts : notLegitCerts;

  return (
    <DashboardLayout requiredRole="faculty">
      <div className="space-y-8">
        <div className="dashboard-header">
          <div><h1 className="dashboard-title">Certificate Verification</h1><p className="dashboard-subtitle">Upload and verify student certificates</p></div>
          <button className="btn btn-outline"><Download /> Export to Excel</button>
        </div>

        <div className="upload-section">
          <div className="upload-header"><div className="upload-icon"><Upload /></div><div><h2 className="upload-title">Upload Certificates</h2><p className="upload-subtitle">Add up to 10 Google Drive links per batch</p></div></div>
          <div className="upload-entries">
            {entries.map((entry, index) => (
              <div key={entry.id} className="upload-entry">
                <div className="upload-entry-number">{index + 1}</div>
                <div className="upload-field"><span className="upload-field-label">Google Drive Link</span><div className="input-with-icon-wrapper"><LinkIcon /><input placeholder="https://drive.google.com/..." value={entry.driveLink} onChange={(e) => updateEntry(entry.id, 'driveLink', e.target.value)} className="input input-with-icon" /></div></div>
                <div className="upload-field"><span className="upload-field-label">Register Number</span><input placeholder="RA2211003XXX" value={entry.registerNumber} onChange={(e) => updateEntry(entry.id, 'registerNumber', e.target.value)} className="input font-mono" /></div>
                <div className="upload-field"><span className="upload-field-label">Section</span><select value={entry.section} onChange={(e) => updateEntry(entry.id, 'section', e.target.value)} className="input"><option value="">Select</option><option value="A">Section A</option><option value="B">Section B</option></select></div>
                <div className="upload-entry-actions"><button className="btn btn-ghost btn-icon" onClick={() => removeEntry(entry.id)} disabled={entries.length === 1} style={{color:'var(--destructive)'}}><Trash2 /></button></div>
              </div>
            ))}
          </div>
          <div className="upload-buttons"><button className="btn btn-outline" onClick={addEntry} disabled={entries.length >= 10}><Plus /> Add More ({entries.length}/10)</button><button className="btn btn-primary" onClick={handleUpload}><Upload /> Upload & Verify</button></div>
        </div>

        <div className="space-y-6">
          <div className="tabs-list">
            <button className={`tab-trigger ${activeTab === 'pending' ? 'active-warning' : ''}`} onClick={() => setActiveTab('pending')}><Clock /> Pending ({pendingCerts.length})</button>
            <button className={`tab-trigger ${activeTab === 'legit' ? 'active-success' : ''}`} onClick={() => setActiveTab('legit')}><CheckCircle2 /> Legitimate ({legitCerts.length})</button>
            <button className={`tab-trigger ${activeTab === 'not_legit' ? 'active-destructive' : ''}`} onClick={() => setActiveTab('not_legit')}><XCircle /> Not Legitimate ({notLegitCerts.length})</button>
          </div>

          {displayCerts.length === 0 ? (
            <div className="empty-state"><ShieldCheck /><p>No certificates in this category</p></div>
          ) : (
            <div className="space-y-4">
              {displayCerts.map((cert) => (
                <div key={cert.id} className="certificate-card">
                  <div className="certificate-card-content">
                    <div className="certificate-info">
                      <div className="certificate-icon"><Cpu /></div>
                      <div className="certificate-details">
                        <h3 className="certificate-student-name">{cert.studentName}</h3>
                        <p className="certificate-register-number">{cert.registerNumber}</p>
                        <div className="certificate-badges"><span className="badge badge-outline">Section {cert.section}</span>{cert.mlScore && <span className={`badge ${cert.mlScore >= 80 ? 'badge-success' : 'badge-warning'}`}>ML Score: {cert.mlScore}%</span>}</div>
                      </div>
                    </div>
                    <div className="certificate-actions">
                      <a href={cert.driveLink} target="_blank" rel="noopener noreferrer" className="btn btn-outline btn-sm"><Eye /> View</a>
                      {activeTab === 'pending' && (<><button className="btn btn-outline btn-sm" style={{borderColor:'var(--success)',color:'var(--success)'}} onClick={() => handleVerify(cert.id, true)}><CheckCircle2 /> Legitimate</button><button className="btn btn-outline btn-sm" style={{borderColor:'var(--destructive)',color:'var(--destructive)'}} onClick={() => handleVerify(cert.id, false)}><XCircle /> Not Legit</button></>)}
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
