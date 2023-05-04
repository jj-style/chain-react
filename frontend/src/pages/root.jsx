import withLayout from "./layout";
import { InstantSearch, Hits } from "react-instantsearch-dom";
import { instantMeiliSearch } from "@meilisearch/instant-meilisearch";
import { SearchBox, ChainGraph } from "../components";
import { useEffect, useState } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Container from "react-bootstrap/esm/Container";
import CloseButton from "react-bootstrap/CloseButton";
import Row from "react-bootstrap/Row";
import InputGroup from "react-bootstrap/InputGroup";
import { Shuffle } from "react-bootstrap-icons";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

const Root = () => {
  const searchClient = instantMeiliSearch(
    process.env.REACT_APP_MEILISEARCH_HOST,
    process.env.REACT_APP_MEILISEARCH_API_KEY,
    {
      placeholderSearch: false,
      primaryKey: "id",
    }
  );

  const queryClient = useQueryClient();

  const [chain, setChain] = useState([]);
  const [newLink, setNewLink] = useState(null);
  const [start, setStart] = useState(null);
  const [end, setEnd] = useState(null);
  const [toSet, setToSet] = useState(null);
  const [verification, setVerification] = useState(null);

  // dev state init
  // useEffect(() => {
  //   setStart({ name: "Mark Hamill", id: 2 });
  //   setEnd({ name: "Harrison Ford", id: 3 });
  //   setChain([{ name: "Carrie Fisher", id: 4 }]);
  //   setVerification({
  //     valid: true,
  //     error: "",
  //     chain: [
  //       {
  //         src: {
  //           id: 2,
  //           name: "Mark Hamill",
  //           Id: 2294,
  //           Title: "Jay and Silent Bob Strike Back",
  //           CreditId: "52fe434bc3a36847f8049419",
  //           Character: "Cocknocker",
  //         },
  //         dest: {
  //           id: 4,
  //           name: "Carrie Fisher",
  //           Id: 2294,
  //           Title: "Jay and Silent Bob Strike Back",
  //           CreditId: "52fe434bc3a36847f804940d",
  //           Character: "Nun",
  //         },
  //       },
  //       {
  //         src: {
  //           id: 4,
  //           name: "Carrie Fisher",
  //           Id: 74849,
  //           Title: "The Star Wars Holiday Special",
  //           CreditId: "52fe48e1c3a368484e10fbdb",
  //           Character: "Princess Leia Organa",
  //         },
  //         dest: {
  //           id: 3,
  //           name: "Harrison Ford",
  //           Id: 74849,
  //           Title: "The Star Wars Holiday Special",
  //           CreditId: "52fe48e1c3a368484e10fbcf",
  //           Character: "Han Solo",
  //         },
  //       },
  //     ],
  //   });
  // }, []);

  // add a search hit to the chain
  let addHit = (hit) => {
    setNewLink(null);
    setChain((curr) => [...curr, hit]);
  };

  // remove a link from the chain by id
  let removeLink = (id) => {
    setChain((curr) => curr.filter((x) => x.id !== id));
  };

  // random url based on whether start/end are selected
  let randomUrlPath =
    start === null && end === null
      ? "/randomActor"
      : `/randomActorNot/${start !== null ? start.id : end.id}`;

  // get random actor hook
  const { isLoading, error, refetch, data } = useQuery({
    queryKey: ["getRandomActor"],
    queryFn: async () => {
      const { data } = await axios.get(`http://localhost:8080${randomUrlPath}`);
      console.log("fetched ", data);
      toSet && toSet(data);
      return data;
    },
    enabled: false,
    refetchOnWindowFocus: false,
    onSettled: () => {
      setToSet(null);
      queryClient.invalidateQueries("getRandomActor");
    },
  });

  // when toSet changes, fetch a random actor to put into the toSet function
  useEffect(() => {
    if (toSet !== null && refetch !== null) refetch();
  }, [toSet, refetch]);

  // when any change made to chain, remove verification data
  useEffect(() => {
    setVerification(null);
  }, [start, end, chain]);

  // get whether chain is valid to send to server for verification
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

  // hook for posting chain for verification
  const { mutate, isPostLoading } = useMutation(postChain, {
    onSuccess: (data) => {
      setVerification(data);
    },
    onError: (err) => {
      //console.log("error", err.response.data);
      setVerification(err.response.data);
    },
    onSettled: () => {
      queryClient.invalidateQueries("create");
    },
  });

  // callback to post the chain for verification
  const doPostChain = () => {
    var x = [
      start.id,
      ...chain.filter((x) => x !== null).map((x) => x.id),
      end.id,
    ];
    mutate({ chain: x });
  };

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
              placeholder="start with actor"
            />

            {/* ACTOR CHAIN */}
            {chain.map((link, index) => {
              return (
                <ListGroup.Item
                  key={index}
                  className="d-flex justify-content-between"
                  variant={
                    verification === null
                      ? null
                      : index < verification?.chain?.length
                      ? "success"
                      : "danger"
                  }
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
              bgVariant={verification?.valid ? "success" : "danger"}
              placeholder="end with actor"
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
              onClick={() => doPostChain()}
            >
              verify
            </Button>
          </ButtonGroup>
        </Row>
        {verification?.valid && (
          <>
            <Row>
              <Button onClick={() => alert("todo")}>View all</Button>
            </Row>
            <Row>
              <ChainGraph data={verification} />
            </Row>
          </>
        )}
      </Container>
    </div>
  );
};

// Component for start/end of chain
const StartEnd = ({
  setToSet,
  currentState,
  setState,
  searchClient,
  bgVariant,
  placeholder,
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
        placeholder={placeholder}
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

// Render each item returned by search
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

// helper to post the chain for verification
const postChain = async (data) => {
  const { data: response } = await axios.post(
    "http://localhost:8080/verify",
    data
  );
  return response;
};

export default withLayout(Root);
