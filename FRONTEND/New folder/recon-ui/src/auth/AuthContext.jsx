import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { API_BASE_URL } from '../components/api/documents';

const AuthContext = createContext(null);

async function fetchCurrentUser(token) {
  const response = await fetch(`${API_BASE_URL}/auth/me`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!response.ok) {
    const detail = (await response.text()).trim();
    const error = new Error(detail || `Could not load user data (${response.status}).`);
    error.status = response.status;
    throw error;
  }

  return response.json();
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    let active = true;
    const savedToken = localStorage.getItem('token') || sessionStorage.getItem('token');

    if (!savedToken) {
      setIsLoading(false);
      return () => {
        active = false;
      };
    }

    fetchCurrentUser(savedToken)
      .then((currentUser) => {
        if (!active) return;
        setToken(savedToken);
        setUser(currentUser);
      })
      .catch((requestError) => {
        if (!active) return;
        setError(requestError.message || 'Could not load the signed-in user.');
        if (requestError.status === 401) {
          localStorage.removeItem('token');
          sessionStorage.removeItem('token');
          localStorage.removeItem('user');
          sessionStorage.removeItem('user');
        }
      })
      .finally(() => {
        if (active) setIsLoading(false);
      });

    return () => {
      active = false;
    };
  }, []);

  const login = useCallback(async (newToken, rememberMe) => {
    const currentUser = await fetchCurrentUser(newToken);
    const storage = rememberMe ? localStorage : sessionStorage;
    const otherStorage = rememberMe ? sessionStorage : localStorage;

    storage.setItem('token', newToken);
    storage.setItem('user', JSON.stringify(currentUser));
    otherStorage.removeItem('token');
    otherStorage.removeItem('user');
    setToken(newToken);
    setUser(currentUser);
    setError('');

    return currentUser;
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    sessionStorage.removeItem('token');
    sessionStorage.removeItem('user');
    setToken(null);
    setUser(null);
    setError('');
  }, []);

  const value = useMemo(
    () => ({ user, token, isLoading, error, login, logout }),
    [user, token, isLoading, error, login, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider.');
  }
  return context;
}
