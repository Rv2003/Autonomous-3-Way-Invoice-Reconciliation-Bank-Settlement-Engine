import { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  Avatar, Box, Button, Card, CardContent, Chip, Divider, LinearProgress, Link as MuiLink, List, ListItem,
  ListItemAvatar, ListItemText, Stack, Typography,
} from '@mui/material';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import ChatIcon from '@mui/icons-material/Chat';
import EmailIcon from '@mui/icons-material/Email';
import DescriptionIcon from '@mui/icons-material/Description';
import PageHeader from '../components/PageHeader';
import StatCard from '../components/StatCard';
import StatusChip from '../components/StatusChip';
import SourceCell from '../components/SourceCell';
import DataTable from '../components/DataTable';
import { docs, whatsappFeed, gmailFeed } from '../data/mock';

const types = ['Auto-detect', 'Supplier invoice', 'Purchase order', 'Bank transfer'];

function PanelHeader({ icon, color, bg, title, sub, tag }) {
  return (
    <Stack direction="row" spacing={1.5} alignItems="center" sx={{ mb: 2 }}>
      <Avatar sx={{ bgcolor: bg, color, width: 40, height: 40 }}>{icon}</Avatar>
      <Box sx={{ flexGrow: 1 }}>
        <Typography fontWeight={500}>{title}</Typography>
        <Typography variant="caption" color="text.secondary">{sub}</Typography>
      </Box>
      {tag}
    </Stack>
  );
}

function Feed({ items, color, bg }) {
  return (
    <List disablePadding>
      {items.map((f, i) => (
        <Box key={f.name}>
          {i > 0 && <Divider component="li" />}
          <ListItem disableGutters secondaryAction={<StatusChip status={f.status} />}>
            <ListItemAvatar><Avatar sx={{ bgcolor: bg, color, width: 36, height: 36 }}><DescriptionIcon fontSize="small" /></Avatar></ListItemAvatar>
            <ListItemText primary={f.name} secondary={f.meta} primaryTypographyProps={{ fontWeight: 500, fontSize: 14 }} sx={{ pr: 11 }} />
          </ListItem>
        </Box>
      ))}
    </List>
  );
}

export default function Dashboard() {
  const [type, setType] = useState('Auto-detect');
  return (
    <>
      <PageHeader
        title="Document Ingestion"
        subtitle="Upload files manually or let WhatsApp and Gmail feed documents into the reconciliation pipeline."
        actions={<Button variant="outlined">Refresh</Button>}
      />
      <Box sx={{ display: 'grid', gap: 2, gridTemplateColumns: { xs: '1fr', sm: 'repeat(2,1fr)', lg: 'repeat(4,1fr)' }, mb: 3 }}>
        <StatCard label="Documents today" value="128" sub="+12% vs yesterday" />
        <StatCard label="In pipeline" value="6" sub="Processing now" />
        <StatCard label="Needs review" value="14" sub="In exceptions queue" color="warning.main" />
        <StatCard label="Failed extraction" value="3" sub="Check source files" color="error.main" />
      </Box>

      <Box sx={{ display: 'grid', gap: 2, gridTemplateColumns: { xs: '1fr', lg: 'repeat(3,1fr)' }, mb: 3 }}>
        <Card><CardContent>
          <PanelHeader icon={<CloudUploadIcon />} color="#1976d2" bg="#e3f2fd" title="Manual Upload" sub="Admin upload to the pipeline" />
          <Box sx={{ border: '2px dashed #90caf9', bgcolor: '#f3f9ff', borderRadius: 2, p: 3, textAlign: 'center' }}>
            <CloudUploadIcon color="primary" sx={{ fontSize: 40 }} />
            <Typography fontWeight={500}>Drag & drop invoices here</Typography>
            <Typography variant="caption" color="text.secondary" display="block" sx={{ mb: 1.5 }}>PDF files only · up to 10 MB each</Typography>
            <Button variant="contained" size="small">Browse files</Button>
          </Box>
          <Typography variant="body2" fontWeight={500} sx={{ mt: 2, mb: 1 }}>Document type</Typography>
          <Stack direction="row" spacing={1} useFlexGap flexWrap="wrap">
            {types.map((t) => (
              <Chip key={t} label={t} size="small" clickable onClick={() => setType(t)}
                color={type === t ? 'primary' : 'default'} variant={type === t ? 'filled' : 'outlined'} />
            ))}
          </Stack>
          <Typography variant="body2" fontWeight={500} sx={{ mt: 2 }}>Upload queue</Typography>
          <Box sx={{ mt: 1 }}>
            <Stack direction="row" justifyContent="space-between"><Typography variant="caption">INV-2044_Cargills.pdf</Typography><Typography variant="caption" color="text.secondary">78%</Typography></Stack>
            <LinearProgress variant="determinate" value={78} sx={{ height: 6, borderRadius: 3, my: 0.5 }} />
            <Stack direction="row" justifyContent="space-between" sx={{ mt: 1 }}><Typography variant="caption">BT-2211_Slip.pdf</Typography><Typography variant="caption" color="text.secondary">Uploaded</Typography></Stack>
            <LinearProgress variant="determinate" value={100} color="success" sx={{ height: 6, borderRadius: 3, my: 0.5 }} />
          </Box>
          <Button fullWidth variant="contained" sx={{ mt: 2 }}>Upload & process</Button>
        </CardContent></Card>

        <Card><CardContent>
          <PanelHeader icon={<ChatIcon />} color="#00897b" bg="#e0f2f1" title="Live WhatsApp Extraction" sub="Incoming documents via webhook" tag={<Chip size="small" label="● Live" sx={{ bgcolor: '#e8f5e9', color: '#2e7d32', fontWeight: 500 }} />} />
          <Box sx={{ bgcolor: '#f5f7fa', borderRadius: 2, p: 1.5, mb: 1 }}>
            <Stack direction="row" justifyContent="space-between"><Typography variant="caption" color="text.secondary">Business number</Typography><Typography variant="caption" fontWeight={500}>+94 77 ••• 4567</Typography></Stack>
            <Stack direction="row" justifyContent="space-between"><Typography variant="caption" color="text.secondary">Webhook</Typography><Typography variant="caption" fontWeight={500}>Active · 38 today</Typography></Stack>
          </Box>
          <Feed items={whatsappFeed} color="#1976d2" bg="#e0f2f1" />
          <Button fullWidth variant="outlined" sx={{ mt: 1.5 }}>View all messages</Button>
        </CardContent></Card>

        <Card><CardContent>
          <PanelHeader icon={<EmailIcon />} color="#e53935" bg="#ffebee" title="Gmail Extraction" sub="Attachments from the inbox" tag={<Chip size="small" label="Synced" sx={{ bgcolor: '#e8f5e9', color: '#2e7d32', fontWeight: 500 }} />} />
          <Box sx={{ bgcolor: '#f5f7fa', borderRadius: 2, p: 1.5, mb: 1 }}>
            <Stack direction="row" justifyContent="space-between"><Typography variant="caption" color="text.secondary">Mailbox</Typography><Typography variant="caption" fontWeight={500}>finance@company.lk</Typography></Stack>
            <Stack direction="row" justifyContent="space-between"><Typography variant="caption" color="text.secondary">Last sync</Typography><Typography variant="caption" fontWeight={500}>1 min ago · 24 today</Typography></Stack>
          </Box>
          <Feed items={gmailFeed} color="#1976d2" bg="#ffebee" />
          <Button fullWidth variant="outlined" sx={{ mt: 1.5 }}>Sync now</Button>
        </CardContent></Card>
      </Box>

      <Card>
        <Stack direction="row" justifyContent="space-between" sx={{ p: 2.5 }}>
          <Typography variant="h6" sx={{ fontSize: 17 }}>Recent documents</Typography>
          <MuiLink component={RouterLink} to="/documents" underline="none" fontWeight={500}>View all</MuiLink>
        </Stack>
        <DataTable
          columns={[
            { key: 'name', label: 'Document', strong: true },
            { key: 'source', label: 'Source', render: (r) => <SourceCell source={r.source} /> },
            { key: 'type', label: 'Type' },
            { key: 'company', label: 'Company' },
            { key: 'status', label: 'Status', render: (r) => <StatusChip status={r.status} /> },
          ]}
          rows={docs.slice(0, 5)}
        />
      </Card>
    </>
  );
}
