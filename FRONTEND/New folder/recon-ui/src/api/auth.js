import { API_BASE_URL } from './documents';

export async function authenticatedFetch(path, options = {}) {
  const headers = new Headers(options.headers);
  return fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
    credentials: 'include',
  });
}

async function readResponse(response, action) {
  if (!response.ok) {
    const detail = (await response.text()).trim();
    throw new Error(detail || `${action} failed (${response.status})`);
  }
  return response.status === 204 ? null : response.json();
}

export async function updateProfile(profile) {
  const response = await authenticatedFetch('/auth/profile', {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(profile),
  });
  return readResponse(response, 'Updating profile');
}

export async function changePassword(passwords) {
  const response = await authenticatedFetch('/auth/password', {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(passwords),
  });
  return readResponse(response, 'Changing password');
}
