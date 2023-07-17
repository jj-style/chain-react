import withLayout from "./layout";
import { InstantSearch, Hits } from "react-instantsearch-dom";
import { SearchBox, Graph } from "../components";
import { useEffect, useState } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import CloseButton from "react-bootstrap/CloseButton";
import Row from "react-bootstrap/Row";
import InputGroup from "react-bootstrap/InputGroup";
import { Shuffle } from "react-bootstrap-icons";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

import { SEARCH_CLIENT } from "../constants";

const Root = () => {
  const queryClient = useQueryClient();

  const [chain, setChain] = useState([]);
  const [newLink, setNewLink] = useState(null);
  const [start, setStart] = useState(null);
  const [end, setEnd] = useState(null);
  const [toSetRandomActor, setToSetRandomActor] = useState(null);
  const [verification, setVerification] = useState(null);

  // dev state init
  // useEffect(() => {
  //   setStart({ name: "Bruce Willis", id: 62 });
  //   setEnd({ name: "Harrison Ford", id: 3 });
  //   setChain([{ name: "Gary Oldman", id: 64 }]);
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
  const {
    isLoading: isLoadingRandomActor,
    error: errorRandomActor,
    refetch: refetchRandomActor,
  } = useQuery({
    queryKey: ["getRandomActor"],
    queryFn: async () => {
      const { data } = await axios.get(
        `${process.env.REACT_APP_SERVER_URL}${randomUrlPath}`
      );
      console.log("fetched ", data);
      toSetRandomActor && toSetRandomActor(data);
      return data;
    },
    enabled: false,
    refetchOnWindowFocus: false,
    onSettled: () => {
      setToSetRandomActor(null);
      queryClient.invalidateQueries("getRandomActor");
    },
  });

  // when toSetRandomActor changes, fetch a random actor to put into the toSet function
  useEffect(() => {
    if (toSetRandomActor !== null && refetchRandomActor !== null)
      refetchRandomActor();
  }, [toSetRandomActor, refetchRandomActor]);

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
  const { mutate: mutateVerification, isLoading: isLoadingVerification } =
    useMutation(postVerifyChain, {
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
  const doPostVerifyChain = () => {
    var x = [
      start.id,
      ...chain.filter((x) => x !== null).map((x) => x.id),
      end.id,
    ];
    mutateVerification({ chain: x });
  };

  return (
    <>
      <Row>
        <ListGroup className="d-flex justify-content-between">
          {/* START ACTOR */}
          <StartEnd
            setToSet={setToSetRandomActor}
            currentState={start}
            setState={setStart}
            searchClient={SEARCH_CLIENT}
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
            <InstantSearch indexName="actors" searchClient={SEARCH_CLIENT}>
              <SearchBox placeholder="find actor..." />
              <Hits
                hitComponent={({ hit }) => <Hit hit={hit} addHit={addHit} />}
              />
            </InstantSearch>
          )}

          {/* END ACTOR */}
          <StartEnd
            setToSet={setToSetRandomActor}
            currentState={end}
            setState={setEnd}
            searchClient={SEARCH_CLIENT}
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
            onClick={() => doPostVerifyChain()}
          >
            verify
          </Button>
        </ButtonGroup>
      </Row>
      {!isLoadingVerification && verification?.valid && start && end && (
        // TODO: add slider to configure max relationship hops
        <Graph
          query={`match p=(a:Actor{id: ${start.id}})-[:ACTED_IN*1..6]-(b:Actor{id:${end.id}}) return p`}
          start={start.id}
          end={end.id}
          chain={chain.filter((x) => x !== null).map((x) => x.id)}
        />
      )}
    </>
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
const postVerifyChain = async (data) => {
  const { data: response } = await axios.post(
    `${process.env.REACT_APP_SERVER_URL}/verifyEdges`,
    data
  );
  return response;
};

export default withLayout(Root, "Home");
