import { authenticatedFetch } from './auth';

export async function fetchGmailMessages(limit = 10) {
  const response = await authenticatedFetch(`/email/messages?limit=${encodeURIComponent(limit)}`);
  if (!response.ok) {
    const detail = (await response.text()).trim();
    throw new Error(detail || `Fetching Gmail messages failed (${response.status})`);
  }
  return response.json();
}
