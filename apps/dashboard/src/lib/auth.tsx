'use client';
import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { apiRequest, getToken } from './api';

interface User {
  id: string;
  email: string;
  username: string;
  company_name: string;
  avatar_url: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  hasToken: boolean;
  login: (email: string, password: string, field?: 'email' | 'phone') => Promise<void>;
  register: (data: { email: string; password: string; username: string; company_name?: string }) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [hasToken, setHasToken] = useState(false);

  useEffect(() => {
    const token = getToken();
    if (token) {
      setHasToken(true);
      // Retry once on transient error before giving up
      let retried = false;
      const attempt = () => {
        apiRequest<{ user: User }>('/users/me')
          .then((res) => setUser(res.data.user))
          .catch((err) => {
            if (err instanceof Error && (err.message.includes('401') || err.message.includes('403') || err.message.includes('unauthorized'))) {
              localStorage.removeItem('token');
              setHasToken(false);
            } else if (!retried) {
              retried = true;
              // Retry after 1 second delay for transient network issues
              setTimeout(attempt, 1000);
              return; // Don't call finally yet
            }
            // Non-auth error after retry: keep token, show page without user data
          })
          .finally(() => setLoading(false));
      };
      attempt();
    } else {
      setHasToken(false);
      setLoading(false);
    }
  }, []);

  const login = async (email: string, password: string, field: 'email' | 'phone' = 'email') => {
    const body: Record<string, string> = { password };
    if (field === 'phone') {
      body.phone = email;
    } else {
      body.email = email;
    }
    const res = await apiRequest<{ token: string; user: User }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(body),
    });
    localStorage.setItem('token', res.data.token);
    setUser(res.data.user);
    setHasToken(true);
  };

  const register = async (data: { email: string; password: string; username: string; company_name?: string }) => {
    const res = await apiRequest<{ token: string; user: User }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    localStorage.setItem('token', res.data.token);
    setUser(res.data.user);
    setHasToken(true);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
    setHasToken(false);
  };

  return (
    <AuthContext.Provider value={{ user, loading, hasToken, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
