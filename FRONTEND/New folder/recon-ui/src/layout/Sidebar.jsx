import { NavLink } from 'react-router-dom';
import { Avatar, Box, Drawer, List, ListItemButton, ListItemIcon, ListItemText, Chip, Stack, Typography } from '@mui/material';
import DashboardIcon from '@mui/icons-material/Dashboard';
import DescriptionIcon from '@mui/icons-material/Description';
import ReportProblemIcon from '@mui/icons-material/ReportProblem';
import MenuBookIcon from '@mui/icons-material/MenuBook';
import PersonIcon from '@mui/icons-material/Person';

export const DRAWER_WIDTH = 248;

const items = [
  { to: '/', label: 'Dashboard', icon: <DashboardIcon /> },
  { to: '/documents', label: 'Documents', icon: <DescriptionIcon /> },
  { to: '/review', label: 'Review Queue', icon: <ReportProblemIcon />, badge: 14 },
  { to: '/ledger', label: 'Ledger', icon: <MenuBookIcon /> },
  { to: '/users', label: 'Users', icon: <PersonIcon /> },
];

const navSx = {
  mb: 0.5,
  borderRadius: 2,
  py: 1,
  color: '#cbd8e8',
  '& .MuiListItemIcon-root': { color: '#91a9c5', minWidth: 40 },
  '& .MuiListItemText-primary': { fontSize: 14, fontWeight: 500 },
  transition: 'background-color .15s, color .15s',
  '&:hover': { bgcolor: 'rgba(255,255,255,0.08)', color: '#ffffff', '& .MuiListItemIcon-root': { color: '#ffffff' } },
  '&.active': {
    bgcolor: 'rgba(87, 155, 232, 0.2)',
    color: '#ffffff',
    boxShadow: 'inset 3px 0 0 #72b5ff',
    '& .MuiListItemIcon-root': { color: '#a9d2ff' },
    '& .MuiListItemText-primary': { fontWeight: 600 },
    '&:hover': { bgcolor: 'rgba(87, 155, 232, 0.26)' },
  },
};

export default function Sidebar({ variant = 'permanent', open = true, onClose }) {
  return (
    <Drawer
      variant={variant}
      open={open}
      onClose={onClose}
      sx={{
        width: DRAWER_WIDTH,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: DRAWER_WIDTH,
          boxSizing: 'border-box',
          color: '#fff',
          borderRight: 0,
          background: 'linear-gradient(180deg, #123b64 0%, #0d2f51 58%, #0a2744 100%)',
        },
      }}
    >
      <Stack direction="row" spacing={1.5} alignItems="center" sx={{ px: 2.5, pt: 3, pb: 2.75 }}>
        <Avatar variant="rounded" sx={{ bgcolor: '#2d6ca7', width: 40, height: 40, fontWeight: 700, fontSize: 15, borderRadius: 2 }}>IR</Avatar>
        <Box>
          <Typography sx={{ fontWeight: 700, fontSize: 17, lineHeight: 1.2 }}>Invoice Recon</Typography>
          <Typography sx={{ fontSize: 12, color: '#a8bdd4' }}>Settlement Engine</Typography>
        </Box>
      </Stack>

      <Typography sx={{ px: 3, pb: 1, fontSize: 10, fontWeight: 700, letterSpacing: 1.3, color: '#8ea8c4' }}>WORKSPACE</Typography>
      <List sx={{ px: 1.5 }}>
        {items.map((it) => (
          <ListItemButton key={it.to} component={NavLink} to={it.to} end={it.to === '/'} onClick={onClose} sx={navSx}>
            <ListItemIcon>{it.icon}</ListItemIcon>
            <ListItemText primary={it.label} />
            {it.badge && (
              <Chip label={it.badge} size="small" sx={{ bgcolor: '#ff9800', color: '#1a1a1a', fontWeight: 700, height: 20, minWidth: 28 }} />
            )}
          </ListItemButton>
        ))}
      </List>

      <Box sx={{ flexGrow: 1 }} />
      <Stack direction="row" spacing={1.25} alignItems="center" sx={{ m: 1.5, p: 1.5, borderRadius: 2, bgcolor: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.08)' }}>
        <Avatar sx={{ bgcolor: '#2d6ca7', width: 34, height: 34, fontSize: 12, fontWeight: 600 }}>AD</Avatar>
        <Box>
          <Typography sx={{ fontSize: 13, fontWeight: 600 }}>Admin User</Typography>
          <Typography sx={{ fontSize: 11, color: '#a8bdd4' }}>Administrator</Typography>
        </Box>
      </Stack>
    </Drawer>
  );
}
