import { useEffect } from "react";
import Graph from "graphology";
import { SigmaContainer, useLoadGraph, useSigma } from "@react-sigma/core";
import "@react-sigma/core/lib/react-sigma.min.css";
import { useLayoutCircular } from "@react-sigma/layout-circular";

const LoadGraph = ({ data }) => {
  const loadGraph = useLoadGraph();
  const { positions, assign } = useLayoutCircular();

  useEffect(() => {
    const graph = new Graph();

    data?.chain?.map((edge, index) => {
      if (index === 0) {
        graph.addNode(edge.src.id, {
          x: 0,
          y: 0,
          size: 15,
          label: edge.src.name,
        });
      }
      graph.addNode(edge.dest.id, {
        x: 0,
        y: 0,
        size: 15,
        label: edge.dest.name,
      });
      graph.addEdge(edge.src.id, edge.dest.id, {});
    });

    loadGraph(graph);
    assign();
  }, [loadGraph, assign]);

  return null;
};

const ChainGraph = ({ data }) => {
  return (
    <SigmaContainer style={{ height: "500px", width: "500px" }}>
      <LoadGraph data={data} />
    </SigmaContainer>
  );
};

export default ChainGraph;
