import { Stack, Typography } from '@mui/material';
import EmailIcon from '@mui/icons-material/Email';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';

export function TelegramLogo({ size = 18 }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      style={{ display: 'inline-block', verticalAlign: 'middle', flexShrink: 0 }}
    >
      <circle cx="12" cy="12" r="12" fill="#0088cc" />
      <path
        d="M5.4 12.05l12.4-5.1c.57-.23 1.07.13.88.94l-2.1 9.9c-.15.7-.58.87-1.16.54l-3.23-2.38-1.56 1.5c-.17.17-.32.32-.65.32l.23-3.3 6.02-5.43c.26-.23-.06-.36-.4-.13L8.4 13.57l-3.2-.99c-.7-.22-.71-.7.15-1.03z"
        fill="#ffffff"
      />
    </svg>
  );
}

export function WhatsAppLogo({ size = 18 }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      style={{ display: 'inline-block', verticalAlign: 'middle', flexShrink: 0 }}
    >
      <circle cx="12" cy="12" r="12" fill="#25D366" />
      <path
        d="M17.5 14.3c-.2-.1-1.3-.7-1.5-.8-.2-.1-.4-.1-.5.1-.2.2-.6.8-.8.9-.1.2-.3.2-.5.1-.2-.1-1-.4-1.9-1.2-.7-.6-1.2-1.4-1.3-1.6-.1-.2 0-.4.1-.5.1-.1.2-.3.4-.4.1-.1.2-.2.2-.4 0-.1 0-.3-.1-.4-.1-.1-.5-1.3-.7-1.8-.2-.5-.4-.4-.5-.4h-.5c-.2 0-.4.1-.6.3-.2.2-.8.8-.8 2 0 1.1.8 2.3 1 2.4.1.2 1.7 2.6 4.1 3.6.6.2 1 .4 1.4.5.6.2 1.1.2 1.5.1.5-.1 1.3-.6 1.5-1.1.2-.6.2-1.1.1-1.2-.1-.1-.3-.2-.5-.3z"
        fill="#ffffff"
      />
    </svg>
  );
}

export default function SourceCell({ source }) {
  const isTelegram =
    source === 'Telegram' ||
    source === 'Documents via Telegram' ||
    (typeof source === 'string' && source.toLowerCase().includes('telegram'));

  const isWhatsApp =
    source === 'WhatsApp' ||
    (typeof source === 'string' && source.toLowerCase().includes('whatsapp'));

  const isGmail =
    source === 'Gmail' ||
    (typeof source === 'string' && source.toLowerCase().includes('gmail'));

  const label = isTelegram ? 'Documents via Telegram' : (source || 'Manual upload');

  return (
    <Stack direction="row" spacing={1} alignItems="center">
      {isTelegram && <TelegramLogo size={18} />}
      {isWhatsApp && <WhatsAppLogo size={18} />}
      {isGmail && <EmailIcon sx={{ fontSize: 18, color: '#e53935' }} />}
      {!isTelegram && !isWhatsApp && !isGmail && (
        <CloudUploadIcon sx={{ fontSize: 18, color: '#2563eb' }} />
      )}
      <Typography
        variant="body2"
        sx={{
          fontWeight: isTelegram ? 600 : 400,
          color: isTelegram ? '#0088cc' : 'inherit',
        }}
      >
        {label}
      </Typography>
    </Stack>
  );
}
