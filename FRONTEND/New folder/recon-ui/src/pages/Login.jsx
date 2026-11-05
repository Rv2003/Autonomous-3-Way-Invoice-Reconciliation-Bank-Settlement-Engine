import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Alert, Box, Button, Checkbox, FormControlLabel, Link, Paper, Stack, Typography,
} from '@mui/material';
import errorIcon from '../assets/auth/error.svg';
import checkIcon from '../assets/auth/check.svg';
import AuthLayout, { AuthField } from '../layout/AuthLayout';
import { API_BASE_URL } from '../api/documents';
import { useAuth } from '../auth/AuthContext';

export default function Login() {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [error, setError] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setIsSubmitting(true);

    try {
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          email: email.trim(),
          password,
          remember_me: rememberMe,
        }),
      });

      if (!response.ok) {
        const detail = (await response.text()).trim();
        setError(response.status === 401
          ? 'Incorrect email or password. Check your details and try again.'
          : detail || `Sign in failed (${response.status}). Please try again.`);
        return;
      }

      await login();
      navigate('/');
    } catch (requestError) {
      setError(requestError instanceof TypeError
        ? `Cannot reach the API at ${API_BASE_URL}. Start the backend and check VITE_API_BASE_URL.`
        : requestError.message || 'Could not connect to the server. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <AuthLayout>
      <Stack spacing={{ xs: 2.5, md: 1.75, lg: 2.5 }} component="form" onSubmit={handleSubmit}>
        <Paper
          variant="outlined"
          sx={{
            width: '100%',
            p: { xs: 3, sm: 5, md: 3.5, lg: 5 },
            borderColor: '#e0e0e0',
            borderRadius: 2,
            boxShadow: '0 1px 3px rgba(0,0,0,0.08)',
          }}
        >
          <Stack spacing={{ xs: 3, md: 2.1, lg: 3 }}>
            <Box>
              <Typography component="h1" sx={{ fontSize: { xs: 28, md: 20, lg: 28 }, lineHeight: 1.2, fontWeight: 700, color: '#212121', mb: { xs: 0.75, md: 0.5, lg: 0.75 } }}>
                Sign in
              </Typography>
              <Typography sx={{ fontSize: { xs: 14, md: 10, lg: 14 }, lineHeight: { xs: '20px', md: '14px', lg: '20px' }, color: '#637381' }}>
                Use your workspace account to open the dashboard.
              </Typography>
            </Box>

            {error && (
              <Alert
                icon={<Box component="img" src={errorIcon} alt="" sx={{ width: 20, height: 20 }} />}
                sx={{
                  px: 1.75,
                  py: 1.5,
                  borderRadius: 1.5,
                  bgcolor: '#fdecea',
                  color: '#5f2120',
                  alignItems: 'flex-start',
                  '& .MuiAlert-icon': { mr: 1.25, mt: 0.25, p: 0 },
                  '& .MuiAlert-message': { p: 0, fontSize: { xs: 14, md: 10, lg: 14 }, lineHeight: { xs: '20px', md: '14px', lg: '20px' } },
                }}
              >
                {error}
              </Alert>
            )}

            <Stack spacing={{ xs: 2.5, md: 1.75, lg: 2.5 }}>
              <AuthField
                label="Email"
                placeholder="name@company.lk"
                type="email"
                autoComplete="email"
                error={Boolean(error)}
                value={email}
                onChange={(event) => setEmail(event.target.value)}
              />
              <AuthField
                label="Password"
                placeholder="Enter your password"
                type="password"
                autoComplete="current-password"
                error={Boolean(error)}
                value={password}
                onChange={(event) => setPassword(event.target.value)}
              />
            </Stack>

            <Stack direction="row" justifyContent="space-between" alignItems="center" gap={1}>
              <FormControlLabel
                control={<Checkbox checked={rememberMe} onChange={(event) => setRememberMe(event.target.checked)} size="small" icon={<Box sx={{ width: 18, height: 18, border: '1px solid #9aa5b1', borderRadius: '3px' }} />} checkedIcon={<Box component="img" src={checkIcon} alt="" sx={{ width: 18, height: 18 }} />} />}
                label="Keep me signed in"
                sx={{
                  m: 0,
                  '& .MuiFormControlLabel-label': { fontSize: { xs: 14, md: 10, lg: 14 }, color: '#212121', whiteSpace: 'nowrap' },
                  '& .MuiCheckbox-root': { p: 0, mr: 0.75 },
                }}
              />
              <Link component="button" type="button" underline="none" sx={{ fontSize: { xs: 14, md: 10, lg: 14 }, fontWeight: 500, color: '#2563eb', whiteSpace: 'nowrap' }}>
                Forgot password?
              </Link>
            </Stack>

            <Button
              type="submit"
              variant="contained"
              fullWidth
              disabled={isSubmitting}
              sx={{ minHeight: { xs: 46, md: 33, lg: 46 }, bgcolor: '#2563eb', borderRadius: 1.5, fontSize: { xs: 14, md: 10, lg: 14 }, '&:hover': { bgcolor: '#1d4ed8' } }}
            >
              {isSubmitting ? 'Signing in…' : 'Sign in'}
            </Button>
          </Stack>
        </Paper>
        <Typography sx={{ textAlign: 'center', fontSize: { xs: 12, md: 9, lg: 12 }, color: '#637381', lineHeight: 1.6 }}>
          Need access? Ask your workspace administrator to add you.
        </Typography>
      </Stack>
    </AuthLayout>
  );
}