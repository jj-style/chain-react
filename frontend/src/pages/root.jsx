import withLayout from "./layout";
import { InstantSearch, Hits } from "react-instantsearch-dom";
import { instantMeiliSearch } from "@meilisearch/instant-meilisearch";
import { SearchBox } from "../components";
import { useState } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import Button from "react-bootstrap/Button";
import Container from "react-bootstrap/esm/Container";
import CloseButton from "react-bootstrap/CloseButton";
import Row from "react-bootstrap/Row";
const Root = () => {
  const searchClient = instantMeiliSearch(
    process.env.REACT_APP_MEILISEARCH_HOST,
    process.env.REACT_APP_MEILISEARCH_API_KEY,
    {
      placeholderSearch: false,
      primaryKey: "Id",
    }
  );

  const [chain, setChain] = useState([]);
  const [newLink, setNewLink] = useState(null);

  let addHit = (hit) => {
    setNewLink(null);
    setChain((curr) => [...curr, hit]);
  };

  let removeLink = (id) => {
    setChain((curr) => curr.filter((x) => x.id !== id));
  };

  return (
    <div id="root">
      <h1>Root</h1>
      <Container>
        <Row>
          <ListGroup>
            {chain.map((link, index) => {
              return (
                <ListGroup.Item
                  key={index}
                  className="d-flex justify-content-between"
                >
                  <span>{link.name}</span>
                  <CloseButton onClick={() => removeLink(link.id)} />
                </ListGroup.Item>
              );
            })}
          </ListGroup>
          {newLink !== null && (
            <InstantSearch indexName="actors" searchClient={searchClient}>
              <SearchBox />
              <Hits
                hitComponent={({ hit }) => <Hit hit={hit} addHit={addHit} />}
              />
            </InstantSearch>
          )}
        </Row>
        <Row>
          <Button variant="primary" onClick={() => setNewLink("")}>
            +
          </Button>
        </Row>
      </Container>
    </div>
  );
};

const Hit = ({ hit, addHit }) => {
  return (
    <Button
      variant="link"
      onClick={() => addHit({ name: hit.Name, id: hit.Id })}
    >
      {hit.Name}
    </Button>
  );
};

export default withLayout(Root);
