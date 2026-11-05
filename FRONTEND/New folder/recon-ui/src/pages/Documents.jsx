import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Alert, Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, InputAdornment, MenuItem, Stack, TablePagination, TextField, Tooltip, Typography } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import VisibilityOutlinedIcon from '@mui/icons-material/VisibilityOutlined';
import PageHeader from '../components/PageHeader';
import DataTable from '../components/DataTable';
import StatusChip from '../components/StatusChip';
import SourceCell from '../components/SourceCell';
import PdfUploadQueue from '../components/PdfUploadQueue';
import usePdfUploadQueue from '../hooks/usePdfUploadQueue';
import { deleteDocument, displayDocumentStatus, fetchDocuments, getDocumentUrl } from '../components/api/documents';

const statuses = ['All', 'Uploaded', 'Processed', 'Processing', 'Failed'];
const sources = ['All', 'WhatsApp', 'Gmail', 'Manual upload'];

function formatReceived(value) {
  if (!value) return '—';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '—' : date.toLocaleString();
}

function formatFileSize(bytes) {
  if (!Number.isFinite(bytes)) return '—';
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

export default function Documents() {
  const [q, setQ] = useState('');
  const [status, setStatus] = useState('All');
  const [source, setSource] = useState('All');
  const [page, setPage] = useState(0);
  const [rpp, setRpp] = useState(10);
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [documentToDelete, setDocumentToDelete] = useState(null);
  const [deletingId, setDeletingId] = useState('');
  const [deleteError, setDeleteError] = useState('');
  const fileInput = useRef(null);
  const uploadQueue = usePdfUploadQueue();

  const loadDocuments = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      setDocuments(await fetchDocuments());
    } catch (err) {
      setError(err.message || 'Could not load documents from the backend.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadDocuments();
  }, [loadDocuments]);

  const handleFileChange = (event) => {
    uploadQueue.addFiles(event.target.files);
    event.target.value = '';
  };

  const handleUploadQueue = async () => {
    await uploadQueue.processQueue();
    setPage(0);
    await loadDocuments();
  };

  const handleRetry = async () => {
    await uploadQueue.retryFailed();
    setPage(0);
    await loadDocuments();
  };

  const handleDeleteDocument = async () => {
    if (!documentToDelete) return;
    setDeleteError('');
    setDeletingId(documentToDelete.id);
    try {
      await deleteDocument(documentToDelete.id);
      setDocuments((current) => current.filter((document) => document.id !== documentToDelete.id));
      setPage((current) => Math.min(current, Math.max(0, Math.ceil((rows.length - 1) / rpp) - 1)));
      setDocumentToDelete(null);
    } catch (err) {
      setDeleteError(err.message || 'Could not delete this document.');
    } finally {
      setDeletingId('');
    }
  };

  const rows = useMemo(
    () => documents.map((document) => ({
      id: document.id,
      name: document.original_name,
      source: 'Manual upload',
      type: 'PDF document',
      company: '—',
      received: formatReceived(document.uploaded_at),
      size: formatFileSize(document.size_bytes),
      status: displayDocumentStatus(document.status),
    })).filter((document) =>
      (status === 'All' || document.status === status) &&
      (source === 'All' || document.source === source) &&
      `${document.name}`.toLowerCase().includes(q.toLowerCase())),
    [documents, q, status, source]
  );

  return (
    <>
      <input
        ref={fileInput}
        type="file"
        accept=".pdf,application/pdf"
        hidden
        multiple
        onChange={handleFileChange}
      />
      <PageHeader
        title="Documents"
        subtitle="Add multiple PDFs to the upload queue and track each file as it uploads sequentially."
        actions={
          <Stack direction="row" spacing={1}>
            <Button variant="outlined" onClick={loadDocuments} disabled={loading || uploadQueue.processing}>Refresh</Button>
            <Button variant="outlined" onClick={() => fileInput.current?.click()}>Add PDFs</Button>
          </Stack>
        }
      />
      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
      {deleteError && <Alert severity="error" sx={{ mb: 2 }} onClose={() => setDeleteError('')}>{deleteError}</Alert>}
      <PdfUploadQueue
        {...uploadQueue}
        onUpload={handleUploadQueue}
        onRetry={handleRetry}
        onClearCompleted={uploadQueue.clearCompleted}
        onClearQueue={uploadQueue.clearQueue}
      />
      <Stack direction={{ xs: 'column', md: 'row' }} spacing={1.5} sx={{ mb: 2 }}>
        <TextField size="small" fullWidth placeholder="Search by file name" value={q}
          onChange={(e) => { setQ(e.target.value); setPage(0); }}
          slotProps={{ input: { startAdornment: <InputAdornment position="start"><SearchIcon fontSize="small" /></InputAdornment> } }}
          sx={{ bgcolor: '#fff' }} />
        <TextField select size="small" label="Source" value={source} onChange={(e) => setSource(e.target.value)} sx={{ minWidth: 160, bgcolor: '#fff' }}>
          {sources.map((s) => <MenuItem key={s} value={s}>{s}</MenuItem>)}
        </TextField>
        <TextField select size="small" label="Status" value={status} onChange={(e) => setStatus(e.target.value)} sx={{ minWidth: 160, bgcolor: '#fff' }}>
          {statuses.map((s) => <MenuItem key={s} value={s}>{s}</MenuItem>)}
        </TextField>
      </Stack>
      <DataTable
        columns={[
          { key: 'name', label: 'Document', strong: true },
          { key: 'source', label: 'Source', render: (r) => <SourceCell source={r.source} /> },
          { key: 'type', label: 'Type' },
          { key: 'company', label: 'Company' },
          { key: 'received', label: 'Received' },
          { key: 'size', label: 'Size' },
          { key: 'status', label: 'Status', render: (r) => <StatusChip status={r.status} /> },
          {
            key: 'view',
            label: '',
            align: 'right',
            render: (r) => (
              <Tooltip title="View document">
                <IconButton
                  aria-label={`View ${r.name}`}
                  component="a"
                  href={getDocumentUrl(r.id)}
                  target="_blank"
                  rel="noreferrer"
                  size="small"
                  onClick={(event) => event.stopPropagation()}
                >
                  <VisibilityOutlinedIcon />
                </IconButton>
              </Tooltip>
            ),
          },
          {
            key: 'actions',
            label: '',
            align: 'right',
            render: (r) => (
              <Tooltip title="Delete document">
                <span>
                  <IconButton
                    aria-label={`Delete ${r.name}`}
                    color="error"
                    size="small"
                    disabled={Boolean(deletingId)}
                    onClick={(event) => {
                      event.stopPropagation();
                      setDeleteError('');
                      setDocumentToDelete(r);
                    }}
                  >
                    <DeleteOutlineIcon />
                  </IconButton>
                </span>
              </Tooltip>
            ),
          },
        ]}
        rows={loading ? [] : rows.slice(page * rpp, page * rpp + rpp)}
      />
      {loading && <Typography variant="body2" color="text.secondary" sx={{ p: 2 }}>Loading documents…</Typography>}
      {!loading && !error && rows.length === 0 && (
        <Typography variant="body2" color="text.secondary" sx={{ p: 2, textAlign: 'center' }}>No documents found.</Typography>
      )}
      <TablePagination component="div" count={rows.length} page={page} rowsPerPage={rpp}
        onPageChange={(_, p) => setPage(p)} onRowsPerPageChange={(e) => { setRpp(+e.target.value); setPage(0); }} />
      <Dialog
        open={Boolean(documentToDelete)}
        onClose={() => !deletingId && setDocumentToDelete(null)}
        fullWidth
        maxWidth="xs"
      >
        <DialogTitle>Delete document?</DialogTitle>
        <DialogContent>
          {documentToDelete ? `This will permanently delete "${documentToDelete.name}".` : ''}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDocumentToDelete(null)} disabled={Boolean(deletingId)}>Cancel</Button>
          <Button color="error" variant="contained" onClick={handleDeleteDocument} disabled={Boolean(deletingId)}>
            {deletingId ? 'Deleting…' : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
