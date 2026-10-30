import { useCallback, useEffect, useMemo, useState } from 'react';
import { apiRequest } from '../api/client';
import AuthContext from './auth-context';

const TOKEN_KEY = 'invoice-reconciliation-token';
const supportedRoles = new Set(['administrator', 'clerk']);

function validateUser(user) {
  if (
    !user
    || typeof user.id !== 'string'
    || typeof user.full_name !== 'string'
    || typeof user.email !== 'string'
    || !supportedRoles.has(user.role)
  ) {
    throw new Error('The server returned an invalid user profile.');
  }

  return user;
}

function readStoredToken() {
  return sessionStorage.getItem(TOKEN_KEY) || localStorage.getItem(TOKEN_KEY);
}

function clearStoredToken() {
  sessionStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(TOKEN_KEY);
}

export function AuthProvider({ children }) {
  const [initialToken] = useState(readStoredToken);
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(initialToken);
  const [loading, setLoading] = useState(() => Boolean(initialToken));
  const [sessionError, setSessionError] = useState('');

  useEffect(() => {
    const storedToken = initialToken;
    if (!storedToken) {
      return undefined;
    }

    let active = true;
    apiRequest('/auth/me', { token: storedToken })
      .then((profile) => {
        if (active) {
          setToken(storedToken);
          setUser(validateUser(profile));
        }
      })
      .catch((error) => {
        if (!active) return;
        if (error.status === 401) {
          clearStoredToken();
          setToken(null);
        } else {
          setSessionError(error.message || 'Could not verify the saved session.');
        }
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, [initialToken]);

  const login = useCallback(async (email, password, rememberMe) => {
    const result = await apiRequest('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });

    if (!result || typeof result.token !== 'string' || !result.token) {
      throw new Error('The server did not return an authentication token.');
    }

    const profile = validateUser(result.user);
    clearStoredToken();
    (rememberMe ? localStorage : sessionStorage).setItem(TOKEN_KEY, result.token);
    setToken(result.token);
    setUser(profile);
    setSessionError('');
  }, []);

  const logout = useCallback(() => {
    clearStoredToken();
    setToken(null);
    setUser(null);
    setSessionError('');
  }, []);

  const value = useMemo(
    () => ({ user, token, loading, sessionError, login, logout }),
    [user, token, loading, sessionError, login, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
