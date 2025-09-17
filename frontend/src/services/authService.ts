import api from './api';
import { Credentials } from '../types/auth';

export const authService = {
  async login(credentials: Credentials): Promise<{ token: string }> {
    const response = await api.post('/auth/login', credentials);
    return response.data;
  },

  async register(credentials: Credentials): Promise<void> {
    await api.post('/auth/register', credentials);
  },
};