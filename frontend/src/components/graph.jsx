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

import { loadGraph } from "./cypherToGraph";

const Graph = ({ start, end, chain, data, verification }) => {
  const [graph, setGraph] = useState(null);

  useEffect(() => {
    loadGraph(data, { id: "Id", labels: "Labels", type: "Type" })
      .then((g) => setGraph(g))
      .catch((err) => console.log("error constructing graph ", err));
  }, [data]);

  const pathEdges = verification?.chain
    .map((c) => [c?.src?.CreditId, c?.dest?.CreditId])
    .flat();

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
            if (d["Labels"][0] === "Actor") {
              d.label = d.name;
              d.color = "red";
              d.size = 5;
              if (d?.id === start || d?.id === end || chain.includes(d?.id)) {
                d.highlighted = true;
              }
            } else if (d["Labels"][0] === "Movie") {
              d.label = d.title;
              d.color = "blue";
              d.size = 5;
            }
            return d;
          },
          edgeReducer: (_, e) => {
            e.label = e.character;
            if (pathEdges.includes(e.id)) {
              e.color = "green";
              e.size = 1;
            }
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
