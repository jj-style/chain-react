import { useEffect, useState } from "react";
import {
  ControlsContainer,
  FullScreenControl,
  SigmaContainer,
  ZoomControl,
} from "@react-sigma/core";
import "@react-sigma/core/lib/react-sigma.min.css";

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

const Graph = ({ query, start, end }) => {
  const [graph, setGraph] = useState(null);

  useEffect(() => {
    cypherToGraph(
      { driver },
      query,
      {},
      { id: "@id", labels: "@labels", type: "@type" }
    )
      .then((g) => {
        setGraph(g);
      })
      .catch((err) => console.error("err getting graph", err));
  }, []);

  if (graph) {
    // Position nodes on a circle, then run Force Atlas 2 for a while to get proper graph layout:
    circular.assign(graph);
    const settings = forceAtlas2.inferSettings(graph);
    forceAtlas2.assign(graph, { settings, iterations: 600 });
    return (
      <SigmaContainer
        className="w-100"
        style={{ height: "500px", width: "500px" }}
        settings={{
          renderEdgeLabels: true,
          renderLabels: true,
          labelRenderedSizeThreshold: 0,
          nodeReducer: (_, d) => {
            if (d["@labels"][0] === "Actor") {
              d.label = d.name;
              d.color = "red";
              d.size = 5;
              if (start !== undefined && end !== undefined) {
                if (d?.id?.low === start || d?.id?.low === end) {
                  d.highlighted = true;
                }
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
