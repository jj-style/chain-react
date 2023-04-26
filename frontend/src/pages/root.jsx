import withLayout from "./layout";
import { InstantSearch, Hits } from "react-instantsearch-dom";
import { instantMeiliSearch } from "@meilisearch/instant-meilisearch";
import { SearchBox } from "../components";
import { useEffect, useState } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Container from "react-bootstrap/esm/Container";
import CloseButton from "react-bootstrap/CloseButton";
import Row from "react-bootstrap/Row";
import InputGroup from "react-bootstrap/InputGroup";
import { Shuffle } from "react-bootstrap-icons";
import { useQuery } from "@tanstack/react-query";

const Root = () => {
  const searchClient = instantMeiliSearch(
    process.env.REACT_APP_MEILISEARCH_HOST,
    process.env.REACT_APP_MEILISEARCH_API_KEY,
    {
      placeholderSearch: false,
      primaryKey: "id",
    }
  );

  const [chain, setChain] = useState([]);
  const [newLink, setNewLink] = useState(null);
  const [start, setStart] = useState(null);
  const [end, setEnd] = useState(null);
  const [toSet, setToSet] = useState(null);

  let addHit = (hit) => {
    setNewLink(null);
    setChain((curr) => [...curr, hit]);
  };

  let removeLink = (id) => {
    setChain((curr) => curr.filter((x) => x.id !== id));
  };

  let randomUrlPath =
    start === null && end === null
      ? "/randomActor"
      : `/randomActorNot/${start !== null ? start.id : end.id}`;

  const { isLoading, error, refetch, data } = useQuery({
    queryKey: ["getRandomActor"],
    queryFn: () =>
      fetch(`http://localhost:8080${randomUrlPath}`)
        .then((res) => res.json())
        .then((data) => {
          console.log("fetched ", data);
          toSet && toSet(data);
          return data;
        }),
    enabled: false,
    refetchOnWindowFocus: false,
    onSettled: () => {
      setToSet(null);
    },
  });

  useEffect(() => {
    if (toSet !== null && refetch !== null) refetch();
  }, [toSet, refetch]);

  let validateChain = () => {
    if (
      start === null ||
      end === null ||
      chain.filter((x) => x !== null).length < 1
    ) {
      return false;
    }
    return true;
  };

  let validChain = validateChain();

  return (
    <div id="root">
      <h1>Root</h1>
      <Container>
        <Row></Row>
        <Row>
          <ListGroup className="d-flex justify-content-between">
            {/* START ACTOR */}
            <StartEnd
              setToSet={setToSet}
              currentState={start}
              setState={setStart}
              searchClient={searchClient}
              bgVariant="success"
            />

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
            <StartEnd
              setToSet={setToSet}
              currentState={end}
              setState={setEnd}
              searchClient={searchClient}
              bgVariant="danger"
            />
          </ListGroup>
        </Row>
        <Row>
          <ButtonGroup className="m-0 p-0">
            <Button variant="outline-primary" onClick={() => setNewLink("")}>
              +
            </Button>
            <Button
              variant="outline-info"
              disabled={!validChain}
              onClick={() => console.log("todo - post")}
            >
              verify
            </Button>
          </ButtonGroup>
        </Row>
      </Container>
    </div>
  );
};

const StartEnd = ({
  setToSet,
  currentState,
  setState,
  searchClient,
  bgVariant,
}) => {
  return currentState !== null ? (
    <InputGroup className="">
      <Button variant="secondary" onClick={() => setToSet(() => setState)}>
        <Shuffle />
      </Button>
      <ListGroup.Item
        variant={`${bgVariant} d-flex justify-content-between`}
        style={{ flexGrow: 1 }}
      >
        <span>{currentState.name}</span>
        <CloseButton onClick={() => setState(null)} />
      </ListGroup.Item>
    </InputGroup>
  ) : (
    <InstantSearch indexName="actors" searchClient={searchClient}>
      <SearchBox
        placeholder="end with actor"
        button={
          <Button variant="secondary" onClick={() => setToSet(() => setState)}>
            <Shuffle />
          </Button>
        }
      />
      <Hits
        hitComponent={({ hit }) => (
          <Hit hit={hit} addHit={(hit) => setState(hit)} />
        )}
      />
    </InstantSearch>
  );
};

const Hit = ({ hit, addHit }) => {
  return (
    <Button
      variant="link"
      onClick={() => addHit({ name: hit.name, id: hit.id })}
    >
      {hit.name}
    </Button>
  );
};

export default withLayout(Root);
