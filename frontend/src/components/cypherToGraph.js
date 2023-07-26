// Adapted from https://github.com/sim51/graphology-neo4j/blob/main/src/cypher-to-graph.ts

import Graph from "graphology";

export function loadGraph(data, opts = { id: "@id", labels: "@labels", type: "@type" }) {
  return new Promise((resolve, reject) => {
      const graph = new Graph({
        multi: true,
        type: "directed",
        allowSelfLoops: true,
      });

      data?.forEach(record => {
        record?.Nodes.forEach((value) => {
          try {
            pushValueInGraph(value, graph, opts);
          } catch (e) {
            reject(e);
          }
        record?.Relationships.forEach((value) => {
          try {
            pushValueInGraph(value, graph, opts);
          } catch (e) {
            reject(e);
          }
        });
      })
      });

      resolve(graph);
    });
}

function pushValueInGraph(value, graph, opts) {
  // check if it's a node
  if (isNode(value)) {
    mergeNodeInGraph(value, graph, opts);
  }
  // check if it's a Relationship
  else if (isRelationship(value)) {
    mergeRelationshipInGraph(value, graph, opts);
  }
  // check if it's a Path
  else if (isPath(value)) {
    const path = value;
    mergeNodeInGraph(path.start, graph, opts);
    path.segments.forEach((seg) => {
      mergeNodeInGraph(seg.start, graph, opts);
      mergeRelationshipInGraph(seg.relationship, graph, opts);
      mergeNodeInGraph(seg.end, graph, opts);
    });
  } else if (Array.isArray(value)) {
    value.forEach((item) => {
      pushValueInGraph(item, graph, opts);
    });
  } else if (Object.prototype.toString.call(value) === "[object Object]") {
    Object.keys(value).forEach((key) => {
      pushValueInGraph(value[key], graph, opts);
    });
  }
}

function mergeNodeInGraph(node, graph, opts) {
  // TODO: cast properties ?
  const vertex = {
    ...node.Props,
    [opts.id]: node.Id.toString(),
    [opts.labels]: node.Labels,
  };
  graph.mergeNode(node.Id.toString(), vertex);
}

function mergeRelationshipInGraph(rel, graph, opts) {
  const edge = {
    ...rel.Props,
    [opts.id]: rel.Id.toString(),
    [opts.type]: rel.Type,
  };
  graph.mergeEdgeWithKey(
    rel.Id.toString(),
    `${rel.StartId}`,
    `${rel.EndId}`,
    edge
  );
}

function isNode(object) {
  let isNode = false;
  if (
    object &&
    Object.prototype.hasOwnProperty.call(object, "ElementId") &&
    Object.prototype.hasOwnProperty.call(object, "Labels")
  ) {
    isNode = true;
  }
  return isNode;
}

function isRelationship(object) {
  let isRel = false;
  if (
    object &&
    Object.prototype.hasOwnProperty.call(object, "ElementId") &&
    Object.prototype.hasOwnProperty.call(object, "Type") &&
    Object.prototype.hasOwnProperty.call(object, "StartElementId") &&
    Object.prototype.hasOwnProperty.call(object, "EndElementId")
  ) {
    isRel = true;
  }
  return isRel;
}

function isPath(object) {
  let isPath = false;
  if (
    object &&
    Object.prototype.hasOwnProperty.call(object, "start") &&
    Object.prototype.hasOwnProperty.call(object, "end") &&
    Object.prototype.hasOwnProperty.call(object, "segments")
  ) {
    isPath = true;
  }
  return isPath;
}
