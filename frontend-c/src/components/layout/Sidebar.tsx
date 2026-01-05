import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { LayoutDashboard, FileCheck, User, LogOut, Shield, GraduationCap } from 'lucide-react';

export function Sidebar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  const facultyLinks = [
    { to: '/faculty/dashboard', icon: LayoutDashboard, label: 'Dashboard' },
    { to: '/faculty/certificates', icon: FileCheck, label: 'Certificate Verification' },
    { to: '/faculty/profile', icon: User, label: 'Profile' },
  ];

  const hodLinks = [
    { to: '/hod/dashboard', icon: LayoutDashboard, label: 'Dashboard' },
  ];

  const links = user?.role === 'hod' ? hodLinks : facultyLinks;

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <div className="sidebar-logo">
          <div className="sidebar-logo-icon"><Shield /></div>
          <div className="sidebar-logo-text">
            <h1>EduVault</h1>
            <p>CSE Department</p>
          </div>
        </div>
      </div>

      <nav className="sidebar-nav">
        {links.map((link) => (
          <Link key={link.to} to={link.to} className={`sidebar-nav-link ${location.pathname === link.to ? 'active' : ''}`}>
            <link.icon />
            <span>{link.label}</span>
          </Link>
        ))}
      </nav>

      <div className="sidebar-footer">
        <div className="sidebar-user">
          <div className="sidebar-user-avatar"><GraduationCap /></div>
          <div className="sidebar-user-info">
            <div className="sidebar-user-name">{user?.name}</div>
            <div className="sidebar-user-role">{user?.role}</div>
          </div>
        </div>
        <button className="sidebar-logout-btn" onClick={handleLogout}>
          <LogOut />
          Sign Out
        </button>
      </div>
    </aside>
  );
}
