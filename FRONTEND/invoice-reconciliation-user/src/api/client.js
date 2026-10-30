export const API_BASE_URL = (
  import.meta.env.VITE_API_BASE_URL || '/api'
).replace(/\/+$/, '');

export class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

export async function apiRequest(path, { token, ...options } = {}) {
  const headers = new Headers(options.headers);
  headers.set('Accept', 'application/json');

  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }
  if (options.body && !(options.body instanceof FormData) && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(`${API_BASE_URL}${path}`, { ...options, headers });
  if (!response.ok) {
    const message = (await response.text()).trim();
    throw new ApiError(message || `Request failed (${response.status})`, response.status);
  }
  if (response.status === 204) {
    return null;
  }

  return response.json();
}
