import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export const ProtectedRoute = () => {
  const { token, isLoading } = useAuth();

  if (isLoading) {
    return <div>Loading...</div>; // Or a spinner component
  }

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
};