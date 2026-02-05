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
    // Check local storage for persisted user
    const storedUser = localStorage.getItem('eduvault_user');
    if (storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        // basic validation
        if (parsedUser && parsedUser.email && parsedUser.role) {
          setUser(parsedUser);
        } else {
          localStorage.removeItem('eduvault_user');
        }
      } catch (error) {
        console.error('Failed to parse user from local storage:', error);
        localStorage.removeItem('eduvault_user');
      }
    }
    setIsLoading(false);
  }, []);

  const signIn = (email: string, role: UserRole): { error: Error | null } => {
    try {
      // Whitelist Validation
      const accessMap: Record<string, { role: UserRole; staffId: string; sections?: string[] }> = {
        'harjeetp.cse2024@citchennai.net': {
          role: 'faculty',
          staffId: 'FAC01',
          sections: ['A', 'B', 'C', 'D', 'E', 'F']
        },
        'hemanm.cse2024@citchennai.net': {
          role: 'faculty',
          staffId: 'FAC02',
          sections: ['G', 'H', 'I', 'J', 'K', 'L']
        },
        'akashkumargouda.cse2024@citchennai.net': {
          role: 'faculty',
          staffId: 'FAC03',
          sections: ['M', 'N', 'O', 'P', 'Q']
        },
        'aadhishs.cse2024@citchennai.net': {
          role: 'hod',
          staffId: 'HOD01'
        }
      };

      const normalizedEmail = email.toLowerCase();
      const userConfig = accessMap[normalizedEmail];

      if (!userConfig) {
        return { error: new Error('This email is not authorized.') };
      }

      if (userConfig.role !== role) {
        return { error: new Error(`This email belongs to ${userConfig.role}, not ${role}.`) };
      }

      const newUser: AuthUser = {
        id: crypto.randomUUID(),
        name: email.split('@')[0],
        email: email,
        role: role,
        staffId: userConfig.staffId,
        position: role === 'hod' ? 'Head of Department' : 'Assistant Professor',
        department: 'Computer Science & Engineering',
        assignedSections: userConfig.sections
      };

      setUser(newUser);
      localStorage.setItem('eduvault_user', JSON.stringify(newUser));
      return { error: null };
    } catch (error: any) {
      return { error };
    }
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('eduvault_user');
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
