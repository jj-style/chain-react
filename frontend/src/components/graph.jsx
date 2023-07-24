import { useEffect, useState } from "react";
import {
  ControlsContainer,
  FullScreenControl,
  SigmaContainer,
  ZoomControl,
} from "@react-sigma/core";
import "@react-sigma/core/lib/react-sigma.min.css";

import Spinner from 'react-bootstrap/Spinner';

import circular from "graphology-layout/circular";
import forceAtlas2 from "graphology-layout-forceatlas2";

import * as neo4j from "neo4j-driver";
import { cypherToGraph } from "graphology-neo4j";

const driver = neo4j.driver(
  process.env.REACT_APP_NEO4J_URL,
  neo4j.auth.basic(
    process.env.REACT_APP_NEO4J_USERNAME,
    process.env.REACT_APP_NEO4J_PASSWORD
  )
);

const Graph = ({ query, start, end, chain }) => {
  const [graph, setGraph] = useState(null);
  const [loading, setLoading] = useState(false);

    useEffect(() => {
    setLoading(true);
    cypherToGraph(
      { driver },
      query,
      {},
      { id: "@id", labels: "@labels", type: "@type" }
    )
      .then((g) => {
          setGraph(g);
          setLoading(false);
      })
      .catch((err) => console.error("err getting graph", err));
  }, [query]);

  if(loading) {
      return <div className="w-100 h-100 d-flex flex-column align-items-center justify-content-center my-5">
              <Spinner animation="border" role="status" variant="success" className="justify-content-center">
              <span className="visually-hidden">Loading...</span>
            </Spinner>
          </div>
  }

  if (graph) {
    // Position nodes on a circle, then run Force Atlas 2 for a while to get proper graph layout:
    circular.assign(graph);
    const settings = forceAtlas2.inferSettings(graph);
    forceAtlas2.assign(graph, { settings, iterations: 600 });
    return (
      <SigmaContainer
        className="w-100"
        style={{
          height: "500px",
          width: "500px",
          backgroundColor: "aliceblue",
        }}
        settings={{
          renderEdgeLabels: true,
          renderLabels: true,
          labelRenderedSizeThreshold: 0,
          nodeReducer: (_, d) => {
            if (d["@labels"][0] === "Actor") {
              d.label = d.name;
              d.color = "red";
              d.size = 5;
              if (
                d?.id?.low === start ||
                d?.id?.low === end ||
                chain.includes(d?.id?.low)
              ) {
                d.highlighted = true;
              }
            } else if (d["@labels"][0] === "Movie") {
              d.label = d.title;
              d.color = "blue";
              d.size = 5;
            }
            return d;
          },
          edgeReducer: (_, e) => {
            e.label = e.character;
            return e;
          },
        }}
        graph={graph}
      >
        <ControlsContainer position={"bottom-right"}>
          <ZoomControl />
          <FullScreenControl />
        </ControlsContainer>
      </SigmaContainer>
    );
  }
};

export default Graph;
