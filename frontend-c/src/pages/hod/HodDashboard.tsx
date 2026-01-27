import { useEffect, useState } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useAuth } from '@/contexts/AuthContext';
import { Modal } from '@/components/ui/Modal';
import { Users, FileCheck, TrendingUp, Download, CheckCircle2, XCircle, Clock } from 'lucide-react';
import { getDashboardOverview, getDashboardSections } from '@/lib/api';
import { useNavigate } from 'react-router-dom';

interface DashboardOverview {
  total_students: number;
  total_certificates: number;
  verified_count: number;
  rejected_count: number;
  pending_count: number;
}

interface SectionStat {
  section: string;
  total_certificates: number;
  verified_count: number;
  rejected_count: number;
  pending_count: number;
  verification_rate: number;
}

export default function HodDashboard() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [overview, setOverview] = useState<DashboardOverview | null>(null);
  const [sections, setSections] = useState<SectionStat[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [, setSelectedSection] = useState<null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadDashboard = async () => {
      try {
        setIsLoading(true);
        setError(null);

        const [overviewRes, sectionsRes] = await Promise.all([
          getDashboardOverview(),
          getDashboardSections(),
        ]);

        if (isMounted) {
          // Check for success flag if API returns standard wrapper, or direct data if api.ts unwraps it.
          // api.ts returns response.data. The controllers return { success: true, data: ... }
          // So overviewRes is ALL of { success, data }.
          if (overviewRes.success) {
            setOverview(overviewRes.data as DashboardOverview);
          }
          if (sectionsRes.success) {
            setSections(Array.isArray(sectionsRes.data) ? (sectionsRes.data as SectionStat[]) : []);
          }
        }

      } catch (err) {
        console.error(err);
        if (isMounted) {
          setError('Failed to load dashboard data. Please try again later.');
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    loadDashboard();

    return () => {
      isMounted = false;
    };
  }, []);

  const totalStudents = overview?.total_students ?? 0;
  const totalCertificates = overview?.total_certificates ?? 0;
  const totalLegit = overview?.verified_count ?? 0;
  const totalPending = overview?.pending_count ?? 0;
  const overallLegitRate =
    totalCertificates > 0 ? Math.round((totalLegit / totalCertificates) * 100) : 0;

  return (
    <DashboardLayout requiredRole="hod">
      <div className="space-y-8">
        <div className="dashboard-header">
          <div>
            <h1 className="dashboard-title">HOD Dashboard</h1>
            <p className="dashboard-subtitle">Welcome, <span>{user?.name}</span> - CSE Department Overview</p>
          </div>
          <div className="flex gap-2">
            <button className="btn btn-outline" onClick={() => navigate('/hod/faculty-stats')}>
              <Users /> Faculty Stats Lookup
            </button>
            <button className="btn btn-outline">
              <Download /> Export All Data
            </button>
          </div>
        </div>

        {isLoading ? (
          <div className="stats-grid">
            <div className="stat-card skeleton" />
            <div className="stat-card skeleton" />
            <div className="stat-card skeleton" />
            <div className="stat-card skeleton" />
          </div>
        ) : error ? (
          <div className="section-card">
            <p className="text-destructive font-medium">{error}</p>
          </div>
        ) : (
          <div className="stats-grid">
            <div className="stat-card stat-card-primary">
              <div className="stat-card-content">
                <div className="stat-card-info">
                  <p className="stat-card-title">Total Students</p>
                  <p className="stat-card-value">{totalStudents}</p>
                  <p className="stat-card-subtitle">All sections</p>
                </div>
                <div className="stat-card-icon"><Users /></div>
              </div>
            </div>
            <div className="stat-card stat-card-warning">
              <div className="stat-card-content">
                <div className="stat-card-info">
                  <p className="stat-card-title">Total Certificates</p>
                  <p className="stat-card-value">{totalCertificates}</p>
                  <p className="stat-card-subtitle">Uploaded this semester</p>
                </div>
                <div className="stat-card-icon"><FileCheck /></div>
              </div>
            </div>
            <div className="stat-card stat-card-success">
              <div className="stat-card-content">
                <div className="stat-card-info">
                  <p className="stat-card-title">Verified Certificates</p>
                  <p className="stat-card-value">{totalLegit}</p>
                  <p className="stat-card-subtitle">{overallLegitRate}% verification rate</p>
                </div>
                <div className="stat-card-icon"><CheckCircle2 /></div>
              </div>
            </div>
            <div className="stat-card stat-card-primary">
              <div className="stat-card-content">
                <div className="stat-card-info">
                  <p className="stat-card-title">Pending Review</p>
                  <p className="stat-card-value">{totalPending}</p>
                  <p className="stat-card-subtitle">Awaiting verification</p>
                </div>
                <div className="stat-card-icon"><Clock /></div>
              </div>
            </div>
          </div>
        )}

        {/* Section-wise Statistics */}
        <div className="content-grid">
          <div className="section-card">
            <div className="section-card-header">
              <h2 className="section-card-title">Section-wise Certificates</h2>
              <span className="section-card-subtitle">Overview</span>
            </div>
            {sections.length === 0 ? (
              <div className="section-progress empty-state">
                <p>No section data available yet.</p>
              </div>
            ) : (
              <div className="section-progress">
                {sections.map((section) => {
                  const sectionVerified =
                    section.total_certificates > 0
                      ? Math.round(
                        (section.verified_count / section.total_certificates) * 100,
                      )
                      : 0;
                  return (
                    <div key={section.section} className="section-progress-item">
                      <div className="section-progress-header">
                        <span className="section-progress-title">Section {section.section}</span>
                        <span className="section-progress-value">
                          {sectionVerified}% verified
                        </span>
                      </div>
                      <div className="section-progress-bar">
                        <div
                          className="section-progress-bar-fill"
                          style={{ width: `${sectionVerified}%` }}
                        />
                      </div>
                      <div className="section-progress-meta">
                        <span>Total: {section.total_certificates} </span>
                        <span>Verified: {section.verified_count} </span>
                        <span>Pending: {section.pending_count} </span>
                        <span>Rejected: {section.rejected_count} </span>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>

          <div className="section-card">
            <div className="section-card-header">
              <h2 className="section-card-title">Certificate Status</h2>
              <TrendingUp />
            </div>
            <div className="certificate-status">
              <div className="certificate-status-main">
                <p className="certificate-status-value">{overallLegitRate}%</p>
                <p className="certificate-status-label">Verification Rate</p>
              </div>
              <div className="certificate-status-items">
                <div className="certificate-status-item warning">
                  <div className="certificate-status-item-label warning"><Clock /> Pending</div>
                  <span className="certificate-status-item-value">{totalPending}</span>
                </div>
                <div className="certificate-status-item success">
                  <div className="certificate-status-item-label success"><CheckCircle2 /> Verified</div>
                  <span className="certificate-status-item-value">{totalLegit}</span>
                </div>
                <div className="certificate-status-item destructive">
                  <div className="certificate-status-item-label destructive"><XCircle /> Rejected</div>
                  <span className="certificate-status-item-value">{overview?.rejected_count ?? 0}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <Modal
        isOpen={false}
        onClose={() => setSelectedSection(null)}
        title=""
      >
        <div>{/* Modal content intentionally disabled until section-level APIs are available */}</div>
      </Modal>
    </DashboardLayout>
  );
}
