import { Helmet } from "react-helmet";
import { NavBar, Footer } from "../components";
import Container from "react-bootstrap/esm/Container";

const withLayout = (Page, title) => {
  return () => (
    <div className="h-100">
      <Helmet>
        <title>{`${title} | chain-react`}</title>
      </Helmet>
      <NavBar />
      <Container>
        <Page />
      </Container>
      <Footer />
    </div>
  );
};

export default withLayout;
