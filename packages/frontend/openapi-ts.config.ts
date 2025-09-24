import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: '../backend/docs/swagger.json',
  output: {
    path: 'src/api/generated',
    format: 'prettier',
  },
  plugins: [
    '@hey-api/typescript',
    {
      name: '@hey-api/sdk',
      client: '@hey-api/client-fetch',
    },
  ],
});
