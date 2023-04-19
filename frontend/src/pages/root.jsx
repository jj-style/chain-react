import withLayout from "./layout";
import { InstantSearch, Hits, Highlight } from "react-instantsearch-dom";
import { instantMeiliSearch } from "@meilisearch/instant-meilisearch";
import { SearchBox } from "../components";

const Root = () => {
  const searchClient = instantMeiliSearch(
    process.env.REACT_APP_MEILISEARCH_HOST,
    process.env.REACT_APP_MEILISEARCH_API_KEY,
    {
      placeholderSearch: false,
      primaryKey: "Id",
    }
  );

  return (
    <div id="root">
      <h1>Root</h1>
      <InstantSearch indexName="actors" searchClient={searchClient}>
        <SearchBox />
        <Hits hitComponent={Hit} />
      </InstantSearch>
    </div>
  );
};

const Hit = ({ hit }) => <Highlight attribute="Name" hit={hit} />;

export default withLayout(Root);
