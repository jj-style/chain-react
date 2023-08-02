import { Bitcoin, Litecoin } from "cryptocons";
import { LightningFill, CupHot } from "react-bootstrap-icons";
import withLayout from "./layout";

const Support = () => {
  return (
    <>
      <h1>Support The Project</h1>
      <p>
        Thank you for visting this site. I hope you have got some enjoyment out
        of it.
        <br />
        If you did, please do consider supporting it in any way you can.
      </p>
      <p>
        You will notice there are no popups for cookies, no ads, no anything.
        <br />I made this to be fun, not to mine and sell your data. I care
        about making something fun and interesting, and I care about your
        privacy. As an example you can see in the{" "}
        <a
          href="https://github.com/jj-style/chain-react/blob/master/backend/src/server/routes.go#L20"
          target="_blank"
          rel="noreferrer"
        >
          server log configuration
        </a>{" "}
        and the{" "}
        <a
          href="https://github.com/jj-style/chain-react/blob/master/docker/Caddyfile#L13-L15"
          target="_blank"
          rel="noreferrer"
        >
          proxy log configuration
        </a>{" "}
        the explicit removal of your IP address and user agent. I'm not storing
        it because I don't want it, I don't care about it, and I don't need it
        to run this site.
      </p>
      <p>
        However, this site does take time to develop and maintain and costs
        money to run. There are many ways to support this project, not just
        financially, and any help you can provide is amazing.
      </p>
      <h2>Requests and Issues</h2>
      <p>
        One awesome way you can contribute is if you have any cool ideas or
        feature requests, or you found a bug that needs fixing, please raise it{" "}
        <a
          href="https://github.com/jj-style/chain-react/issues/new/choose"
          target="_blank"
          rel="noreferrer"
        >
          here
        </a>
        .
      </p>
      <h2>Contributing</h2>
      <p>
        If you are a developer or know how to code, consider getting involved in
        the project, fixing bugs or implementing the next cool feature. You can
        find the repository and how to get started at{" "}
        <a
          href="https://github.com/jj-style/chain-react"
          target="_blank"
          rel="noreferrer"
        >
          https://github.com/jj-style/chain-react
        </a>
        .
      </p>
      <h2>Donation</h2>
      <p>
        This site costs me time and money to run. I develop this in the evenings
        and weekends after work. <br />
        <em>Any</em> financial contribution you can make, no matter how small
        will go a huge way to keeping this project alive.
      </p>
      If you would like to make a donation, you can do so via:
      <ul>
        <li>
          <Litecoin /> Litecoin: {process.env.REACT_APP_LITECOIN_ADDRESS}
        </li>
        <li>
          <Bitcoin /> Bitcoin: {process.env.REACT_APP_BITCOIN_ADDRESS}
        </li>
        <li>
          <LightningFill color="gold" size={24} /> Lightning Network:{" "}
          <a
            href={process.env.REACT_APP_ALBY_ADDRESS}
            target="_blank"
            rel="noreferrer"
          >
            {process.env.REACT_APP_ALBY_ADDRESS}
          </a>
        </li>
        <li>
          <CupHot size={24} />{" "}
          <a
            href={process.env.REACT_APP_BMAC_ADDRESS}
            target="_blank"
            rel="noreferrer"
          >
            bymeacoffee
          </a>
        </li>
      </ul>
      <br className="mb-5" />
    </>
  );
};

export default withLayout(Support, "Support The Project");
