import { useState } from 'react';
import { Box, IconButton, InputAdornment, Paper, Stack, TextField, Typography } from '@mui/material';
import emailIcon from '../assets/auth/email.svg';
import lockIcon from '../assets/auth/lock.svg';
import visibilityIcon from '../assets/auth/visibility.svg';
import invoiceIcon from '../assets/auth/invoice.svg';
import purchaseOrderIcon from '../assets/auth/purchase-order.svg';
import bankTransferIcon from '../assets/auth/bank-transfer.svg';

const documents = [
  { icon: invoiceIcon, name: 'INV-2043', type: 'Supplier invoice' },
  { icon: purchaseOrderIcon, name: 'PO-8812', type: 'Purchase order' },
  { icon: bankTransferIcon, name: 'TRX-20261025', type: 'Bank transfer' },
];

export function AuthField({
  label,
  type = 'text',
  placeholder,
  error = false,
  autoComplete,
  name,
  value,
  onChange,
  required = false,
}) {
  const [visible, setVisible] = useState(false);
  const isPassword = type === 'password';

  return (
    <Box>
      <Typography sx={{ mb: { xs: 0.75, md: 0.5, lg: 0.75 }, fontSize: { xs: 14, md: 10, lg: 14 }, fontWeight: 500, color: '#212121' }}>{label}</Typography>
      <TextField
        fullWidth
        name={name}
        type={isPassword && !visible ? 'password' : 'text'}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        required={required}
        error={error}
        autoComplete={autoComplete}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <Box
                component="img"
                src={isPassword ? lockIcon : emailIcon}
                alt=""
                sx={{ width: { xs: 20, md: 14, lg: 20 }, height: { xs: 20, md: 14, lg: 20 }, display: 'block' }}
              />
            </InputAdornment>
          ),
          ...(isPassword && {
            endAdornment: (
              <InputAdornment position="end">
                <IconButton
                  aria-label={visible ? 'Hide password' : 'Show password'}
                  onClick={() => setVisible((current) => !current)}
                  edge="end"
                  size="small"
                >
                  <Box component="img" src={visibilityIcon} alt="" sx={{ width: { xs: 20, md: 14, lg: 20 }, height: { xs: 20, md: 14, lg: 20 }, display: 'block' }} />
                </IconButton>
              </InputAdornment>
            ),
          }),
        }}
        sx={{
          '& .MuiOutlinedInput-root': {
            height: { xs: 44, md: 31, lg: 44 },
            borderRadius: '6px',
            bgcolor: '#fff',
            fontSize: { xs: 14, md: 10, lg: 14 },
            '& fieldset': { borderColor: error ? '#d32f2f' : '#b8c7db' },
          },
          '& input::placeholder': { color: '#9aa5b1', opacity: 1 },
        }}
      />
    </Box>
  );
}

export default function AuthLayout({ children, isRegister = false }) {
  const headline = isRegister
    ? 'Bring every financial document into one clear picture.'
    : 'Every invoice matched to its order and its payment.';
  const description = isRegister
    ? 'Create your workspace to automate invoice matching, review discrepancies and track bank settlements.'
    : 'Documents arrive by upload, WhatsApp or Gmail. Anything that does not match goes to the review queue.';

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'grid',
        gridTemplateColumns: { xs: '1fr', md: 'minmax(380px, 41.6667%) minmax(0, 1fr)' },
        bgcolor: '#f5f7fa',
      }}
    >
      <Box
        component="aside"
        sx={{
          minHeight: { xs: 'auto', md: '100vh' },
          p: { xs: 3, sm: 4, md: 5, lg: 7 },
          bgcolor: '#1d4ed8',
          color: '#fff',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'space-between',
          gap: { xs: 5, md: 3 },
        }}
      >
        <Stack direction="row" spacing={1.5} alignItems="center">
          <Box sx={{ width: { xs: 36, md: 26, lg: 36 }, height: { xs: 36, md: 26, lg: 36 }, borderRadius: 2, bgcolor: '#fff', display: 'grid', placeItems: 'center' }}>
            <Typography sx={{ color: '#2563eb', fontWeight: 700, fontSize: { xs: 14, md: 10, lg: 14 } }}>IR</Typography>
          </Box>
          <Box>
            <Typography sx={{ fontWeight: 700, fontSize: { xs: 16, md: 11, lg: 16 }, lineHeight: 1.15 }}>Invoice Recon</Typography>
            <Typography sx={{ color: '#dbeafe', fontSize: { xs: 11, md: 8, lg: 11 } }}>Settlement Engine</Typography>
          </Box>
        </Stack>

        <Box sx={{ width: '100%', maxWidth: 488, alignSelf: 'center' }}>
          <Typography sx={{ fontSize: { xs: 28, md: 24, lg: 34 }, lineHeight: { xs: 1.2, md: '30px', lg: '42px' }, fontWeight: 700, mb: { xs: 2.5, md: 1.75, lg: 2.5 } }}>
            {headline}
          </Typography>
          <Typography sx={{ color: '#dbeafe', fontSize: { xs: 16, md: 12, lg: 16 }, lineHeight: { xs: '24px', md: '18px', lg: '24px' }, mb: { xs: 2.5, md: 1.75, lg: 2.5 } }}>
            {description}
          </Typography>
          <Box
            sx={{
              p: { xs: 2.5, md: 1.75, lg: 2.5 },
              borderRadius: 3,
              bgcolor: 'rgba(255,255,255,0.1)',
              border: '1px solid rgba(255,255,255,0.22)',
              display: 'flex',
              flexDirection: 'column',
              gap: { xs: 1.5, md: 1, lg: 1.5 },
            }}
          >
            {documents.map((document) => (
              <Paper
                key={document.name}
                elevation={0}
                sx={{ px: { xs: 1.5, md: 1, lg: 1.5 }, py: { xs: 1.25, md: 0.875, lg: 1.25 }, borderRadius: 2, display: 'flex', alignItems: 'center', gap: { xs: 1.5, md: 1, lg: 1.5 } }}
              >
                <Box sx={{ width: { xs: 36, md: 26, lg: 36 }, height: { xs: 36, md: 26, lg: 36 }, borderRadius: '50%', bgcolor: '#eff6ff', display: 'grid', placeItems: 'center', flexShrink: 0 }}>
                  <Box component="img" src={document.icon} alt="" sx={{ width: { xs: 18, md: 13, lg: 18 }, height: { xs: 18, md: 13, lg: 18 } }} />
                </Box>
                <Box sx={{ minWidth: 0, flexGrow: 1 }}>
                  <Typography sx={{ fontSize: { xs: 14, md: 10, lg: 14 }, fontWeight: 500, color: '#212121', lineHeight: 1.3 }}>{document.name}</Typography>
                  <Typography sx={{ fontSize: { xs: 12, md: 8.5, lg: 12 }, color: '#637381', lineHeight: 1.3 }}>{document.type}</Typography>
                </Box>
                <Typography sx={{ fontSize: { xs: 14, md: 10, lg: 14 }, fontWeight: 500, color: '#212121', whiteSpace: 'nowrap' }}>LKR 482,500.00</Typography>
              </Paper>
            ))}
            <Stack direction="row" spacing={1.25} alignItems="center">
              <Box sx={{ bgcolor: '#e8f5e9', color: '#2e7d32', borderRadius: 4, px: { xs: 1.25, md: 0.875, lg: 1.25 }, py: { xs: 0.5, md: 0.35, lg: 0.5 } }}>
                <Typography sx={{ fontSize: { xs: 12, md: 8.5, lg: 12 }, lineHeight: 1.2, fontWeight: 500 }}>Three-way match</Typography>
              </Box>
              <Typography sx={{ fontSize: { xs: 12, md: 8.5, lg: 12 }, color: '#dbeafe' }}>Ready to settle</Typography>
            </Stack>
          </Box>
        </Box>

        <Typography sx={{ color: '#bfdbfe', fontSize: { xs: 12, md: 8.5, lg: 12 } }}>
          Autonomous 3-way invoice reconciliation and bank settlement
        </Typography>
      </Box>

      <Box
        component="main"
        sx={{
          minWidth: 0,
          minHeight: { xs: 'auto', md: '100vh' },
          p: { xs: 2.5, sm: 4, md: 3.5, lg: 5 },
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
        }}
      >
        <Box sx={{ width: '100%', maxWidth: { xs: 440, md: 313, lg: 440 } }}>
          {children}
        </Box>
      </Box>
    </Box>
  );
}
