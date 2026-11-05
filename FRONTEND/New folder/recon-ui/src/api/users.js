import { authenticatedFetch } from './auth';

async function readResponse(response, action) {
  if (!response.ok) {
    const detail = (await response.text()).trim();
    throw new Error(detail || `${action} failed (${response.status})`);
  }
  return response.json();
}

export async function fetchUsers() {
  const response = await authenticatedFetch('/users');
  return readResponse(response, 'Loading users');
}

export async function createUser(user) {
  const response = await authenticatedFetch('/users', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(user),
  });
  return readResponse(response, 'Adding user');
}

export async function updateUserRole(id, role) {
  const response = await authenticatedFetch(`/users/${encodeURIComponent(id)}/role`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ role }),
  });
  return readResponse(response, 'Changing user role');
}
