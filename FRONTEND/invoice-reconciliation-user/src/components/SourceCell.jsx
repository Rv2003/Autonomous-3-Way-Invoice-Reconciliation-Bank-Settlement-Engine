import { Box, Stack, Typography } from '@mui/material';
import { sourceColors } from '../theme';

export default function SourceCell({ source }) {
  return (
    <Stack direction="row" spacing={1} alignItems="center">
      <Box sx={{ width: 8, height: 8, borderRadius: '50%', bgcolor: sourceColors[source] }} />
      <Typography variant="body2">{source}</Typography>
    </Stack>
  );
}
