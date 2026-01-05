import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { useToast } from '@/contexts/ToastContext';
import { UserRole } from '@/types';
import { Shield, GraduationCap, Users, Loader2, Mail } from 'lucide-react';

export default function Login() {
  const [selectedRole, setSelectedRole] = useState<UserRole | null>(null);
  const [email, setEmail] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { signIn, isAuthenticated, user, isLoading } = useAuth();
  const navigate = useNavigate();
  const { toast } = useToast();

  useEffect(() => {
    if (!isLoading && isAuthenticated && user) {
      navigate(`/${user.role}/dashboard`);
    }
  }, [isAuthenticated, user, isLoading, navigate]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!selectedRole) {
      toast({ title: 'Select Role', description: 'Please select your role to continue.', variant: 'destructive' });
      return;
    }

    if (!email.endsWith('@citchennai.net')) {
      toast({ title: 'Invalid Email', description: 'Only @citchennai.net email addresses are allowed.', variant: 'destructive' });
      return;
    }

    setIsSubmitting(true);
    const { error } = signIn(email, selectedRole);
    
    if (error) {
      toast({ title: 'Sign In Failed', description: error.message, variant: 'destructive' });
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return (
      <div className="login-loading">
        <Loader2 className="spinner" />
      </div>
    );
  }

  return (
    <div className="login-page">
      <div className="login-branding">
        <div className="login-branding-shapes">
          <div className="login-branding-shape login-branding-shape-1" />
          <div className="login-branding-shape login-branding-shape-2" />
          <div className="login-branding-shape login-branding-shape-3" />
        </div>
        
        <div className="login-branding-header">
          <div className="login-branding-logo">
            <div className="login-branding-logo-icon"><Shield /></div>
            <div>
              <div className="login-branding-title">EduVault</div>
              <div className="login-branding-subtitle">Certificate Verification System</div>
            </div>
          </div>
        </div>

        <div className="login-branding-content">
          <h2 className="login-branding-heading">
            Secure Academic<br /><span>Certificate Management</span>
          </h2>
          <p className="login-branding-description">
            Streamlined verification workflow with ML-powered duplicate detection for the CSE Department.
          </p>
          <div className="login-branding-badges">
            <div className="login-branding-badge">
              <div className="login-branding-badge-label">Department</div>
              <div className="login-branding-badge-value">CSE</div>
            </div>
            <div className="login-branding-badge">
              <div className="login-branding-badge-label">Institution</div>
              <div className="login-branding-badge-value">CIT Chennai</div>
            </div>
          </div>
        </div>

        <div className="login-branding-footer">© 2024 EduVault. All rights reserved.</div>
      </div>

      <div className="login-form-container">
        <div className="login-form-wrapper">
          <div className="login-mobile-logo">
            <div className="login-mobile-logo-inner">
              <div className="login-mobile-logo-icon"><Shield /></div>
              <h1 className="login-mobile-logo-title">EduVault</h1>
            </div>
          </div>

          <div className="login-form-header">
            <h2>Welcome Back</h2>
            <p>Sign in to access the certificate verification portal</p>
          </div>

          <form onSubmit={handleSubmit} className="login-form">
            <div className="login-role-section">
              <label className="label">Select Your Role</label>
              <div className="login-role-grid">
                <button type="button" onClick={() => setSelectedRole('faculty')} className={`login-role-button ${selectedRole === 'faculty' ? 'selected' : ''}`}>
                  <GraduationCap />
                  <span>Faculty</span>
                </button>
                <button type="button" onClick={() => setSelectedRole('hod')} className={`login-role-button ${selectedRole === 'hod' ? 'selected' : ''}`}>
                  <Users />
                  <span>HOD</span>
                </button>
              </div>
            </div>

            <div className="login-email-section">
              <label className="label" htmlFor="email">Email</label>
              <div className="login-email-input-wrapper">
                <Mail />
                <input id="email" type="email" placeholder="name@citchennai.net" value={email} onChange={(e) => setEmail(e.target.value)} className="input input-with-icon" required />
              </div>
            </div>

            <button type="submit" disabled={isSubmitting} className="btn btn-primary login-submit-btn">
              {isSubmitting && <Loader2 className="spinner" style={{ width: '1.25rem', height: '1.25rem' }} />}
              Sign In
            </button>

            <p className="login-email-hint">Only @citchennai.net emails are allowed</p>
          </form>
        </div>
      </div>
    </div>
  );
}
