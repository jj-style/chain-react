import { instantMeiliSearch } from "@meilisearch/instant-meilisearch";
import { MeiliSearch } from 'meilisearch';

export const SEARCH_CLIENT = instantMeiliSearch(
  process.env.REACT_APP_MEILISEARCH_URL,
  process.env.REACT_APP_MEILISEARCH_API_KEY,
  {
    placeholderSearch: false,
    primaryKey: "id",
  }
);


export const MEILI_CLIENT = new MeiliSearch({
  host: process.env.REACT_APP_MEILISEARCH_URL,
  apiKey: process.env.REACT_APP_MEILISEARCH_API_KEY,
});