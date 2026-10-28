import { createTheme } from '@mui/material/styles';

export const statusColors = {
  Processed: ['#e8f5e9', '#2e7d32'],
  Matched: ['#e8f5e9', '#2e7d32'],
  Verified: ['#e8f5e9', '#2e7d32'],
  Active: ['#e8f5e9', '#2e7d32'],
  Processing: ['#e3f2fd', '#1565c0'],
  'In review': ['#e3f2fd', '#1565c0'],
  Invited: ['#e3f2fd', '#1565c0'],
  'Needs review': ['#fff3e0', '#e65100'],
  Pending: ['#fff3e0', '#e65100'],
  Open: ['#fff3e0', '#e65100'],
  Failed: ['#ffebee', '#c62828'],
  Mismatch: ['#ffebee', '#c62828'],
  Escalated: ['#ffebee', '#c62828'],
  Disabled: ['#eceff3', '#5f6b7a'],
};

export const sourceColors = { WhatsApp: '#25d366', Gmail: '#e53935', 'Manual upload': '#1976d2' };

export default createTheme({
  palette: {
    primary: { main: '#2469a6', dark: '#174d7d', light: '#e8f1fa' },
    success: { main: '#2e7d32' },
    warning: { main: '#e65100' },
    error: { main: '#c62828' },
    background: { default: '#f4f6f9', paper: '#ffffff' },
    text: { primary: '#1a2027', secondary: '#5f6b7a' },
    divider: '#e5e9ef',
  },
  shape: { borderRadius: 10 },
  typography: {
    fontFamily: 'Roboto, "Helvetica Neue", Arial, sans-serif',
    h4: { fontWeight: 650, fontSize: '1.625rem', letterSpacing: '-0.025em' },
    h6: { fontWeight: 600, letterSpacing: '-0.01em' },
    button: { textTransform: 'none', fontWeight: 600, letterSpacing: 0 },
  },
  components: {
    MuiButton: { defaultProps: { disableElevation: true }, styleOverrides: { root: { borderRadius: 8, minHeight: 38 } } },
    MuiCard: {
      defaultProps: { variant: 'outlined' },
      styleOverrides: { root: { borderColor: '#e5e9ef', borderRadius: 12, boxShadow: '0 1px 2px rgba(16, 24, 40, 0.03)' } },
    },
    MuiTableCell: {
      styleOverrides: {
        head: { fontSize: 11, fontWeight: 600, color: '#697586', textTransform: 'uppercase', letterSpacing: '0.04em', backgroundColor: '#f7f8fa' },
      },
    },
  },
});
