import { Link as RouterLink } from 'react-router-dom';
import { Box, Button, Link, Paper, Stack, Typography } from '@mui/material';
import AuthLayout, { AuthField } from '../layout/AuthLayout';

export default function Register() {
  const handleSubmit = (event) => {
    event.preventDefault();
  };

  return (
    <AuthLayout isRegister>
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
                Create account
              </Typography>
              <Typography sx={{ fontSize: { xs: 14, md: 10, lg: 14 }, lineHeight: { xs: '20px', md: '14px', lg: '20px' }, color: '#637381' }}>
                Set up your workspace account to get started.
              </Typography>
            </Box>

            <Stack spacing={{ xs: 2.5, md: 1.75, lg: 2.5 }}>
              <AuthField label="Work email" placeholder="name@company.lk" type="email" autoComplete="email" />
              <AuthField label="Password" placeholder="Enter your password" type="password" autoComplete="new-password" />
              <AuthField label="Full name" placeholder="Enter your full name" autoComplete="name" />
              <AuthField label="Confirm password" placeholder="Re-enter your password" type="password" autoComplete="new-password" />
            </Stack>

            <Button
              type="submit"
              variant="contained"
              fullWidth
              sx={{ minHeight: { xs: 46, md: 33, lg: 46 }, bgcolor: '#2563eb', borderRadius: 1.5, fontSize: { xs: 14, md: 10, lg: 14 }, '&:hover': { bgcolor: '#1d4ed8' } }}
            >
              Create account
            </Button>
          </Stack>
        </Paper>
        <Typography sx={{ textAlign: 'center', fontSize: { xs: 12, md: 9, lg: 12 }, color: '#637381', lineHeight: 1.6 }}>
          By creating an account, you agree to your workspace terms and security policy.
        </Typography>
        <Typography sx={{ textAlign: 'center', fontSize: { xs: 12, md: 9, lg: 12 }, color: '#637381' }}>
          Already have an account?{' '}
          <Link component={RouterLink} to="/login" underline="none" sx={{ color: '#2563eb', fontWeight: 500 }}>
            Sign in
          </Link>
        </Typography>
      </Stack>
    </AuthLayout>
  );
}