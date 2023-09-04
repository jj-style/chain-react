import { createBrowserRouter } from "react-router-dom";
import { Root, Support, Cheat } from "./pages";

const Router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
  },
  {
    path: "/support-the-project",
    element: <Support />,
  },
  { path: "/cheat", element: <Cheat /> },
]);

export default Router;
