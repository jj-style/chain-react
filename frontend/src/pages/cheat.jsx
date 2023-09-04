import withLayout from "./layout";
import { Graph, StartEnd } from "../components";
import { useEffect, useState } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Row from "react-bootstrap/Row";
import Spinner from "react-bootstrap/Spinner";
import Modal from "react-bootstrap/Modal";
import Slider from "rc-slider";
import "rc-slider/assets/index.css";
import axios from "axios";

const Cheat = () => {
  const [start, setStart] = useState(null);
  const [end, setEnd] = useState(null);
  const [verification, setVerification] = useState(null);
  const [graphLength, setGraphLength] = useState(4);
  const [graphData, setGraphData] = useState(null);
  const [isLoadingGraphData, setIsLoadingGraphData] = useState(false);
  const [showModal, setShowModal] = useState(false);

  const handleCloseModalNo = () => setShowModal(false);
  const handleCloseModalYes = () => {
    setShowModal(false);
    fetchGraphData();
  };

  // when any change made to, remove verification data
  useEffect(() => {
    setVerification(null);
    setGraphLength(4);
  }, [start, end]);

  const fetchGraphData = () => {
    setIsLoadingGraphData(true);
    axios
      .post(`${process.env.REACT_APP_SERVER_URL}/api/graph`, {
        chain: [start?.id, end?.id],
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
  };

  return (
    <>
      <p>Explore paths between two actors without guessing a link.</p>
      <Modal show={showModal} onHide={handleCloseModalNo}>
        <Modal.Header closeButton>
          <Modal.Title>Warning</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          Are you sure you want to ignore the purpose of the game and jump
          straight to the answer?
        </Modal.Body>
        <Modal.Footer>
          <Button variant="success" onClick={handleCloseModalNo}>
            No
          </Button>
          <Button variant="danger" onClick={handleCloseModalYes}>
            Yes
          </Button>
        </Modal.Footer>
      </Modal>
      <Row className="my-2 mx-1">
        <ListGroup className="d-flex justify-content-between">
          {/* START ACTOR */}
          <StartEnd
            currentState={start}
            setState={setStart}
            bgVariant="success"
            placeholder="start with actor"
          />

          {/* END ACTOR */}
          <StartEnd
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
            variant="outline-primary"
            disabled={start === null || end === null}
            onClick={() => setShowModal(true)}
          >
            cheat
          </Button>
        </ButtonGroup>
      </Row>

      {isLoadingGraphData && (
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
      {!isLoadingGraphData && start && end && graphData !== null && (
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
            chain={[]}
            data={graphData}
            verification={verification}
          />
        </div>
      )}
    </>
  );
};

export default withLayout(Cheat, "Cheat");
