import MadeWithLove from "react-made-with-love";
import Container from "react-bootstrap/Container";
const Footer = () => {
  return (
    <footer className="footer mt-auto py-3 bg-light fixed-bottom">
      <Container className="d-flex justify-content-center">
        <MadeWithLove by="JJ Style" emoji link="https://github.com/jj-style" />
      </Container>
    </footer>
  );
};

export default Footer;
