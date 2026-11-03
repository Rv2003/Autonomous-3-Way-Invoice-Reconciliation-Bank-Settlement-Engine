import { Routes, Route } from 'react-router-dom';
import AppLayout from './layout/AppLayout';
import Dashboard from './pages/Dashboard';
import Documents from './pages/Documents';
import ReviewQueue from './pages/ReviewQueue';
import MatchDetail from './pages/MatchDetail';
import Ledger from './pages/Ledger';
import Users from './pages/Users';
import Login from './pages/Login';
import Register from './pages/Register';

export default function App() {
  return (
    <Routes>
      <Route path="login" element={<Login />} />
      <Route path="register" element={<Register />} />
      <Route element={<AppLayout />}>
        <Route index element={<Dashboard />} />
        <Route path="documents" element={<Documents />} />
        <Route path="review" element={<ReviewQueue />} />
        <Route path="review/:id" element={<MatchDetail />} />
        <Route path="ledger" element={<Ledger />} />
        <Route path="users" element={<Users />} />
      </Route>
    </Routes>
  );
}
