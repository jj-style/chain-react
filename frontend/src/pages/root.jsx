import withLayout from "./layout";
import { Graph } from "../components";
import { useEffect, useState } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Row from "react-bootstrap/Row";
import InputGroup from "react-bootstrap/InputGroup";
import Spinner from "react-bootstrap/Spinner";
import { Shuffle, Trash } from "react-bootstrap-icons";
import Slider from "rc-slider";
import "rc-slider/assets/index.css";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

import AsyncSelect from "react-select/async";

import { MEILI_CLIENT } from "../constants";

const Root = () => {
  const queryClient = useQueryClient();

  const [chain, setChain] = useState([]);
  const [start, setStart] = useState(null);
  const [end, setEnd] = useState(null);
  const [toSetRandomActor, setToSetRandomActor] = useState(null);
  const [verification, setVerification] = useState(null);
  const [graphLength, setGraphLength] = useState(4);
  const [graphData, setGraphData] = useState(null);
  const [isLoadingGraphData, setIsLoadingGraphData] = useState(false);

  // dev state init
  // useEffect(() => {
  //   setStart({ name: "Bruce Willis", id: 62 });
  //   setEnd({ name: "Harrison Ford", id: 3 });
  //   setChain([{ name: "Gary Oldman", id: 64 }]);
  // }, []);

  // random url based on whether start/end are selected
  let randomUrlPath =
    start === null && end === null
      ? "/api/randomActor"
      : `/api/randomActorNot/${start !== null ? start.id : end.id}`;

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
    setGraphLength(4);
  }, [start, end, chain]);

  // when verification valid, fetch graph
  useEffect(() => {
    if (verification?.valid && start && end) {
      setIsLoadingGraphData(true);
      axios
        .post(`${process.env.REACT_APP_SERVER_URL}/api/graph`, {
          chain: [
            start?.id,
            ...chain.filter((x) => x !== null).map((x) => x.id),
            end?.id,
          ],
          length: graphLength,
        })
        .then((data) => {
          setGraphData(data?.data?.result);
          setIsLoadingGraphData(false);
        })
        .catch((err) => {
          setIsLoadingGraphData(false);
          setGraphData(null);
          // TODO - show errors nicer
          console.log(err.response);
          alert(`${err?.response?.status}: ${err?.response?.data?.error}`);
        });
    }
  }, [verification, start, end, graphLength, setIsLoadingGraphData]);

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
      <Row className="my-2 mx-1">
        <ListGroup className="d-flex justify-content-between">
          {/* START ACTOR */}
          <StartEnd
            setToSet={setToSetRandomActor}
            currentState={start}
            setState={setStart}
            bgVariant="success"
            placeholder="start with actor"
          />

          {/* ACTOR CHAIN */}
          {chain.map((link, index) => {
            return (
              <InputGroup
                key={index}
                className="d-flex justify-content-between align-items-center my-1"
              >
                <AsyncSelect
                  loadOptions={loadOptions}
                  cacheOptions={true}
                  defaultOptions={true}
                  isClearable={true}
                  className="flex-fill"
                  placeholder="find actor..."
                  value={
                    link === null
                      ? null
                      : { label: link?.name, value: link?.id }
                  }
                  onChange={(n) =>
                    setChain((curr) => {
                      var next = [...curr];
                      next[index] =
                        n === null ? n : { id: n.value, name: n.label };
                      return next;
                    })
                  }
                  // https://stackoverflow.com/a/63898744
                  menuPortalTarget={document.body}
                  styles={{
                    menuPortal: (base) => ({ ...base, zIndex: 9999 }),
                    control: (base, props) => ({
                      ...base,
                      backgroundColor: `var(--bs-${
                        verification === null
                          ? "gray"
                          : index < verification?.chain?.length
                          ? "success"
                          : "danger"
                      })`,
                    }),
                    singleValue: (provided) => ({
                      ...provided,
                      color: "white",
                    }),
                    placeholder: (provided) => ({
                      ...provided,
                      color: "white",
                    }),
                  }}
                />
                <Button
                  variant="secondary"
                  onClick={() => {
                    setChain((curr) => {
                      var next = [...curr];
                      next.splice(index, 1);
                      return next;
                    });
                  }}
                >
                  <Trash />
                </Button>
              </InputGroup>
            );
          })}

          {/* END ACTOR */}
          <StartEnd
            setToSet={setToSetRandomActor}
            currentState={end}
            setState={setEnd}
            bgVariant={verification?.valid ? "success" : "danger"}
            placeholder="end with actor"
          />
        </ListGroup>
      </Row>
      <Row>
        <ButtonGroup>
          <Button
            disabled={isLoadingGraphData}
            variant="outline-primary"
            onClick={() => setChain((c) => [...c, null])}
          >
            +
          </Button>
          <Button
            variant="outline-success"
            disabled={!validChain || isLoadingGraphData}
            onClick={() => doPostVerifyChain()}
          >
            verify
          </Button>
        </ButtonGroup>
      </Row>

      {isLoadingGraphData && verification?.valid && (
        <div className="w-100 h-100 d-flex flex-column align-items-center justify-content-center my-5">
          <Spinner
            animation="border"
            role="status"
            variant="success"
            className="justify-content-center"
          >
            <span className="visually-hidden">Loading...</span>
          </Spinner>
        </div>
      )}
      {!isLoadingVerification &&
        !isLoadingGraphData &&
        verification?.valid &&
        start &&
        end &&
        graphData !== null && (
          <div>
            {/* <Slider
              className="mt-2 mb-4"
              value={graphLength}
              onChange={(n) => setGraphLength(n)}
              dots={true}
              min={4}
              max={10}
              step={2}
              marks={[4, 6, 8, 10].reduce(
                (prev, curr) => ({ ...prev, [curr]: curr }),
                {}
              )}
            /> */}
            <Graph
              start={start.id}
              end={end.id}
              chain={chain.filter((x) => x !== null).map((x) => x.id)}
              data={graphData}
            />
          </div>
        )}
    </>
  );
};

// Component for start/end of chain
const StartEnd = ({
  setToSet,
  currentState,
  setState,
  bgVariant,
  placeholder,
}) => {
  return (
    <InputGroup>
      <AsyncSelect
        loadOptions={loadOptions}
        cacheOptions={true}
        defaultOptions={false}
        isClearable={true}
        placeholder={placeholder}
        value={
          currentState === null
            ? null
            : { label: currentState?.name, value: currentState?.id }
        }
        onChange={(n) =>
          setState(n === null ? n : { id: n.value, name: n.label })
        }
        className="flex-fill"
        menuPortalTarget={document.body}
        styles={{
          menuPortal: (base) => ({ ...base, zIndex: 9999 }),
          control: (base, props) => ({
            ...base,
            backgroundColor: `var(--bs-${bgVariant})`,
          }),
          singleValue: (provided) => ({
            ...provided,
            color: "white",
          }),
          placeholder: (provided) => ({
            ...provided,
            color: "white",
          }),
        }}
      />
      <Button variant="secondary" onClick={() => setToSet(() => setState)}>
        <Shuffle />
      </Button>
    </InputGroup>
  );
};

// helper to post the chain for verification
const postVerifyChain = async (data) => {
  const { data: response } = await axios.post(
    `${process.env.REACT_APP_SERVER_URL}/api/verifyEdges`,
    data
  );
  return response;
};

// load selectable actor options based on search query
const loadOptions = (inputValue, callback) => {
  setTimeout(() => {
    MEILI_CLIENT.index("actors")
      .search(inputValue, { sort: ["popularity:desc"] })
      .then((resp) => {
        const data = resp.hits.map((h, _) => ({
          value: h?.id,
          label: h?.name,
        }));
        callback(data);
      });
  }, 500);
};

export default withLayout(Root, "Home");
