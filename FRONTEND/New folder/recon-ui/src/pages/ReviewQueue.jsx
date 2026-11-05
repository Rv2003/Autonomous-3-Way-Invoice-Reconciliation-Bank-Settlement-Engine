import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Tab, Tabs } from '@mui/material';
import PageHeader from '../components/PageHeader';
import StatCard from '../components/StatCard';
import DataTable from '../components/DataTable';
import StatusChip from '../components/StatusChip';
import { exceptions } from '../data/mock';

const tabs = [
  ['All', null],
  ['Amount mismatch', 'Amount mismatch'],
  ['Missing PO', 'Missing PO'],
  ['Duplicate', 'Duplicate invoice'],
  ['Unmatched payment', 'Unmatched payment'],
];

export default function ReviewQueue() {
  const nav = useNavigate();
  const [tab, setTab] = useState(0);
  const filter = tabs[tab][1];
  const rows = filter ? exceptions.filter((e) => e.issue === filter) : exceptions;
  const count = (f) => (f ? exceptions.filter((e) => e.issue === f).length : exceptions.length);

  return (
    <>
      <PageHeader
        title="Review Queue"
        subtitle="Exceptions the engine could not resolve automatically. Approve, reject or request a correction."
        actions={<Button variant="outlined">Refresh</Button>}
      />
      <Box sx={{ display: 'grid', gap: 2, gridTemplateColumns: { xs: '1fr', md: 'repeat(3,1fr)' }, mb: 3 }}>
        <StatCard label="Open exceptions" value="14" sub="5 added today" color="warning.main" />
        <StatCard label="Amount at risk" value="LKR 4.82M" sub="Across 14 invoices" />
        <StatCard label="Avg. resolution time" value="3.2 h" sub="-18% vs last week" />
      </Box>
      <Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        {tabs.map(([label, f]) => <Tab key={label} label={`${label}  ${count(f)}`} />)}
      </Tabs>
      <DataTable
        onRowClick={(r) => nav(`/review/${r.id}`)}
        columns={[
          { key: 'id', label: 'Exception', strong: true },
          { key: 'doc', label: 'Document' },
          { key: 'company', label: 'Company' },
          { key: 'issue', label: 'Issue' },
          { key: 'variance', label: 'Variance', align: 'right', render: (r) => <Box component="span" sx={{ color: r.variance === '—' ? 'text.secondary' : 'error.main', fontWeight: 500 }}>{r.variance}</Box> },
          { key: 'age', label: 'Age' },
          { key: 'status', label: 'Status', render: (r) => <StatusChip status={r.status} /> },
          { key: 'a', label: '', render: (r) => <Button size="small" variant="outlined" onClick={(e) => { e.stopPropagation(); nav(`/review/${r.id}`); }}>Review</Button> },
        ]}
        rows={rows}
      />
    </>
  );
}
