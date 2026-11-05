import { API_BASE_URL } from './documents';

export async function authenticatedFetch(path, options = {}) {
  const headers = new Headers(options.headers);
  return fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
    credentials: 'include',
  });
}
