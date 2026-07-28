'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { apiRequest, getToken } from './api';

interface Admin {
  id: string;
  email: string;
  username: string;
  role: string;
  status: string;
}

interface AuthContextType {
  admin: Admin | null;
  loading: boolean;
  hasToken: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [admin, setAdmin] = useState<Admin | null>(null);
  const [loading, setLoading] = useState(true);
  const [hasToken, setHasToken] = useState(false);

  useEffect(() => {
    const token = getToken();
    if (token) {
      setHasToken(true);
      let retried = false;
      const controller = new AbortController();
      const timeoutId = window.setTimeout(() => controller.abort(), 8000);

      const attempt = () => {
        apiRequest<Admin>('/auth/me', { signal: controller.signal })
          .then((res) => setAdmin(res.data))
          .catch((err) => {
            if (err instanceof Error && (err.message.includes('401') || err.message.includes('403') || err.message.includes('unauthorized'))) {
              localStorage.removeItem('admin_token');
              setHasToken(false);
            } else if (!retried && !controller.signal.aborted) {
              retried = true;
              window.clearTimeout(timeoutId);
              setTimeout(() => attempt(), 1000);
              return;
            }
          })
          .finally(() => {
            window.clearTimeout(timeoutId);
            setLoading(false);
          });
      };
      attempt();
    } else {
      setHasToken(false);
      setLoading(false);
    }
  }, []);

  const login = async (email: string, password: string) => {
    const res = await apiRequest<{ token: string; admin: Admin }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    localStorage.setItem('admin_token', res.data.token);
    setAdmin(res.data.admin);
    setHasToken(true);
  };

  const logout = () => {
    localStorage.removeItem('admin_token');
    setAdmin(null);
    setHasToken(false);
  };

  return (
    <AuthContext.Provider value={{ admin, loading, hasToken, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
