import { useEffect, useState } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { useAuth } from '@/contexts/AuthContext';
import { Users, FileCheck, Clock, CheckCircle2, XCircle, TrendingUp, ChevronDown, ChevronUp } from 'lucide-react';
import {
  getDashboardOverview,
  getDashboardSections,
  getRecentActivity,
  DashboardOverview,
  SectionStat,
  RecentActivity
} from '@/api';

export default function FacultyDashboard() {
  const { user } = useAuth();
  const [overview, setOverview] = useState<DashboardOverview | null>(null);
  const [sections, setSections] = useState<SectionStat[]>([]);
  const [activities, setActivities] = useState<RecentActivity[]>([]);
  const [activityFilter, setActivityFilter] = useState<string>('ALL');
  const [showAllActivities, setShowAllActivities] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadDashboard = async () => {
      try {
        setIsLoading(true);
        setError(null);


        const [overviewRes, sectionsRes, activityRes] = await Promise.all([
          getDashboardOverview(),
          getDashboardSections(),
          getRecentActivity(),
        ]);

        if (isMounted) {
          if (overviewRes.success) {
            setOverview(overviewRes.data as DashboardOverview);
          }
          if (sectionsRes.success) {
            setSections(Array.isArray(sectionsRes.data) ? (sectionsRes.data as SectionStat[]) : []);
          }
          if (activityRes.success) {
            setActivities(Array.isArray(activityRes.data) ? (activityRes.data as RecentActivity[]) : []);
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

  const handleExport = async () => {
    try {
      const response = await import('@/api/dashboard.api').then(m => m.exportMyCertificates());

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;

      const contentDisp = response.headers['content-disposition'];
      let customFilename = `My_Certificates.xlsx`;
      if (contentDisp && contentDisp.indexOf('filename=') !== -1) {
        customFilename = contentDisp.split('filename=')[1].replace(/"/g, '');
      }

      link.setAttribute('download', customFilename);
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (error) {
      console.error('Failed to export:', error);
      setError('Failed to export certificates.');
    }
  };

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
          <div className="flex gap-4 items-center">
            <button onClick={handleExport} className="btn btn-outline btn-sm gap-2">
              <TrendingUp className="w-4 h-4" /> Export Excel
            </button>
            <div className="dashboard-date">
              <Clock />
              <span>{new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</span>
            </div>
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
          <div className="flex items-center justify-between mb-6">
            <h2 className="section-card-title">Recent Activity (Today)</h2>
            <select
              value={activityFilter}
              onChange={(e) => setActivityFilter(e.target.value)}
            >
              <option value="ALL">All Actions</option>
              <option value="UPLOADED">Uploaded</option>
              <option value="VERIFIED">Verified</option>
              <option value="REJECTED">Rejected</option>
            </select>
          </div>

          <div className="space-y-3">
            {(() => {
              const filteredActivities = activities.filter(a => activityFilter === 'ALL' || a.action === activityFilter);
              const visibleActivities = showAllActivities ? filteredActivities : filteredActivities.slice(0, 5);
              const hasMore = filteredActivities.length > 5;

              if (filteredActivities.length === 0) {
                return <p className="text-muted-foreground text-sm text-center py-6 bg-muted/20 rounded-lg border border-dashed">No recent activity found for today.</p>;
              }

              return (
                <>
                  <div className="space-y-3">
                    {visibleActivities.map((activity) => (
                      <div
                        key={activity.id}
                        className="activity-item group"
                      >
                        {/* Status Strip */}
                        <div className={`activity-strip ${activity.action === 'VERIFIED' ? 'strip-verified' :
                          activity.action === 'REJECTED' ? 'strip-rejected' :
                            'strip-uploaded'
                          }`} />

                        <div className="flex items-center gap-4 pl-2">
                          <div className={`activity-icon ${activity.action === 'VERIFIED' ? 'bg-emerald-50 text-emerald-600 border-emerald-200' :
                            activity.action === 'REJECTED' ? 'bg-rose-50 text-rose-600 border-rose-200' :
                              'bg-blue-50 text-blue-600 border-blue-200'
                            }`}>
                            {activity.action === 'VERIFIED' ? <CheckCircle2 size={20} strokeWidth={2.5} className="text-emerald-600" /> :
                              activity.action === 'REJECTED' ? <XCircle size={20} strokeWidth={2.5} className="text-rose-600" /> :
                                <FileCheck size={20} strokeWidth={2.5} className="text-blue-600" />}
                          </div>

                          <div className="space-y-1">
                            <div className="flex items-center gap-2">
                              <p className="activity-name group-hover:text-primary transition-colors">
                                {activity.student_name}
                              </p>
                              <span className={`activity-tag ${activity.action === 'VERIFIED' ? 'tag-verified' :
                                activity.action === 'REJECTED' ? 'tag-rejected' :
                                  'tag-uploaded'
                                }`}>
                                {activity.action}
                              </span>
                            </div>

                            <div className="activity-details">
                              <span className="font-mono text-xs">{activity.reg_no}</span>
                              <div className="w-1 h-1 rounded-full bg-border" />
                              <span>Section {activity.section}</span>
                            </div>
                          </div>
                        </div>

                        <div className="text-right pl-4">
                          <div className="activity-time inline-flex items-center gap-1.5">
                            <Clock size={12} className="opacity-70" />
                            {new Date(activity.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>

                  {hasMore && (
                    <button
                      onClick={() => setShowAllActivities(!showAllActivities)}
                      className="w-full flex items-center justify-center gap-2 p-3 mt-4 text-sm font-medium text-primary hover:text-primary/80 hover:bg-primary/5 rounded-lg border border-transparent hover:border-primary/10 transition-all duration-200 group"
                    >
                      {showAllActivities ? (
                        <>Show Less <ChevronUp size={16} className="group-hover:-translate-y-0.5 transition-transform" /></>
                      ) : (
                        <>Show More ({filteredActivities.length - 5} hidden) <ChevronDown size={16} className="group-hover:translate-y-0.5 transition-transform" /></>
                      )}
                    </button>
                  )}
                </>
              );
            })()}
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
