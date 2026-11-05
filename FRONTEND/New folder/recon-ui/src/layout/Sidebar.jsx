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
  color: '#dbeafe',
  '& .MuiListItemIcon-root': { color: '#bfdbfe', minWidth: 40 },
  '& .MuiListItemText-primary': { fontSize: 14, fontWeight: 500 },
  transition: 'background-color .15s, color .15s',
  '&:hover': { bgcolor: 'rgba(255,255,255,0.08)', color: '#ffffff', '& .MuiListItemIcon-root': { color: '#ffffff' } },
  '&.active': {
    bgcolor: 'rgba(255,255,255,0.16)',
    color: '#ffffff',
    boxShadow: 'inset 3px 0 0 #bfdbfe',
    '& .MuiListItemIcon-root': { color: '#ffffff' },
    '& .MuiListItemText-primary': { fontWeight: 600 },
    '&:hover': { bgcolor: 'rgba(255,255,255,0.22)' },
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
          background: 'linear-gradient(180deg, #1d4ed8 0%, #1e40af 58%, #1e3a8a 100%)',
        },
      }}
    >
      <Stack direction="row" spacing={1.5} alignItems="center" sx={{ px: 2.5, pt: 3, pb: 2.75 }}>
        <Avatar variant="rounded" sx={{ bgcolor: '#2563eb', width: 40, height: 40, fontWeight: 700, fontSize: 15, borderRadius: 2 }}>IR</Avatar>
        <Box>
          <Typography sx={{ fontWeight: 700, fontSize: 17, lineHeight: 1.2 }}>Invoice Recon</Typography>
          <Typography sx={{ fontSize: 12, color: '#a8bdd4' }}>Settlement Engine</Typography>
        </Box>
      </Stack>

      <Typography sx={{ px: 3, pb: 1, fontSize: 10, fontWeight: 700, letterSpacing: 1.3, color: '#bfdbfe' }}>WORKSPACE</Typography>
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
        <Avatar sx={{ bgcolor: '#2563eb', width: 34, height: 34, fontSize: 12, fontWeight: 600 }}>AD</Avatar>
        <Box>
          <Typography sx={{ fontSize: 13, fontWeight: 600 }}>Admin User</Typography>
          <Typography sx={{ fontSize: 11, color: '#bfdbfe' }}>Administrator</Typography>
        </Box>
      </Stack>
    </Drawer>
  );
}
