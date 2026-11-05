import { useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { Alert, AppBar, Avatar, Badge, Box, Divider, IconButton, Stack, Toolbar, Typography, useMediaQuery, useTheme } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import MenuIcon from '@mui/icons-material/Menu';
import Sidebar from './Sidebar';
import { useAuth } from '../auth/AuthContext';

const titles = { '/': 'Dashboard', '/documents': 'Documents', '/review': 'Review Queue', '/ledger': 'Ledger', '/users': 'Users' };

export default function AppLayout() {
  const { user, logout, updateUser } = useAuth();
  const { pathname } = useLocation();
  const navigate = useNavigate();
  const theme = useTheme();
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'));
  const [mobileOpen, setMobileOpen] = useState(false);
  const [logoutError, setLogoutError] = useState('');
  const title = titles[pathname] ?? (pathname.startsWith('/review') ? 'Review Queue' : 'Invoice Recon');
  const userName = user?.full_name || 'User';
  const userInitials = userName.split(/\s+/).map((part) => part[0]).join('').slice(0, 2).toUpperCase();

  const handleLogout = async () => {
    setLogoutError('');
    try {
      await logout();
      navigate('/login', { replace: true });
    } catch (error) {
      setLogoutError(error.message);
    }
  };

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh', bgcolor: 'background.default' }}>
      <Sidebar
        variant={isDesktop ? 'permanent' : 'temporary'}
        open={isDesktop || mobileOpen}
        onClose={() => setMobileOpen(false)}
        onLogout={handleLogout}
      />
      <Box component="main" sx={{ flexGrow: 1, minWidth: 0 }}>
        <AppBar position="sticky" color="inherit" elevation={0} sx={{ borderBottom: 1, borderColor: 'divider', bgcolor: 'rgba(255,255,255,0.94)', backdropFilter: 'blur(12px)' }}>
          <Toolbar sx={{ minHeight: { xs: 64, sm: 72 }, px: { xs: 2, sm: 3, lg: 4 }, gap: 1.5 }}>
            {!isDesktop && (
              <IconButton
                edge="start"
                color="inherit"
                aria-label="Open navigation"
                onClick={() => setMobileOpen(true)}
                sx={{ color: 'text.secondary' }}
              >
                <MenuIcon />
              </IconButton>
            )}
            <Box sx={{ flexGrow: 1, minWidth: 0 }}>
              <Typography variant="h6" noWrap sx={{ fontSize: { xs: 17, sm: 19 }, fontWeight: 600, lineHeight: 1.25 }}>
                {title}
              </Typography>
              <Typography variant="caption" color="text.secondary" sx={{ display: { xs: 'none', sm: 'block' }, mt: 0.25 }}>
                Settlement workspace
              </Typography>
            </Box>
            <Divider orientation="vertical" flexItem sx={{ my: 1.5, display: { xs: 'none', sm: 'block' } }} />
            <Stack direction="row" spacing={{ xs: 0.75, sm: 1.25 }} alignItems="center" sx={{ pl: { xs: 0, sm: 0.75 } }}>
              <Stack
                component={RouterLink}
                to="/profile"
                direction="row"
                spacing={{ xs: 0, sm: 1 }}
                alignItems="center"
                aria-label="Open profile"
                sx={{
                  color: 'inherit',
                  textDecoration: 'none',
                  borderRadius: 2,
                  p: 0.5,
                  '&:hover': { bgcolor: 'action.hover' },
                }}
              >
                <Badge
                  overlap="circular"
                  variant="dot"
                  anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                  sx={{ '& .MuiBadge-badge': { bgcolor: '#22a06b', boxShadow: '0 0 0 2px #fff' } }}
                >
                  <Avatar sx={{ width: 38, height: 38, bgcolor: 'primary.light', color: 'primary.dark', fontSize: 13, fontWeight: 700 }}>
                    {userInitials}
                  </Avatar>
                </Badge>
                <Box sx={{ display: { xs: 'none', sm: 'block' }, minWidth: 0 }}>
                  <Typography variant="body2" sx={{ fontWeight: 600, lineHeight: 1.25, whiteSpace: 'nowrap' }}>
                    {userName}
                  </Typography>
                  <Typography variant="caption" color="text.secondary" sx={{ whiteSpace: 'nowrap' }}>
                    Profile
                  </Typography>
                </Box>
              </Stack>
            </Stack>
          </Toolbar>
        </AppBar>
        <Box sx={{ width: '100%', maxWidth: 1600, mx: 'auto', px: { xs: 2, sm: 3, lg: 4 }, py: { xs: 2.5, sm: 3.5 } }}>
          {logoutError && <Alert severity="error" sx={{ mb: 2 }}>{logoutError}</Alert>}
          <Outlet context={{ user, updateUser }} />
        </Box>
      </Box>
    </Box>
  );
}
