import { useState } from 'react';
import { Box, Button, Card, CardContent, Tab, Tabs, Typography } from '@mui/material';
import PageHeader from '../components/PageHeader';
import StatCard from '../components/StatCard';
import DataTable from '../components/DataTable';
import StatusChip from '../components/StatusChip';
import { ledger, settlements } from '../data/mock';

const hash = (r) => <Box component="span" sx={{ fontFamily: 'ui-monospace, Menlo, monospace', fontSize: 13 }}>{r.hash}</Box>;

export default function Ledger() {
  const [tab, setTab] = useState(0);
  return (
    <>
      <PageHeader
        title="Ledger"
        subtitle="Hash-chained record of every reconciliation and settlement event."
        actions={[<Button key="e" variant="outlined">Export</Button>, <Button key="v" variant="contained">Verify chain</Button>]}
      />
      <Box sx={{ display: 'grid', gap: 2, gridTemplateColumns: { xs: '1fr', sm: 'repeat(2,1fr)', lg: 'repeat(4,1fr)' }, mb: 3 }}>
        <StatCard label="Ledger entries" value="12,480" sub="+86 today" />
        <StatCard label="Chain integrity" value="Verified" sub="Last check 2 min ago" color="success.main" />
        <StatCard label="Settled today" value="LKR 8.4M" sub="31 payments" />
        <StatCard label="Unsettled payments" value="7" sub="Awaiting bank match" color="warning.main" />
      </Box>
      <Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Tab label="Ledger entries" />
        <Tab label="Bank settlement" />
      </Tabs>
      {tab === 0 ? (
        <DataTable
          columns={[
            { key: 'id', label: 'Block', strong: true },
            { key: 'time', label: 'Time' },
            { key: 'event', label: 'Event' },
            { key: 'invoice', label: 'Invoice' },
            { key: 'amount', label: 'Amount (LKR)', align: 'right' },
            { key: 'hash', label: 'Hash', render: hash },
            { key: 'status', label: 'Status', render: (r) => <StatusChip status={r.status} /> },
          ]}
          rows={ledger}
        />
      ) : (
        <Card><CardContent>
          <Typography fontWeight={500} sx={{ mb: 1.5 }}>Bank settlement matching</Typography>
          <DataTable
            columns={[
              { key: 'id', label: 'Bank transaction', strong: true },
              { key: 'payer', label: 'Payer' },
              { key: 'amount', label: 'Amount (LKR)', align: 'right' },
              { key: 'invoice', label: 'Matched invoice' },
              { key: 'date', label: 'Value date' },
              { key: 'status', label: 'Status', render: (r) => <StatusChip status={r.status} /> },
            ]}
            rows={settlements}
          />
        </CardContent></Card>
      )}
    </>
  );
}
