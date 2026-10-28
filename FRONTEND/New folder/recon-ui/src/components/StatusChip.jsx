import { Chip } from '@mui/material';
import { statusColors } from '../theme';

export default function StatusChip({ status }) {
  const [bg, fg] = statusColors[status] ?? statusColors.Processing;
  return <Chip size="small" label={status} sx={{ bgcolor: bg, color: fg, fontWeight: 500, height: 24 }} />;
}
