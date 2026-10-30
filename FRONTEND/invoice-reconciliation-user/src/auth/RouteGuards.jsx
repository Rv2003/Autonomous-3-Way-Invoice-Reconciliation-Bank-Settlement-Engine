import { Box, CircularProgress } from '@mui/material';
import { Navigate, useLocation } from 'react-router-dom';
import useAuth from './useAuth';

function LoadingScreen() {
  return (
    <Box sx={{ minHeight: '100vh', display: 'grid', placeItems: 'center' }}>
      <CircularProgress aria-label="Checking your session" />
    </Box>
  );
}

export function RequireAuth({ children }) {
  const { loading, user } = useAuth();
  const location = useLocation();

  if (loading) return <LoadingScreen />;
  if (!user) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return children;
}

export function GuestOnly({ children }) {
  const { loading, user } = useAuth();
  const location = useLocation();

  if (loading) return <LoadingScreen />;
  if (user) {
    return <Navigate to={location.state?.from?.pathname || '/'} replace />;
  }

  return children;
}

export function RequireRole({ roles, children }) {
  const { user } = useAuth();
  if (!roles.includes(user?.role)) {
    return <Navigate to="/" replace />;
  }
  return children;
}
