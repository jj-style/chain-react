import { createBrowserRouter } from "react-router-dom";
import { Root, Support } from "./pages";

const Router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
  },
  {
    path: "/support-the-project",
    element: <Support />,
  },
]);

export default Router;
