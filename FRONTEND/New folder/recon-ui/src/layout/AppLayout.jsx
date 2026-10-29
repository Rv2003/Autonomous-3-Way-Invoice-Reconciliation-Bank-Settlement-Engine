import { useState } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import { AppBar, Avatar, Badge, Box, Divider, IconButton, Stack, Toolbar, Typography, useMediaQuery, useTheme } from '@mui/material';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import MenuIcon from '@mui/icons-material/Menu';
import Sidebar from './Sidebar';

const titles = { '/': 'Dashboard', '/documents': 'Documents', '/review': 'Review Queue', '/ledger': 'Ledger', '/users': 'Users' };

export default function AppLayout() {
  const { pathname } = useLocation();
  const theme = useTheme();
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'));
  const [mobileOpen, setMobileOpen] = useState(false);
  const title = titles[pathname] ?? (pathname.startsWith('/review') ? 'Review Queue' : 'Invoice Recon');

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
                <Avatar sx={{ width: 38, height: 38, bgcolor: '#e8eef8', color: '#264c7a', fontSize: 13, fontWeight: 700 }}>
                  AD
                </Avatar>
              </Badge>
              <Box sx={{ display: { xs: 'none', sm: 'block' }, minWidth: 0 }}>
                <Typography variant="body2" sx={{ fontWeight: 600, lineHeight: 1.25, whiteSpace: 'nowrap' }}>
                  Admin User
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ whiteSpace: 'nowrap' }}>
                  Workspace admin
                </Typography>
              </Box>
              <KeyboardArrowDownIcon sx={{ display: { xs: 'none', sm: 'block' }, color: 'text.secondary', fontSize: 20 }} />
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
