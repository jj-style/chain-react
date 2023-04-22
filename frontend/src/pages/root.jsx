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
  const [start, setStart] = useState(null);
  const [end, setEnd] = useState(null);

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
        <Row></Row>
        <Row>
          <ListGroup className="d-flex justify-content-between">
            {/* START ACTOR */}
            {start !== null ? (
              <ListGroup.Item variant="success d-flex justify-content-between">
                <span>{start.name}</span>
                <CloseButton onClick={() => setStart(null)} />
              </ListGroup.Item>
            ) : (
              <InstantSearch indexName="actors" searchClient={searchClient}>
                <SearchBox placeholder="start with actor" />
                <Hits
                  hitComponent={({ hit }) => (
                    <Hit hit={hit} addHit={(hit) => setStart(hit)} />
                  )}
                />
              </InstantSearch>
            )}

            {/* ACTOR CHAIN */}
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

            {/* NEW LINK IN CHAIN TEXT FIELD */}
            {newLink !== null && (
              <InstantSearch indexName="actors" searchClient={searchClient}>
                <SearchBox placeholder="find actor..." />
                <Hits
                  hitComponent={({ hit }) => <Hit hit={hit} addHit={addHit} />}
                />
              </InstantSearch>
            )}

            {/* END ACTOR */}
            {end !== null ? (
              <ListGroup.Item variant="danger d-flex justify-content-between">
                <span>{end.name}</span>
                <CloseButton onClick={() => setEnd(null)} />
              </ListGroup.Item>
            ) : (
              <InstantSearch indexName="actors" searchClient={searchClient}>
                <SearchBox placeholder="end with actor" />
                <Hits
                  hitComponent={({ hit }) => (
                    <Hit hit={hit} addHit={(hit) => setEnd(hit)} />
                  )}
                />
              </InstantSearch>
            )}
          </ListGroup>
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
