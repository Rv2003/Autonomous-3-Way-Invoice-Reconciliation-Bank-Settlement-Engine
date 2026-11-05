import { useEffect, useState } from 'react';
import { useOutletContext } from 'react-router-dom';
import {
  Alert, Avatar, Box, Button, Dialog, DialogActions, DialogContent, DialogTitle,
  Paper, Stack, TextField, Typography,
} from '@mui/material';
import LockResetIcon from '@mui/icons-material/LockReset';
import { changePassword, updateProfile } from '../api/auth';

const emptyPassword = { current_password: '', new_password: '', confirm_password: '' };

function profileFromUser(user) {
  return {
    full_name: user?.full_name || '',
    email: user?.email || '',
    employee_number: user?.employee_number || '',
    phone_number: user?.phone_number || '',
  };
}

export default function Profile() {
  const { user, updateUser } = useOutletContext();
  const [profile, setProfile] = useState(() => profileFromUser(user));
  const [password, setPassword] = useState(emptyPassword);
  const [profileError, setProfileError] = useState('');
  const [passwordError, setPasswordError] = useState('');
  const [profileSuccess, setProfileSuccess] = useState('');
  const [passwordSuccess, setPasswordSuccess] = useState('');
  const [passwordDialogOpen, setPasswordDialogOpen] = useState(false);
  const [isEditingProfile, setIsEditingProfile] = useState(false);
  const [isSavingProfile, setIsSavingProfile] = useState(false);
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const initials = user?.full_name?.split(/\s+/).map((part) => part[0]).slice(0, 2).join('').toUpperCase() || 'U';

  useEffect(() => {
    setProfile(profileFromUser(user));
  }, [user]);

  const handleProfileSubmit = async (event) => {
    event.preventDefault();
    setProfileError('');
    setProfileSuccess('');
    setIsSavingProfile(true);
    try {
      const updatedUser = await updateProfile(profile);
      updateUser(updatedUser);
      setProfileSuccess('Your profile has been updated.');
      setIsEditingProfile(false);
    } catch (error) {
      setProfileError(error.message);
    } finally {
      setIsSavingProfile(false);
    }
  };

  const handlePasswordSubmit = async (event) => {
    event.preventDefault();
    setPasswordError('');
    setPasswordSuccess('');
    if (password.new_password !== password.confirm_password) {
      setPasswordError('New password and confirmation do not match.');
      return;
    }

    setIsChangingPassword(true);
    try {
      await changePassword(password);
      setPassword(emptyPassword);
      setPasswordSuccess('Your password has been changed.');
      setPasswordDialogOpen(false);
    } catch (error) {
      setPasswordError(error.message);
    } finally {
      setIsChangingPassword(false);
    }
  };

  return (
    <Box sx={{ width: '100%', maxWidth: 760, mx: 'auto' }}>
      <Stack spacing={0.75} alignItems="flex-start" sx={{ mb: 3 }}>
        <Typography variant="h4" sx={{ fontSize: { xs: '1.4rem', sm: '1.625rem' } }}>Profile</Typography>
        <Typography variant="body2" color="text.secondary">
          Update your personal and account information.
        </Typography>
      </Stack>
      <Stack spacing={3}>
        <Paper sx={{ p: { xs: 3, sm: 4 }, borderRadius: 3, color: '#fff', background: 'linear-gradient(120deg, #1d4ed8 0%, #4338ca 100%)' }}>
          <Stack spacing={1.5} alignItems="center" textAlign="center">
            <Avatar sx={{ width: 64, height: 64, bgcolor: 'rgba(255,255,255,0.18)', color: '#fff', fontSize: 22, fontWeight: 700 }}>{initials}</Avatar>
            <Stack spacing={0.25} alignItems="center" sx={{ minWidth: 0 }}>
              <Typography variant="h5" fontWeight={700} noWrap>{user?.full_name}</Typography>
              <Typography variant="body2" sx={{ color: 'rgba(255,255,255,0.82)' }} noWrap>{user?.email}</Typography>
            </Stack>
          </Stack>
        </Paper>

        <Paper component="form" variant="outlined" onSubmit={handleProfileSubmit} sx={{ p: { xs: 2.5, sm: 3 }, borderRadius: 3 }}>
          <Stack spacing={0.25} sx={{ mb: 3 }}>
            <Typography variant="h6" fontWeight={700}>Personal information</Typography>
            <Typography variant="body2" color="text.secondary">Keep your contact and employee details current.</Typography>
          </Stack>
          <Stack spacing={2}>
            {profileError && <Alert severity="error">{profileError}</Alert>}
            {profileSuccess && <Alert severity="success">{profileSuccess}</Alert>}
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
              <TextField
                fullWidth
                label="Full name"
                autoComplete="name"
                required
                value={profile.full_name}
                onChange={(event) => setProfile((current) => ({ ...current, full_name: event.target.value }))}
                InputProps={{ readOnly: !isEditingProfile }}
              />
              <TextField
                fullWidth
                label="Work email"
                type="email"
                autoComplete="email"
                required
                value={profile.email}
                onChange={(event) => setProfile((current) => ({ ...current, email: event.target.value }))}
                InputProps={{ readOnly: !isEditingProfile }}
              />
            </Stack>
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
              <TextField
                fullWidth
                label="Employee number"
                autoComplete="off"
                value={profile.employee_number}
                InputProps={{ readOnly: true }}
                helperText="Contact your administrator to update this number."
              />
              <TextField
                fullWidth
                label="Phone number"
                type="tel"
                autoComplete="tel"
                value={profile.phone_number}
                onChange={(event) => setProfile((current) => ({ ...current, phone_number: event.target.value }))}
                InputProps={{ readOnly: !isEditingProfile }}
              />
            </Stack>
            <Stack direction="row" justifyContent="flex-end">
              <Button
                type={isEditingProfile ? 'submit' : 'button'}
                variant="contained"
                disabled={isSavingProfile}
                onClick={() => {
                  if (!isEditingProfile) {
                    setProfileError('');
                    setProfileSuccess('');
                    setIsEditingProfile(true);
                  }
                }}
              >
                {isSavingProfile ? 'Saving…' : isEditingProfile ? 'Save profile' : 'Edit profile'}
              </Button>
            </Stack>
          </Stack>
        </Paper>

        <Paper variant="outlined" sx={{ p: { xs: 2.5, sm: 3 }, borderRadius: 3 }}>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ xs: 'stretch', sm: 'center' }} justifyContent="space-between">
            <Stack spacing={0.25}>
              <Typography variant="h6" fontWeight={700}>Password and security</Typography>
              <Typography variant="body2" color="text.secondary">Update your password to keep your account secure.</Typography>
            </Stack>
            <Button
              variant="outlined"
              startIcon={<LockResetIcon />}
              onClick={() => {
                setPasswordError('');
                setPasswordDialogOpen(true);
              }}
            >
              Change password
            </Button>
          </Stack>
          {passwordSuccess && <Alert severity="success" sx={{ mt: 2 }}>{passwordSuccess}</Alert>}
        </Paper>
      </Stack>

      <Dialog
        open={passwordDialogOpen}
        onClose={() => !isChangingPassword && setPasswordDialogOpen(false)}
        fullWidth
        maxWidth="sm"
        PaperProps={{ sx: { borderRadius: 3, overflow: 'hidden' } }}
      >
        <Box component="form" onSubmit={handlePasswordSubmit}>
          <DialogTitle sx={{ px: 3, py: 2.5, color: '#fff', background: 'linear-gradient(120deg, #1d4ed8, #4338ca)' }}>
            <Stack direction="row" spacing={1.5} alignItems="center">
              <Avatar sx={{ bgcolor: 'rgba(255,255,255,0.16)', color: '#fff' }}><LockResetIcon /></Avatar>
              <Box>
                <Typography variant="h6" component="div" fontWeight={700}>Change password</Typography>
                <Typography variant="body2" sx={{ color: 'rgba(255,255,255,0.8)' }}>
                  Confirm your current password to continue.
                </Typography>
              </Box>
            </Stack>
          </DialogTitle>
          <DialogContent sx={{ pt: '24px !important' }}>
            <Stack spacing={2}>
              {passwordError && <Alert severity="error">{passwordError}</Alert>}
              <TextField
                label="Current password"
                type="password"
                autoComplete="current-password"
                required
                value={password.current_password}
                onChange={(event) => setPassword((current) => ({ ...current, current_password: event.target.value }))}
                autoFocus
              />
              <TextField
                label="New password"
                type="password"
                autoComplete="new-password"
                required
                value={password.new_password}
                onChange={(event) => setPassword((current) => ({ ...current, new_password: event.target.value }))}
              />
              <TextField
                label="Confirm new password"
                type="password"
                autoComplete="new-password"
                required
                value={password.confirm_password}
                onChange={(event) => setPassword((current) => ({ ...current, confirm_password: event.target.value }))}
                error={Boolean(password.confirm_password) && password.new_password !== password.confirm_password}
                helperText={Boolean(password.confirm_password) && password.new_password !== password.confirm_password ? 'Passwords do not match' : ' '}
              />
            </Stack>
          </DialogContent>
          <DialogActions sx={{ px: 3, pb: 2 }}>
            <Button onClick={() => setPasswordDialogOpen(false)} disabled={isChangingPassword}>Cancel</Button>
            <Button type="submit" variant="contained" disabled={isChangingPassword}>
              {isChangingPassword ? 'Changing…' : 'Save new password'}
            </Button>
          </DialogActions>
        </Box>
      </Dialog>
    </Box>
  );
}
