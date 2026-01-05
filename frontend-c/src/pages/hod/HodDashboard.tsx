import { useEffect, useState } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useAuth } from '@/contexts/AuthContext';
import { Modal } from '@/components/ui/Modal';
import { Users, FileCheck, TrendingUp, Download, CheckCircle2, XCircle, Clock, ChevronRight } from 'lucide-react';

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

// Prefer env; fall back to local backend default to avoid blank UI when env is missing locally.
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export default function HodDashboard() {
  const { user } = useAuth();
  const [overview, setOverview] = useState<DashboardOverview | null>(null);
  const [sections, setSections] = useState<SectionStat[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedSection, setSelectedSection] = useState<null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadDashboard = async () => {
      try {
        if (!API_BASE_URL) {
          throw new Error('VITE_API_BASE_URL is not configured');
        }

        setIsLoading(true);
        setError(null);

        const [overviewRes, sectionsRes] = await Promise.all([
          fetch(`${API_BASE_URL}/dashboard/overview`),
          fetch(`${API_BASE_URL}/dashboard/sections`),
        ]);

        if (!overviewRes.ok) {
          throw new Error(`Overview request failed with status ${overviewRes.status}`);
        }
        if (!sectionsRes.ok) {
          throw new Error(`Sections request failed with status ${sectionsRes.status}`);
        }

        const overviewJson = await overviewRes.json();
        const sectionsJson = await sectionsRes.json();

        if (!overviewJson.success) {
          throw new Error(overviewJson.error?.message || 'Failed to load dashboard overview');
        }
        if (!sectionsJson.success) {
          throw new Error(sectionsJson.error?.message || 'Failed to load section statistics');
        }

        if (isMounted) {
          setOverview(overviewJson.data as DashboardOverview);
          setSections(Array.isArray(sectionsJson.data) ? (sectionsJson.data as SectionStat[]) : []);
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
          <button className="btn btn-outline">
            <Download /> Export All Data
          </button>
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
                        <span>Total: {section.total_certificates}</span>
                        <span>Verified: {section.verified_count}</span>
                        <span>Pending: {section.pending_count}</span>
                        <span>Rejected: {section.rejected_count}</span>
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
