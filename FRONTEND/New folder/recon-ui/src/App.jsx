import { Box, CircularProgress } from '@mui/material';
import { Navigate, Outlet, Routes, Route, useLocation } from 'react-router-dom';
import { useAuth } from './auth/AuthContext';
import AppLayout from './layout/AppLayout';
import Dashboard from './pages/Dashboard';
import Documents from './pages/Documents';
import ReviewQueue from './pages/ReviewQueue';
import MatchDetail from './pages/MatchDetail';
import Ledger from './pages/Ledger';
import Users from './pages/Users';
import Profile from './pages/Profile';
import Login from './pages/Login';

function RequireAuth() {
  const { user, isLoading } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <Box sx={{ minHeight: '100vh', display: 'grid', placeItems: 'center' }}>
        <CircularProgress aria-label="Checking sign-in" />
      </Box>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  return <Outlet />;
}

export default function App() {
  return (
    <Routes>
      <Route path="login" element={<Login />} />
      <Route element={<RequireAuth />}>
        <Route element={<AppLayout />}>
          <Route index element={<Dashboard />} />
          <Route path="documents" element={<Documents />} />
          <Route path="review" element={<ReviewQueue />} />
          <Route path="review/:id" element={<MatchDetail />} />
          <Route path="ledger" element={<Ledger />} />
          <Route path="users" element={<Users />} />
          <Route path="profile" element={<Profile />} />
        </Route>
      </Route>
    </Routes>
  );
}
