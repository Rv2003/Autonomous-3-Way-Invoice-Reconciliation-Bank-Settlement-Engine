const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || '/api').replace(/\/+$/, '');

async function ensureOk(response, action) {
  if (!response.ok) {
    const detail = await response.text();
    throw new Error(detail || `${action} failed (${response.status})`);
  }
}

export async function fetchDocuments() {
  const response = await fetch(`${API_BASE_URL}/pdf`);
  await ensureOk(response, 'Loading documents');
  const documents = await response.json();
  return Array.isArray(documents) ? documents : [];
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

export function displayDocumentStatus(status) {
  const labels = { uploaded: 'Uploaded', processing: 'Processing', done: 'Processed', failed: 'Failed' };
  return labels[status?.toLowerCase()] ?? status ?? 'Unknown';
}
