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
  signIn: (email: string, password: string, role: UserRole) => Promise<{ error: Error | null }>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);
import { apiClient } from '@/api/client';

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check local storage for persisted user
    const storedUser = localStorage.getItem('eduvault_user');
    const storedToken = localStorage.getItem('eduvault_token');

    if (storedUser && storedToken) {
      try {
        const parsedUser = JSON.parse(storedUser);
        if (parsedUser && parsedUser.email && parsedUser.role) {
          setUser(parsedUser);
        } else {
          localStorage.removeItem('eduvault_user');
          localStorage.removeItem('eduvault_token');
        }
      } catch (error) {
        console.error('Failed to parse user from local storage:', error);
        localStorage.removeItem('eduvault_user');
        localStorage.removeItem('eduvault_token');
      }
    }
    setIsLoading(false);
  }, []);

  const signIn = async (email: string, password: string, role: UserRole): Promise<{ error: Error | null }> => {
    try {
      const response = await apiClient.post('/auth/login', {
        email,
        password,
        role
      });

      const { token, user: apiUser } = response.data;

      // Map API user to AuthUser context
      const newUser: AuthUser = {
        id: apiUser.faculty_id || 'HOD01',  // Use faculty_id from backend or fallback
        name: apiUser.email.split('@')[0],
        email: apiUser.email,
        role: apiUser.role,
        staffId: apiUser.faculty_id,
        position: apiUser.role === 'hod' ? 'Head of Department' : 'Assistant Professor',
        department: 'Computer Science & Engineering',
        assignedSections: [], // You might want the backend to return this too, or map it here if needed
      };

      // Temporary Hardcoded Section Map until backend sends it
      // Map Faculty Details based on ID
      if (newUser.id === 'CSE245') {
        newUser.name = 'R. Poornima Lakshmi';
        newUser.assignedSections = ['L', 'M', 'N', 'O', 'P', 'Q'];
      }
      if (newUser.id === 'CSE086') {
        newUser.name = 'Selvajothi M';
        newUser.assignedSections = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q'];
      }
      if (newUser.id === 'CSE262') {
        newUser.name = 'Nishanthini S';
        newUser.assignedSections = ['C', 'D', 'E', 'F', 'G', 'H', 'I'];
      }
      if (newUser.id === 'CSE345') {
        newUser.name = 'M. Ayeesha Nasreen';
        newUser.assignedSections = ['A', 'B', 'I', 'J', 'K'];
      }


      setUser(newUser);
      localStorage.setItem('eduvault_user', JSON.stringify(newUser));
      localStorage.setItem('eduvault_token', token);

      return { error: null };
    } catch (error: any) {
      console.error("Login error", error);
      const msg = error.response?.data?.error || "Login failed";
      return { error: new Error(msg) };
    }
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('eduvault_user');
    localStorage.removeItem('eduvault_token');
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
