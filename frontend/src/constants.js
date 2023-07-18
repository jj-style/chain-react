import { MeiliSearch } from 'meilisearch';

export const MEILI_CLIENT = new MeiliSearch({
  host: process.env.REACT_APP_MEILISEARCH_URL,
  apiKey: process.env.REACT_APP_MEILISEARCH_API_KEY,
});