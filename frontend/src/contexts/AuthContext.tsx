import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import { authService } from '../services/authService';
import { setAuthToken } from '../services/api';
import { Credentials } from '../types/auth';
import { jwtDecode } from 'jwt-decode';

interface AuthContextType {
  token: string | null;
  user: { username: string } | null;
  login: (credentials: Credentials) => Promise<void>;
  register: (credentials: Credentials) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));
  const [user, setUser] = useState<{ username: string } | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (token) {
      try {
        const decoded: { username: string; exp: number } = jwtDecode(token);
        if (decoded.exp * 1000 > Date.now()) {
          setUser({ username: decoded.username });
          setAuthToken(token);
        } else {
          // Token expired
          localStorage.removeItem('token');
          setToken(null);
          setUser(null);
          setAuthToken(null);
        }
      } catch (error) {
        console.error("Invalid token:", error);
        localStorage.removeItem('token');
        setToken(null);
        setUser(null);
        setAuthToken(null);
      }
    }
    setIsLoading(false);
  }, [token]);

  const login = async (credentials: Credentials) => {
    const { token: newToken } = await authService.login(credentials);
    localStorage.setItem('token', newToken);
    setAuthToken(newToken);
    setToken(newToken);
  };

  const register = async (credentials: Credentials) => {
    await authService.register(credentials);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
    setAuthToken(null);
  };

  return (
    <AuthContext.Provider value={{ token, user, login, register, logout, isLoading }}>
      {!isLoading && children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};