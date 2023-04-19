import { NavBar, Footer } from "../components";

const withLayout = (Page) => {
  return () => (
    <div>
      <NavBar />
      <Page />
      <Footer />
    </div>
  );
};

export default withLayout;
