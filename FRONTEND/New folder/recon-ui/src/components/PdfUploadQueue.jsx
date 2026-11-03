import {
  Alert, AlertTitle, Box, Button, LinearProgress, List, ListItem, ListItemText, Stack, Typography,
} from '@mui/material';
import StatusChip from './StatusChip';

export default function PdfUploadQueue({
  queue,
  processing,
  queuedCount,
  failedCount,
  retryableCount,
  completedCount,
  onUpload,
  onRetry,
  onClearCompleted,
  onClearQueue,
}) {
  if (!queue.length) return null;

  return (
    <Box sx={{ mt: 2 }}>
      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
        <Typography variant="body2" fontWeight={500}>Upload queue ({queue.length})</Typography>
        <Stack direction="row" spacing={1}>
          {retryableCount > 0 && (
            <Button size="small" onClick={onRetry} disabled={processing}>Retry failed ({retryableCount})</Button>
          )}
          {completedCount > 0 && (
            <Button size="small" onClick={onClearCompleted} disabled={processing}>Clear completed</Button>
          )}
          {failedCount > 0 && (
            <Button size="small" color="error" onClick={onClearQueue} disabled={processing}>Clear queue</Button>
          )}
        </Stack>
      </Stack>
      {processing && (
        <Alert severity="info" sx={{ mb: 1 }}>
          <AlertTitle>Uploading files in order</AlertTitle>
          Files are uploaded one at a time. Please keep this page open.
          <LinearProgress sx={{ mt: 1 }} />
        </Alert>
      )}
      <List dense disablePadding>
        {queue.map((entry) => (
          <ListItem
            key={entry.id}
            divider
            secondaryAction={<StatusChip status={entry.status} />}
            sx={{ pr: 11 }}
          >
            <ListItemText
              primary={entry.name}
              secondary={entry.message || (entry.status === 'Queued' ? 'Waiting in queue' : undefined)}
              primaryTypographyProps={{ variant: 'body2', noWrap: true }}
              secondaryTypographyProps={{
                color: entry.status === 'Failed' ? 'error.main' : 'text.secondary',
                variant: 'caption',
              }}
            />
          </ListItem>
        ))}
      </List>
      <Button
        fullWidth
        variant="contained"
        sx={{ mt: 1.5 }}
        onClick={onUpload}
        disabled={!queuedCount || processing}
      >
        {processing ? 'Uploading queue…' : `Upload ${queuedCount} queued file${queuedCount === 1 ? '' : 's'}`}
      </Button>
    </Box>
  );
}
