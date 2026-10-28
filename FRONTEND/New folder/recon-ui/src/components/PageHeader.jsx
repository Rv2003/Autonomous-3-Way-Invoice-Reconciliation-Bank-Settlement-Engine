import { Stack, Typography, Box } from '@mui/material';

export default function PageHeader({ title, subtitle, actions }) {
  return (
    <Stack
      direction={{ xs: 'column', sm: 'row' }}
      justifyContent="space-between"
      alignItems={{ xs: 'stretch', sm: 'center' }}
      spacing={{ xs: 1.5, sm: 2 }}
      sx={{ mb: 3, minWidth: 0 }}
    >
      <Box sx={{ minWidth: 0 }}>
        <Typography variant="h4" sx={{ fontSize: { xs: '1.4rem', sm: '1.625rem' } }}>{title}</Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5, maxWidth: 760 }}>{subtitle}</Typography>
      </Box>
      {actions && <Stack direction="row" spacing={1.5} sx={{ alignSelf: { xs: 'flex-start', sm: 'auto' }, flexShrink: 0 }}>{actions}</Stack>}
    </Stack>
  );
}
