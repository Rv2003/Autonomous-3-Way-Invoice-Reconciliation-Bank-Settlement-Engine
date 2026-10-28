import { Card, CardContent, Typography } from '@mui/material';

export default function StatCard({ label, value, sub, color = 'text.primary' }) {
  return (
    <Card>
      <CardContent sx={{ '&:last-child': { pb: 2 } }}>
        <Typography variant="body2" color="text.secondary">{label}</Typography>
        <Typography sx={{ fontSize: 32, fontWeight: 700, lineHeight: 1.3, color }}>{value}</Typography>
        <Typography variant="caption" color="text.secondary">{sub}</Typography>
      </CardContent>
    </Card>
  );
}
