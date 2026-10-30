import { Navigate, Route, Routes } from 'react-router-dom';
import { GuestOnly, RequireAuth, RequireRole } from './auth/RouteGuards';
import AppLayout from './layout/AppLayout';
import Dashboard from './pages/Dashboard';
import Documents from './pages/Documents';
import ReviewQueue from './pages/ReviewQueue';
import MatchDetail from './pages/MatchDetail';
import Ledger from './pages/Ledger';
import Users from './pages/Users';
import Login from './pages/Login';

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<GuestOnly><Login /></GuestOnly>} />
      <Route element={<RequireAuth><AppLayout /></RequireAuth>}>
        <Route index element={<Dashboard />} />
        <Route path="documents" element={<Documents />} />
        <Route path="review" element={<ReviewQueue />} />
        <Route path="review/:id" element={<MatchDetail />} />
        <Route path="ledger" element={<Ledger />} />
        <Route
          path="users"
          element={<RequireRole roles={['administrator']}><Users /></RequireRole>}
        />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
