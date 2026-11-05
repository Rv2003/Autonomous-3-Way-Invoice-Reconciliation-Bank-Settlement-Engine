import { useEffect, useState } from 'react';
import { useOutletContext } from 'react-router-dom';
import {
  Alert, Avatar, Button, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle,
  FormControl, InputLabel, MenuItem, Select, Stack, TextField, Typography,
} from '@mui/material';
import PageHeader from '../components/PageHeader';
import DataTable from '../components/DataTable';
import { createUser, fetchUsers, updateUserRole } from '../components/api/users';

const roleLabels = { administrator: 'Administrator', clerk: 'Clerk' };

const emptyForm = { full_name: '', email: '', password: '', role: 'clerk' };

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
      >
        <form onSubmit={handleCreateUser}>
          <DialogTitle>Add user</DialogTitle>
          <DialogContent>
            <Stack spacing={2} sx={{ pt: 1 }}>
              {createError && <Alert severity="error">{createError}</Alert>}
              <TextField
                label="Full name"
                value={form.full_name}
                onChange={(event) => setForm((current) => ({ ...current, full_name: event.target.value }))}
                autoComplete="name"
                required
                autoFocus
              />
              <TextField
                label="Email"
                type="email"
                value={form.email}
                onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))}
                autoComplete="email"
                required
              />
              <TextField
                label="Initial password"
                type="password"
                value={form.password}
                onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))}
                autoComplete="new-password"
                required
              />
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
