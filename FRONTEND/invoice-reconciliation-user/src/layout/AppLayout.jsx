import { useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { AppBar, Avatar, Badge, Box, Button, Divider, IconButton, Stack, Toolbar, Typography, useMediaQuery, useTheme } from '@mui/material';
import LogoutIcon from '@mui/icons-material/Logout';
import MenuIcon from '@mui/icons-material/Menu';
import useAuth from '../auth/useAuth';
import Sidebar from './Sidebar';

const titles = { '/': 'Dashboard', '/documents': 'Documents', '/review': 'Review Queue', '/ledger': 'Ledger', '/users': 'Users' };

export default function AppLayout() {
  const { pathname } = useLocation();
  const theme = useTheme();
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'));
  const [mobileOpen, setMobileOpen] = useState(false);
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const title = titles[pathname] ?? (pathname.startsWith('/review') ? 'Review Queue' : 'Invoice Recon');
  const initials = user.full_name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
    .toUpperCase();
  const roleLabel = user.role === 'administrator' ? 'Administrator' : 'Clerk';

  const handleLogout = () => {
    logout();
    navigate('/login', { replace: true });
  };

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh', bgcolor: 'background.default' }}>
      <Sidebar
        variant={isDesktop ? 'permanent' : 'temporary'}
        open={isDesktop || mobileOpen}
        onClose={() => setMobileOpen(false)}
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
              <Badge
                overlap="circular"
                variant="dot"
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                sx={{ '& .MuiBadge-badge': { bgcolor: '#22a06b', boxShadow: '0 0 0 2px #fff' } }}
              >
                <Avatar sx={{ width: 38, height: 38, bgcolor: 'primary.light', color: 'primary.dark', fontSize: 13, fontWeight: 700 }}>
                  {initials}
                </Avatar>
              </Badge>
              <Box sx={{ display: { xs: 'none', sm: 'block' }, minWidth: 0 }}>
                <Typography variant="body2" sx={{ fontWeight: 600, lineHeight: 1.25, whiteSpace: 'nowrap' }}>
                  {user.full_name}
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ whiteSpace: 'nowrap' }}>
                  {roleLabel}
                </Typography>
              </Box>
              <Button
                color="inherit"
                size="small"
                startIcon={<LogoutIcon />}
                onClick={handleLogout}
                aria-label="Sign out"
                sx={{ color: 'text.secondary', ml: { xs: 0, sm: 1 } }}
              >
                <Box component="span" sx={{ display: { xs: 'none', sm: 'inline' } }}>Sign out</Box>
              </Button>
            </Stack>
          </Toolbar>
        </AppBar>
        <Box sx={{ width: '100%', maxWidth: 1600, mx: 'auto', px: { xs: 2, sm: 3, lg: 4 }, py: { xs: 2.5, sm: 3.5 } }}>
          <Outlet />
        </Box>
      </Box>
    </Box>
  );
}
