import Container from "react-bootstrap/esm/Container";
import { Bitcoin, Litecoin } from "cryptocons";
import { LightningFill, CupHot } from "react-bootstrap-icons";
import withLayout from "./layout";

const Support = () => {
  return (
    <>
      <h1>Support The Project</h1>
      <p>
        If you enjoy using this site, please do consider supporting it in any
        way you can.
      </p>
      <p>
        You will notice there are no popups for cookies, no ads, no login, no
        anything. I made this to be fun, not to mine your data and make money.
        However, this site does take time to develop and maintain and costs
        money to run so any help you can provide would be hugely appreciated.
      </p>
      <h2>Requests and Issues</h2>
      <p>
        One awesome way you can contribute is if you have any cool
        ideas/features, or you found a bug that needs fixing, please raise it{" "}
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
        The site runs on a server which costs $5 per month. In addition, I need
        to pay for the domain registration. <br />
        Any financial contribution you can make, no matter how small will go a
        huge way to keeping this project alive.
      </p>
      If you would like to make a donation, you can do so via:
      <ul>
        <li>
          <Litecoin /> Litecoin: &lt;insert address here&gt;
        </li>
        <li>
          <Bitcoin /> Bitcoin: &lt;insert address here&gt;
        </li>
        <li>
          <LightningFill color="gold" size={24} /> Lightning Network: &lt;insert
          address here&gt;
        </li>
        <li>
          <CupHot size={24} /> <a href="#">bymeacoffee</a>
        </li>
      </ul>
      <br className="mb-5" />
    </>
  );
};

export default withLayout(Support, "Support The Project");
