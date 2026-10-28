import { useMemo, useState } from 'react';
import { Button, InputAdornment, MenuItem, Stack, TablePagination, TextField } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import PageHeader from '../components/PageHeader';
import DataTable from '../components/DataTable';
import StatusChip from '../components/StatusChip';
import SourceCell from '../components/SourceCell';
import { docs } from '../data/mock';

const statuses = ['All', 'Processed', 'Processing', 'Needs review', 'Failed'];
const sources = ['All', 'WhatsApp', 'Gmail', 'Manual upload'];

export default function Documents() {
  const [q, setQ] = useState('');
  const [status, setStatus] = useState('All');
  const [source, setSource] = useState('All');
  const [page, setPage] = useState(0);
  const [rpp, setRpp] = useState(10);

  const rows = useMemo(
    () => docs.filter((d) =>
      (status === 'All' || d.status === status) &&
      (source === 'All' || d.source === source) &&
      `${d.name} ${d.company}`.toLowerCase().includes(q.toLowerCase())),
    [q, status, source]
  );

  return (
    <>
      <PageHeader
        title="Documents"
        subtitle="Every document received from WhatsApp, Gmail and manual upload, with its classification and status."
        actions={<Button variant="contained">Upload document</Button>}
      />
      <Stack direction={{ xs: 'column', md: 'row' }} spacing={1.5} sx={{ mb: 2 }}>
        <TextField size="small" fullWidth placeholder="Search by file name, company or invoice number" value={q}
          onChange={(e) => { setQ(e.target.value); setPage(0); }}
          slotProps={{ input: { startAdornment: <InputAdornment position="start"><SearchIcon fontSize="small" /></InputAdornment> } }}
          sx={{ bgcolor: '#fff' }} />
        <TextField select size="small" label="Source" value={source} onChange={(e) => setSource(e.target.value)} sx={{ minWidth: 160, bgcolor: '#fff' }}>
          {sources.map((s) => <MenuItem key={s} value={s}>{s}</MenuItem>)}
        </TextField>
        <TextField select size="small" label="Status" value={status} onChange={(e) => setStatus(e.target.value)} sx={{ minWidth: 160, bgcolor: '#fff' }}>
          {statuses.map((s) => <MenuItem key={s} value={s}>{s}</MenuItem>)}
        </TextField>
        <Button variant="outlined">Export CSV</Button>
      </Stack>
      <DataTable
        columns={[
          { key: 'name', label: 'Document', strong: true },
          { key: 'source', label: 'Source', render: (r) => <SourceCell source={r.source} /> },
          { key: 'type', label: 'Type' },
          { key: 'company', label: 'Company' },
          { key: 'received', label: 'Received' },
          { key: 'confidence', label: 'Confidence' },
          { key: 'status', label: 'Status', render: (r) => <StatusChip status={r.status} /> },
        ]}
        rows={rows.slice(page * rpp, page * rpp + rpp)}
      />
      <TablePagination component="div" count={rows.length} page={page} rowsPerPage={rpp}
        onPageChange={(_, p) => setPage(p)} onRowsPerPageChange={(e) => { setRpp(+e.target.value); setPage(0); }} />
    </>
  );
}
