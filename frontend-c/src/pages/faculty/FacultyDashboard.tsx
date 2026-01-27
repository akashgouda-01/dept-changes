import { useEffect, useState } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useAuth } from '@/contexts/AuthContext';
import { Users, FileCheck, Clock, CheckCircle2, XCircle, TrendingUp } from 'lucide-react';
import { getDashboardOverview, getDashboardSections } from '@/lib/api';

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

export default function FacultyDashboard() {
  const { user } = useAuth();
  const [overview, setOverview] = useState<DashboardOverview | null>(null);
  const [sections, setSections] = useState<SectionStat[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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
  const verified = overview?.verified_count ?? 0;
  const pending = overview?.pending_count ?? 0;
  const rejected = overview?.rejected_count ?? 0;
  const verificationPercentage =
    totalCertificates > 0 ? Math.round((verified / totalCertificates) * 100) : 0;

  return (
    <DashboardLayout requiredRole="faculty">
      <div className="space-y-8">
        <div className="dashboard-header">
          <div>
            <h1 className="dashboard-title">Dashboard</h1>
            <p className="dashboard-subtitle">Welcome back, <span>{user?.name}</span></p>
          </div>
          <div className="dashboard-date">
            <Clock />
            <span>{new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</span>
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
            <p className="text-destructive font-medium">{
              error
            }</p>
          </div>
        ) : (
          <div className="stats-grid">
            <div className="stat-card stat-card-primary">
              <div className="stat-card-content">
                <div className="stat-card-info">
                  <p className="stat-card-title">Total Students</p>
                  <p className="stat-card-value">{totalStudents}</p>
                  <p className="stat-card-subtitle">Across all sections</p>
                </div>
                <div className="stat-card-icon"><Users /></div>
              </div>
            </div>
            <div className="stat-card stat-card-warning">
              <div className="stat-card-content">
                <div className="stat-card-info">
                  <p className="stat-card-title">Total Certificates</p>
                  <p className="stat-card-value">{totalCertificates}</p>
                  <p className="stat-card-subtitle">{pending} pending review</p>
                </div>
                <div className="stat-card-icon"><FileCheck /></div>
              </div>
            </div>
            <div className="stat-card stat-card-success">
              <div className="stat-card-content">
                <div className="stat-card-info">
                  <p className="stat-card-title">Verified</p>
                  <p className="stat-card-value">{verified}</p>
                  <p className="stat-card-subtitle">{verificationPercentage}% verification rate</p>
                </div>
                <div className="stat-card-icon"><CheckCircle2 /></div>
              </div>
            </div>
            <div className="stat-card stat-card-destructive">
              <div className="stat-card-content">
                <div className="stat-card-info">
                  <p className="stat-card-title">Rejected</p>
                  <p className="stat-card-value">{rejected}</p>
                  <p className="stat-card-subtitle">Requires resubmission</p>
                </div>
                <div className="stat-card-icon"><XCircle /></div>
              </div>
            </div>
          </div>
        )}

        {/* Section-wise and recent activity UI kept, but no mock numbers; will be wired to real APIs later */}
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
                <p className="certificate-status-value">{verificationPercentage}%</p>
                <p className="certificate-status-label">Verification Rate</p>
              </div>
              <div className="certificate-status-items">
                <div className="certificate-status-item warning">
                  <div className="certificate-status-item-label warning"><Clock /> Pending</div>
                  <span className="certificate-status-item-value">{pending}</span>
                </div>
                <div className="certificate-status-item success">
                  <div className="certificate-status-item-label success"><CheckCircle2 /> Verified</div>
                  <span className="certificate-status-item-value">{verified}</span>
                </div>
                <div className="certificate-status-item destructive">
                  <div className="certificate-status-item-label destructive"><XCircle /> Rejected</div>
                  <span className="certificate-status-item-value">{rejected}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="section-card">
          <h2 className="section-card-title mb-6">Recent Activity</h2>
          <p className="text-muted-foreground">Recent activity feed will be powered by real data in a future update.</p>
        </div>
      </div>
    </DashboardLayout>
  );
}
