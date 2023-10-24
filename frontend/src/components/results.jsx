import Tabs from "react-bootstrap/Tabs";
import Tab from "react-bootstrap/Tabs";

const Results = ({ Summary, Graph }) => {
  return (
    <Tabs
      defaultActiveKey="results"
      id="results-tabs"
      className="mb-3"
      justify
      variant="underline"
      mountOnEnter={true}
    >
      <Tab eventKey="results" title="Results">
        {Summary}
      </Tab>
      <Tab eventKey="graph" title="Graph">
        {Graph}
      </Tab>
    </Tabs>
  );
};

export default Results;
