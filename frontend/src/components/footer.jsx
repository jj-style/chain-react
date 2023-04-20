import MadeWithLove from "react-made-with-love";
import Container from "react-bootstrap/Container";
import Image from "react-bootstrap/Image";
import tmdbLogo from "../assets/tmdb.svg";

const Footer = () => {
  return (
    <footer className="footer mt-auto py-3 bg-light fixed-bottom">
      <Container className="d-flex justify-content-between">
        <a href="https://www.themoviedb.org/" target="_blank" rel="noreferrer">
          <Image src={tmdbLogo} width={48} />
        </a>
        <MadeWithLove by="JJ Style" emoji link="https://github.com/jj-style" />
        <a href="https://www.themoviedb.org/" target="_blank" rel="noreferrer">
          <Image src={tmdbLogo} width={48} />
        </a>
      </Container>
    </footer>
  );
};

export default Footer;
