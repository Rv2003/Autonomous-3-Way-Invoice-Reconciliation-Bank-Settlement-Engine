import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  Alert, Box, Button, Checkbox, FormControlLabel, Paper, Stack, Typography,
} from '@mui/material';
import errorIcon from '../assets/auth/error.svg';
import checkIcon from '../assets/auth/check.svg';
import useAuth from '../auth/useAuth';
import AuthLayout, { AuthField } from '../layout/AuthLayout';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const { login, sessionError } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setSubmitting(true);
    try {
      await login(email.trim(), password, rememberMe);
      navigate(location.state?.from?.pathname || '/', { replace: true });
    } catch (loginError) {
      setError(loginError.message || 'Sign in failed. Please try again.');
    } finally {
      setSubmitting(false);
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

            {(error || sessionError) && (
              <Alert
                severity="error"
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
                {error || `Could not verify your saved session: ${sessionError}`}
              </Alert>
            )}

            <Stack spacing={{ xs: 2.5, md: 1.75, lg: 2.5 }}>
              <AuthField
                label="Email"
                name="email"
                placeholder="name@company.lk"
                type="email"
                autoComplete="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                required
                error={Boolean(error)}
              />
              <AuthField
                label="Password"
                name="password"
                placeholder="Enter your password"
                type="password"
                autoComplete="current-password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                required
                error={Boolean(error)}
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
              <Typography sx={{ fontSize: { xs: 12, md: 9, lg: 12 }, color: '#637381', textAlign: 'right' }}>
                Contact your administrator to reset your password.
              </Typography>
            </Stack>

            <Button
              type="submit"
              variant="contained"
              fullWidth
              disabled={submitting}
              sx={{ minHeight: { xs: 46, md: 33, lg: 46 }, bgcolor: '#2563eb', borderRadius: 1.5, fontSize: { xs: 14, md: 10, lg: 14 }, '&:hover': { bgcolor: '#1d4ed8' } }}
            >
              {submitting ? 'Signing in…' : 'Sign in'}
            </Button>
          </Stack>
        </Paper>
        <Typography sx={{ textAlign: 'center', fontSize: { xs: 12, md: 9, lg: 12 }, color: '#637381', lineHeight: 1.6 }}>
          Need access? Ask your workspace administrator to create your account.
        </Typography>
      </Stack>
    </AuthLayout>
  );
}