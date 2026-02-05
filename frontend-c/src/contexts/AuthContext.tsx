import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { UserRole } from '@/types';
import { auth, googleProvider } from '@/config/firebase';
import { signInWithPopup, signOut, onAuthStateChanged, User as FirebaseUser } from 'firebase/auth';

interface AuthUser {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  staffId?: string;
  position?: string;
  department?: string;
  assignedSections?: string[];
  token?: string;
}

interface AuthContextType {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  signInWithGoogle: () => Promise<{ error: Error | null }>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (currentUser) => {
      setIsLoading(true);
      if (currentUser) {
        try {
          const token = await currentUser.getIdToken();
          const email = currentUser.email || '';
          const derivedUser = deriveUserFromEmail(email, currentUser.uid, token);

          if (derivedUser) {
            setUser(derivedUser);
            localStorage.setItem('eduvault_token', token);
          } else {
            // Invalid domain or role
            setUser(null);
            await signOut(auth);
          }
        } catch (error) {
          console.error("Auth state change error", error);
          setUser(null);
        }
      } else {
        setUser(null);
        localStorage.removeItem('eduvault_token');
      }
      setIsLoading(false);
    });

    return () => unsubscribe();
  }, []);

  const deriveUserFromEmail = (email: string, uid: string, token: string): AuthUser | null => {
    if (!email.endsWith('@citchennai.net')) {
      return null;
    }

    const prefix = email.split('@')[0].toLowerCase();
    let role: UserRole = 'faculty'; // Default or need distinct logic
    let staffId = 'FAC01';
    let assignedSections: string[] | undefined;

    if (prefix === 'hod') {
      role = 'hod';
      staffId = 'HOD01';
    } else if (prefix.startsWith('faculty')) {
      role = 'faculty';
      if (prefix === 'faculty1') {
        assignedSections = ['A', 'B', 'C', 'D', 'E', 'F'];
        staffId = 'FAC01';
      } else if (prefix === 'faculty2') {
        assignedSections = ['G', 'H', 'I', 'J', 'K', 'L'];
        staffId = 'FAC02';
      } else if (prefix === 'faculty3') {
        assignedSections = ['M', 'N', 'O', 'P', 'Q'];
        staffId = 'FAC03';
      } else {
        assignedSections = ['A', 'B']; // Fallback
        staffId = 'FAC00';
      }
    } else {
      // Assume student or unknown
      // For now, if we only strictly allowing faculty/hod:
      // return null; 
      // But let's allow basic access if needed or just return null
      return null;
    }

    return {
      id: uid,
      name: currentUserDisplayName(auth.currentUser) || email.split('@')[0],
      email: email,
      role: role,
      staffId: staffId,
      position: role === 'hod' ? 'Head of Department' : 'Assistant Professor',
      department: 'Computer Science & Engineering',
      assignedSections,
      token: token
    };
  }

  const currentUserDisplayName = (u: FirebaseUser | null) => {
    return u?.displayName || '';
  }

  const signInWithGoogle = async (): Promise<{ error: Error | null }> => {
    try {
      await signInWithPopup(auth, googleProvider);
      return { error: null };
    } catch (error: any) {
      return { error };
    }
  };

  const logout = async () => {
    await signOut(auth);
    setUser(null);
    localStorage.removeItem('eduvault_token');
  };

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, signInWithGoogle, logout }}>
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
