export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

async function ensureOk(response, action) {
  if (!response.ok) {
    const detail = await response.text();
    throw new Error(detail || `${action} failed (${response.status})`);
  }
}

export async function fetchDocuments() {
  const [pdfRes, tgRes] = await Promise.allSettled([
    fetch(`${API_BASE_URL}/pdf`),
    fetch(`${API_BASE_URL}/telegram/feed`),
  ]);

  let pdfDocs = [];
  if (pdfRes.status === 'fulfilled' && pdfRes.value.ok) {
    try {
      const data = await pdfRes.value.json();
      if (Array.isArray(data)) pdfDocs = data;
    } catch {
      // ignore
    }
  }

  let tgDocs = [];
  if (tgRes.status === 'fulfilled' && tgRes.value.ok) {
    try {
      const data = await tgRes.value.json();
      if (Array.isArray(data)) {
        tgDocs = data.map((item) => ({
          id: `tg-${item.id}`,
          original_name: item.original_filename || `Document_${item.id}.pdf`,
          size_bytes: item.size_bytes || 0,
          status: item.status?.toLowerCase() === 'received' ? 'uploaded' : (item.status?.toLowerCase() || 'uploaded'),
          uploaded_at: item.received_at,
          source: 'Documents via Telegram',
          sender: item.sender,
          sha256: item.sha256,
        }));
      }
    } catch {
      // ignore
    }
  }

  const seen = new Set();
  const result = [];

  for (const doc of tgDocs) {
    if (doc.sha256) seen.add(doc.sha256);
    seen.add(doc.original_name);
    result.push(doc);
  }

  for (const doc of pdfDocs) {
    const isTg =
      doc.source === 'telegram' ||
      doc.source === 'Telegram' ||
      doc.source === 'Documents via Telegram' ||
      (doc.original_name && doc.original_name.toLowerCase().includes('telegram'));

    const cleanName = doc.original_name ? doc.original_name.replace(/^\[telegram\]\s*/i, '') : '';

    if (doc.sha256 && seen.has(doc.sha256)) continue;
    if (seen.has(doc.original_name) || (cleanName && seen.has(cleanName))) continue;

    result.push({
      ...doc,
      original_name: cleanName || doc.original_name,
      source: isTg ? 'Documents via Telegram' : (doc.source || 'Manual upload'),
    });
  }

  result.sort((a, b) => new Date(b.uploaded_at || 0) - new Date(a.uploaded_at || 0));
  return result;
}

export async function uploadPdf(file) {
  const formData = new FormData();
  formData.append('file', file);
  const response = await fetch(`${API_BASE_URL}/pdf/submit`, {
    method: 'POST',
    body: formData,
  });
  await ensureOk(response, 'Uploading PDF');
  return response.json();
}

export async function fetchTelegramFeed() {
  try {
    const response = await fetch(`${API_BASE_URL}/telegram/feed`);
    if (!response.ok) return [];
    const items = await response.json();
    return Array.isArray(items) ? items : [];
  } catch {
    return [];
  }
}

export function displayDocumentStatus(status) {
  const labels = { uploaded: 'Uploaded', processing: 'Processing', done: 'Processed', failed: 'Failed', received: 'Received' };
  return labels[status?.toLowerCase()] ?? status ?? 'Unknown';
}
