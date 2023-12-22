import withLayout from "./layout";
import { Graph, loadOptions, Summary, Results } from "../components";
import { useEffect, useState } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import Button from "react-bootstrap/Button";
import Collapse from "react-bootstrap/Collapse";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Row from "react-bootstrap/Row";
import InputGroup from "react-bootstrap/InputGroup";
import Modal from "react-bootstrap/Modal";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Tooltip from "react-bootstrap/Tooltip";
import Spinner from "react-bootstrap/Spinner";
import { Trash, ArrowRightShort, ArrowDownShort } from "react-bootstrap-icons";
import Slider from "rc-slider";
import "rc-slider/assets/index.css";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

import AsyncSelect from "react-select/async";
import { copyToClipboard } from "../utils";

const Daily = () => {
  const queryClient = useQueryClient();

  const [chain, setChain] = useState([null]);
  const [start, setStart] = useState(null);
  const [end, setEnd] = useState(null);
  const [verification, setVerification] = useState(null);
  const [graphLength, setGraphLength] = useState(4);
  const [graphData, setGraphData] = useState(null);
  const [isLoadingGraphData, setIsLoadingGraphData] = useState(false);
  const [openHelp, setOpenHelp] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);
  const handleCloseSuccess = () => setShowSuccess(false);

  const { isLoading, isFetching, error, refetch } = useQuery(
    ["dailyGame"],
    async () => {
      const response = await fetch(
        `${process.env.REACT_APP_SERVER_URL}/api/managedGame`
      );
      if (!response.ok) {
        console.log(response?.data?.error);
        alert(`${response?.status}: ${response?.data?.error}`);
      } else {
        let data = await response.json();
        setStart(data?.start);
        setEnd(data?.end);
        return data;
      }
    },
    {
      refetchOnWindowFocus: false,
      enabled: false,
    }
  );

  useEffect(() => {
    refetch();
  }, [refetch]);

  // when any change made to chain, remove verification data
  useEffect(() => {
    setVerification(null);
    setGraphLength(4);
  }, [start, end, chain]);

  // when verification valid, fetch graph
  useEffect(() => {
    if (verification?.valid && start && end) {
      setShowSuccess(true);
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
  }, [
    verification,
    start,
    end,
    chain,
    graphLength,
    setIsLoadingGraphData,
    setShowSuccess,
  ]);

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

  if (isLoading || isFetching) {
    return (
      <div className="d-flex justify-content-center">
        <Spinner className="my-5" />
      </div>
    );
  }

  if (error !== null) {
    return (
      <div className="d-flex justify-content-center">
        <p>error loading game: {error}</p>
      </div>
    );
  }

  return (
    <>
      <Modal show={showSuccess} onHide={handleCloseSuccess}>
        <Modal.Header closeButton>
          <Modal.Title>Congratulations</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          You completed today's challenge with a chain of {chain.length + 2}!
        </Modal.Body>
        <Modal.Footer>
          <Button variant="primary" onClick={handleCloseSuccess}>
            Close and browse solutions
          </Button>
          <OverlayTrigger
            placement="right"
            delay={{ show: 250, hide: 400 }}
            overlay={(props) => (
              <Tooltip {...props}>Copied to clipboard.</Tooltip>
            )}
          >
            <Button
              variant="success"
              onClick={() => {
                copyToClipboard(
                  `https://chain-react.xyz/daily\nI solved it with a chain of ${
                    chain.length + 2
                  }, what can you do?!`
                );
              }}
            >
              Share
            </Button>
          </OverlayTrigger>
        </Modal.Footer>
      </Modal>
      <span
        onClick={() => setOpenHelp(!openHelp)}
        aria-controls="collapse-help"
        aria-expanded={openHelp}
      >
        Instructions {(openHelp && <ArrowDownShort />) || <ArrowRightShort />}
      </span>
      <Collapse in={openHelp}>
        <div id="collapse-help">
          <ol>
            <li>
              Find a link between the start and end actors via common actors
              they have both acted with in movies.
            </li>
            <li>The shorter the path the better!</li>
          </ol>
        </div>
      </Collapse>
      <Row className="my-2 mx-1">
        <ListGroup className="d-flex justify-content-between">
          {/* START ACTOR */}
          <InputGroup.Text className={`text-white correct`}>
            {start?.name}
          </InputGroup.Text>

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
          <InputGroup.Text
            className={`text-white ${
              verification?.valid ? "correct" : "incorrect"
            }`}
          >
            {end?.name}
          </InputGroup.Text>
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
      {verification?.valid &&
        start &&
        end &&
        !isLoadingGraphData &&
        !isLoadingVerification &&
        graphData !== null && (
          <Results
            Summary={<Summary data={verification} />}
            Graph={
              <Graph
                start={start.id}
                end={end.id}
                chain={chain.filter((x) => x !== null).map((x) => x.id)}
                data={graphData}
                verification={verification}
              />
            }
          />
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
          </div>
        )}
    </>
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

export default withLayout(Daily, "Daily Game");
