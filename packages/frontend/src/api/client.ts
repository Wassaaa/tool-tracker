// API Client using the modern @hey-api/openapi-ts generated code
import { client } from './generated/client.gen';

// Configure the API client
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api';

client.setConfig({
  baseUrl: API_BASE_URL,
});

// Re-export everything from the generated SDK for convenience
export * from './generated';

// Re-export the configured client
export { client };
