import { useEffect, useState } from "react";
import Graph, { MultiGraph } from "graphology";
import {
  SigmaContainer,
  useLoadGraph,
  ControlsContainer,
  ZoomControl,
  FullScreenControl,
  SearchControl,
  useRegisterEvents,
  useSigma,
} from "@react-sigma/core";
import "@react-sigma/core/lib/react-sigma.min.css";
import { useLayoutRandom } from "@react-sigma/layout-random";
import { useLayoutCircular } from "@react-sigma/layout-circular";

const LoadGraph = ({ data }) => {
  const loadGraph = useLoadGraph();
  const { positions, assign } = useLayoutRandom();

  useEffect(() => {
    const graph = new MultiGraph();

    data?.chains?.forEach((chain) => {
      chain?.map((edge, index) => {
        if (!graph.hasNode(edge.src.id)) {
          graph.addNode(edge.src.id, {
            x: 0,
            y: 0,
            size: 15,
            label: edge.src.name,
            data: edge.src,
          });
        }
        if (!graph.hasNode(edge.dest.id)) {
          graph.addNode(edge.dest.id, {
            x: 0,
            y: 0,
            size: 15,
            label: edge.dest.name,
            data: edge.dest,
          });
        }
        graph.addEdge(edge.src.id, edge.dest.id, { label: edge.src.Title });
      });
    });

    loadGraph(graph);
    assign();
  }, [loadGraph, assign, data]);

  const registerEvents = useRegisterEvents();
  const sigma = useSigma();

  let [downNode, setDownNode] = useState(null);

  useEffect(() => {
    // Register the events
    registerEvents({
      // node events
      enterNode: (e) => {
        e.preventSigmaDefault();
        let nodeData = sigma.getGraph().getNodeAttribute(e.node, "data");
        sigma.getGraph().setNodeAttribute(e.node, "label", nodeData.Character);
      },
      leaveNode: (e) => {
        e.preventSigmaDefault();
        let nodeData = sigma.getGraph().getNodeAttribute(e.node, "data");
        sigma.getGraph().setNodeAttribute(e.node, "label", nodeData.name);
        sigma.getGraph().removeNodeAttribute(e.node, "highlighted");
      },
      downNode: (e) => {
        setDownNode(e.node);
        let nodeData = sigma.getGraph().getNodeAttribute(e.node, "data");
        sigma.getGraph().setNodeAttribute(e.node, "label", nodeData.Character);
        sigma.getGraph().setNodeAttribute(e.node, "highlighted", true);
      },
      touchup: (e) => {
        if (downNode) {
          let nodeData = sigma.getGraph().getNodeAttribute(downNode, "data");
          sigma.getGraph().setNodeAttribute(downNode, "label", nodeData.name);
          sigma.getGraph().removeNodeAttribute(downNode, "highlighted");
          setDownNode(null);
        }
      },

      // edge events
    });
  }, [registerEvents, sigma, downNode]);

  return null;
};

const ChainGraph = ({ data }) => {
  return (
    <SigmaContainer
      className="w-100"
      style={{ height: "500px", width: "500px" }}
      settings={{ renderEdgeLabels: true }}
      graph={MultiGraph}
    >
      <LoadGraph data={data} />
      <ControlsContainer position={"bottom-right"}>
        <ZoomControl />
        <FullScreenControl />
      </ControlsContainer>
      <ControlsContainer position={"top-right"}>
        <SearchControl style={{ width: "200px" }} />
      </ControlsContainer>
    </SigmaContainer>
  );
};

export default ChainGraph;
