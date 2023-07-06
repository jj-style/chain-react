import { instantMeiliSearch } from "@meilisearch/instant-meilisearch";

export const SEARCH_CLIENT = instantMeiliSearch(
  process.env.REACT_APP_MEILISEARCH_URL,
  process.env.REACT_APP_MEILISEARCH_API_KEY,
  {
    placeholderSearch: false,
    primaryKey: "id",
  }
);

