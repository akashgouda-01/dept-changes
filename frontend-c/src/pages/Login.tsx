import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { useToast } from '@/contexts/ToastContext';
import { Shield, Loader2 } from 'lucide-react';

export default function Login() {
  const { signInWithGoogle, isAuthenticated, user, isLoading } = useAuth();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [isSigningIn, setIsSigningIn] = useState(false);

  useEffect(() => {
    if (!isLoading && isAuthenticated && user) {
      if (user.role) {
        navigate(`/${user.role}/dashboard`);
      } else {
        // Fallback if role is missing but authenticated
        toast({ title: 'Access Denied', description: 'Your account does not have an assigned role.', variant: 'destructive' });
      }
    }
  }, [isAuthenticated, user, isLoading, navigate]);

  const handleGoogleSignIn = async () => {
    setIsSigningIn(true);
    const { error } = await signInWithGoogle();
    if (error) {
      console.error(error);
      toast({ title: 'Sign In Failed', description: 'Could not sign in with Google.', variant: 'destructive' });
      setIsSigningIn(false);
    }
    // onAuthStateChanged in context handles the rest
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

        <div className="login-branding-footer">© 2026 EduVault. All rights reserved.</div>
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
            <p>Sign in with your institutional email to access the portal</p>
          </div>

          <div className="login-form">
            <button
              type="button"
              onClick={handleGoogleSignIn}
              disabled={isSigningIn}
              className="btn btn-primary login-submit-btn"
              style={{ height: '3.5rem', fontSize: '1rem', gap: '0.75rem' }}
            >
              {isSigningIn ? (
                <Loader2 className="spinner" />
              ) : (
                <svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
                  <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
                  <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
                  <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
                  <path d="M1 1h22v22H1z" fill="none" />
                </svg>
              )}
              {isSigningIn ? 'Signing In...' : 'Sign in with Google'}
            </button>

            <p className="login-email-hint" style={{ textAlign: 'center', marginTop: '1.5rem' }}>
              Please use your @citchennai.net email address.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
