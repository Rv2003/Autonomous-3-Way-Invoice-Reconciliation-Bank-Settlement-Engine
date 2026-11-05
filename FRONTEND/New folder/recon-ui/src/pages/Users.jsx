import { useEffect, useState } from 'react';
import { useOutletContext } from 'react-router-dom';
import {
  Alert, Avatar, Box, Button, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle,
  FormControl, InputLabel, MenuItem, Select, Stack, TextField, Typography,
} from '@mui/material';
import PersonAddAlt1Icon from '@mui/icons-material/PersonAddAlt1';
import PageHeader from '../components/PageHeader';
import DataTable from '../components/DataTable';
import { createUser, fetchUsers, updateUserRole } from '../components/api/users';

const roleLabels = { administrator: 'Administrator', clerk: 'Clerk' };

const emptyForm = {
  full_name: '',
  email: '',
  employee_number: '',
  password: '',
  confirm_password: '',
  role: 'clerk',
};

export default function Users() {
  const { user } = useOutletContext();
  const isAdministrator = user?.role === 'administrator';
  const [users, setUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState('');
  const [roleError, setRoleError] = useState('');
  const [roleBeingUpdated, setRoleBeingUpdated] = useState('');
  const [dialogOpen, setDialogOpen] = useState(false);
  const [form, setForm] = useState(emptyForm);
  const [isCreating, setIsCreating] = useState(false);
  const [createError, setCreateError] = useState('');

  useEffect(() => {
    let active = true;
    fetchUsers()
      .then((result) => {
        if (active) setUsers(result);
      })
      .catch((error) => {
        if (active) setLoadError(error.message);
      })
      .finally(() => {
        if (active) setIsLoading(false);
      });

    return () => {
      active = false;
    };
  }, []);

  const handleRoleChange = async (target, role) => {
    setRoleError('');
    setRoleBeingUpdated(target.id);
    try {
      await updateUserRole(target.id, role);
      setUsers((current) => current.map((item) => item.id === target.id ? { ...item, role } : item));
    } catch (error) {
      setRoleError(error.message);
    } finally {
      setRoleBeingUpdated('');
    }
  };

  const handleCreateUser = async (event) => {
    event.preventDefault();
    setCreateError('');
    if (form.password !== form.confirm_password) {
      setCreateError('Passwords do not match.');
      return;
    }
    setIsCreating(true);
    try {
      const newUser = await createUser(form);
      setUsers((current) => [newUser, ...current]);
      setForm(emptyForm);
      setDialogOpen(false);
    } catch (error) {
      setCreateError(error.message);
    } finally {
      setIsCreating(false);
    }
  };

  const columns = [
    {
      key: 'full_name',
      label: 'Name',
      strong: true,
      render: (row) => (
        <Stack direction="row" spacing={1.5} alignItems="center">
          <Avatar sx={{ width: 28, height: 28, fontSize: 12, bgcolor: 'primary.main' }}>
            {row.full_name.split(/\s+/).map((part) => part[0]).slice(0, 2).join('').toUpperCase()}
          </Avatar>
          <span>{row.full_name}</span>
        </Stack>
      ),
    },
    { key: 'email', label: 'Email' },
    { key: 'employee_number', label: 'Employee no.', render: (row) => row.employee_number || '—' },
    { key: 'phone_number', label: 'Phone', render: (row) => row.phone_number || '—' },
    {
      key: 'role',
      label: 'Role',
      render: (row) => isAdministrator ? (
        <FormControl size="small" sx={{ minWidth: 160 }}>
          <InputLabel id={`role-label-${row.id}`}>Role</InputLabel>
          <Select
            labelId={`role-label-${row.id}`}
            label="Role"
            value={row.role}
            disabled={roleBeingUpdated === row.id || row.id === user.id}
            onChange={(event) => handleRoleChange(row, event.target.value)}
          >
            <MenuItem value="administrator">{roleLabels.administrator}</MenuItem>
            <MenuItem value="clerk">{roleLabels.clerk}</MenuItem>
          </Select>
        </FormControl>
      ) : roleLabels[row.role] ?? row.role,
    },
    {
      key: 'created_at',
      label: 'Added',
      render: (row) => new Date(row.created_at).toLocaleDateString(),
    },
  ];

  return (
    <>
      <PageHeader
        title="Users"
        subtitle="Manage workspace access and roles."
        actions={isAdministrator && <Button variant="contained" onClick={() => setDialogOpen(true)}>Add user</Button>}
      />
      {loadError && <Alert severity="error" sx={{ mb: 2 }}>{loadError}</Alert>}
      {roleError && <Alert severity="error" sx={{ mb: 2 }}>{roleError}</Alert>}
      {isLoading ? (
        <Stack alignItems="center" sx={{ py: 6 }}><CircularProgress size={28} /></Stack>
      ) : !loadError && users.length === 0 ? (
        <Alert severity="info">No users have been added yet.</Alert>
      ) : !loadError ? (
        <DataTable columns={columns} rows={users} />
      ) : null}

      <Dialog
        open={dialogOpen}
        onClose={() => !isCreating && setDialogOpen(false)}
        fullWidth
        maxWidth="sm"
        PaperProps={{ sx: { borderRadius: 3, overflow: 'hidden' } }}
      >
        <form onSubmit={handleCreateUser}>
          <DialogTitle sx={{ px: 3, py: 2.5, color: '#fff', background: 'linear-gradient(120deg, #1d4ed8, #4338ca)' }}>
            <Stack direction="row" spacing={1.5} alignItems="center">
              <Avatar sx={{ bgcolor: 'rgba(255,255,255,0.16)', color: '#fff' }}><PersonAddAlt1Icon /></Avatar>
              <Box>
                <Typography variant="h6" component="div" fontWeight={700}>Add a team member</Typography>
                <Typography variant="body2" sx={{ color: 'rgba(255,255,255,0.8)' }}>
                  Set up their account and workspace details.
                </Typography>
              </Box>
            </Stack>
          </DialogTitle>
          <DialogContent sx={{ pt: '24px !important' }}>
            <Stack spacing={2}>
              {createError && <Alert severity="error">{createError}</Alert>}
              <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: '1fr 1fr' }, gap: 2 }}>
                <TextField
                  label="Full name"
                  value={form.full_name}
                  onChange={(event) => setForm((current) => ({ ...current, full_name: event.target.value }))}
                  autoComplete="name"
                  required
                  autoFocus
                />
                <TextField
                  label="Work email"
                  type="email"
                  value={form.email}
                  onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))}
                  autoComplete="email"
                  required
                />
                <TextField
                  label="Employee number"
                  value={form.employee_number}
                  onChange={(event) => setForm((current) => ({ ...current, employee_number: event.target.value }))}
                  autoComplete="off"
                />
                <TextField
                  label="Initial password"
                  type="password"
                  value={form.password}
                  onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))}
                  autoComplete="new-password"
                  required
                />
                <TextField
                  label="Confirm password"
                  type="password"
                  value={form.confirm_password}
                  onChange={(event) => setForm((current) => ({ ...current, confirm_password: event.target.value }))}
                  autoComplete="new-password"
                  required
                  error={Boolean(form.confirm_password) && form.password !== form.confirm_password}
                  helperText={Boolean(form.confirm_password) && form.password !== form.confirm_password ? 'Passwords do not match' : ' '}
                />
              </Box>
              <FormControl fullWidth>
                <InputLabel id="new-user-role-label">Role</InputLabel>
                <Select
                  labelId="new-user-role-label"
                  label="Role"
                  value={form.role}
                  onChange={(event) => setForm((current) => ({ ...current, role: event.target.value }))}
                >
                  <MenuItem value="administrator">{roleLabels.administrator}</MenuItem>
                  <MenuItem value="clerk">{roleLabels.clerk}</MenuItem>
                </Select>
              </FormControl>
              <Typography variant="caption" color="text.secondary">
                Share the initial password with the user securely.
              </Typography>
            </Stack>
          </DialogContent>
          <DialogActions sx={{ px: 3, pb: 2 }}>
            <Button onClick={() => setDialogOpen(false)} disabled={isCreating}>Cancel</Button>
            <Button type="submit" variant="contained" disabled={isCreating}>
              {isCreating ? 'Adding…' : 'Add user'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </>
  );
}
