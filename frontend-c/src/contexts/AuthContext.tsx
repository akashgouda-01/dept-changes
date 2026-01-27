import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { UserRole } from '@/types';

interface AuthUser {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  staffId?: string;
  position?: string;
  department?: string;
  assignedSections?: string[];
}

interface AuthContextType {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  signIn: (email: string, role: UserRole) => { error: Error | null };
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check for existing session in localStorage
    const storedUser = localStorage.getItem('eduvault_user');
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
    setIsLoading(false);
  }, []);

  const signIn = (email: string, role: UserRole): { error: Error | null } => {
    if (!email.endsWith('@citchennai.net')) {
      return { error: new Error('Only @citchennai.net email addresses are allowed') };
    }

    const authUser: AuthUser = {
      id: crypto.randomUUID(),
      name: email.split('@')[0],
      email: email,
      role: role,
      staffId: role === 'hod' ? 'HOD-01' : 'FAC-01',
      position: role === 'hod' ? 'Head of Department' : 'Assistant Professor',
      department: 'Computer Science & Engineering',
      assignedSections: role === 'faculty' ? ['A', 'B'] : undefined,
    };

    localStorage.setItem('eduvault_user', JSON.stringify(authUser));
    setUser(authUser);
    return { error: null };
  };

  const logout = () => {
    localStorage.removeItem('eduvault_user');
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, signIn, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
