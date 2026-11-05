import { Avatar, Button, Stack } from '@mui/material';
import PageHeader from '../components/PageHeader';
import DataTable from '../components/DataTable';
import StatusChip from '../components/StatusChip';
import { users } from '../data/mock';

export default function Users() {
  return (
    <>
      <PageHeader title="Users" subtitle="Manage who can upload, review and approve documents." actions={<Button variant="contained">Invite user</Button>} />
      <DataTable
        columns={[
          { key: 'name', label: 'Name', strong: true, render: (r) => (
            <Stack direction="row" spacing={1.5} alignItems="center">
              <Avatar sx={{ width: 28, height: 28, fontSize: 12, bgcolor: 'primary.main' }}>{r.name.split(' ').map((w) => w[0]).slice(0, 2).join('')}</Avatar>
              <span>{r.name}</span>
            </Stack>
          ) },
          { key: 'email', label: 'Email' },
          { key: 'role', label: 'Role' },
          { key: 'last', label: 'Last active' },
          { key: 'status', label: 'Status', render: (r) => <StatusChip status={r.status} /> },
        ]}
        rows={users}
      />
    </>
  );
}
