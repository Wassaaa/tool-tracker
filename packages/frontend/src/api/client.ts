// API Client using the modern @hey-api/openapi-ts generated code
import { client } from './generated/client.gen';

// Configure the API client - in containerized setup, /api is proxied by Caddy
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api';

client.setConfig({
  baseUrl: API_BASE_URL,
});

// Re-export everything from the generated SDK for convenience
export * from './generated';

// Re-export the configured client
export { client };
