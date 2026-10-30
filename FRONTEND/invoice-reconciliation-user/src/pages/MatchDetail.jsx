import { useState } from 'react';
import { Link as RouterLink, useParams } from 'react-router-dom';
import { Alert, AlertTitle, Box, Breadcrumbs, Button, Card, CardContent, Link, Snackbar, Stack, TextField, Typography } from '@mui/material';
import PageHeader from '../components/PageHeader';
import StatusChip from '../components/StatusChip';
import DataTable from '../components/DataTable';
import { lineItems } from '../data/mock';

const KV = ({ label, value, color }) => (
  <Stack direction="row" justifyContent="space-between" sx={{ py: 0.75 }}>
    <Typography variant="body2" color="text.secondary">{label}</Typography>
    <Typography variant="body2" fontWeight={500} color={color}>{value}</Typography>
  </Stack>
);

const DocCard = ({ title, status, rows }) => (
  <Card><CardContent>
    <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
      <Typography fontWeight={500}>{title}</Typography>
      <StatusChip status={status} />
    </Stack>
    {rows.map(([l, v, c]) => <KV key={l} label={l} value={v} color={c} />)}
  </CardContent></Card>
);

export default function MatchDetail() {
  const { id } = useParams();
  const [note, setNote] = useState('');
  const [msg, setMsg] = useState('');
  const decide = (label) => setMsg(`${label} · ${id}`);

  return (
    <>
      <Breadcrumbs sx={{ mb: 1, fontSize: 13 }}>
        <Link component={RouterLink} to="/review" underline="hover" color="text.secondary">Review Queue</Link>
        <Typography fontSize={13}>{id}</Typography>
      </Breadcrumbs>
      <PageHeader
        title={`${id} · Keells Super (Pvt) Ltd`}
        subtitle="Supplier invoice INV-2043 compared with purchase order PO-8812 and bank transfer BT-2210."
        actions={[
          <Button key="c" variant="outlined" onClick={() => decide('Correction requested')}>Request correction</Button>,
          <Button key="r" variant="outlined" color="error" onClick={() => decide('Rejected')}>Reject</Button>,
          <Button key="a" variant="contained" color="success" onClick={() => decide('Approved')}>Approve with note</Button>,
        ]}
      />
      <Alert severity="warning" sx={{ mb: 2, alignItems: 'center' }}
        action={<Box sx={{ textAlign: 'right', pr: 1 }}><Typography variant="caption" color="text.secondary">Match confidence</Typography><Typography sx={{ fontSize: 22, fontWeight: 700, color: 'warning.main' }}>82%</Typography></Box>}>
        <AlertTitle sx={{ fontWeight: 700 }}>Amount mismatch detected</AlertTitle>
        The invoice total is LKR 12,500.00 higher than the purchase order. The bank transfer matches the invoice total.
      </Alert>
      <Box sx={{ display: 'grid', gap: 2, gridTemplateColumns: { xs: '1fr', lg: 'repeat(3,1fr)' }, mb: 2 }}>
        <DocCard title="Purchase order" status="Matched" rows={[['Number', 'PO-8812'], ['Supplier', 'Keells Super (Pvt) Ltd'], ['Issued', 'Oct 12, 2026'], ['Currency', 'LKR'], ['Total', '487,500.00']]} />
        <DocCard title="Supplier invoice" status="Mismatch" rows={[['Number', 'INV-2043'], ['Supplier', 'Keells Super (Pvt) Ltd'], ['Issued', 'Oct 20, 2026'], ['PO reference', 'PO-8812'], ['Total', '500,000.00', 'error.main']]} />
        <DocCard title="Bank transfer" status="Matched" rows={[['Reference', 'BT-2210'], ['Payer', 'Dialog Axiata PLC'], ['Value date', 'Oct 25, 2026'], ['Bank', 'Commercial Bank'], ['Amount', '500,000.00']]} />
      </Box>
      <Card sx={{ mb: 2 }}><CardContent>
        <Typography fontWeight={500} sx={{ mb: 1.5 }}>Line item comparison</Typography>
        <DataTable
          columns={[
            { key: 'item', label: 'Item', strong: true },
            { key: 'poQty', label: 'PO qty', align: 'right' },
            { key: 'poPrice', label: 'PO price', align: 'right' },
            { key: 'invQty', label: 'Invoice qty', align: 'right' },
            { key: 'invPrice', label: 'Invoice price', align: 'right' },
            { key: 'variance', label: 'Variance', align: 'right', render: (r) => <Box component="span" sx={{ color: r.variance === '0.00' ? 'text.secondary' : 'error.main', fontWeight: 500 }}>{r.variance}</Box> },
          ]}
          rows={lineItems.map((l) => ({ ...l, id: l.item }))}
        />
      </CardContent></Card>
      <Card><CardContent>
        <Typography fontWeight={500}>Reviewer note</Typography>
        <Typography variant="caption" color="text.secondary">This note is stored in the hash-chained ledger with your decision.</Typography>
        <TextField fullWidth multiline minRows={2} size="small" placeholder="Add a note for the audit trail…" value={note} onChange={(e) => setNote(e.target.value)} sx={{ mt: 1.5 }} />
      </CardContent></Card>
      <Snackbar open={!!msg} autoHideDuration={3000} onClose={() => setMsg('')} message={msg} />
    </>
  );
}
