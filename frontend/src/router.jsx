import { createBrowserRouter } from "react-router-dom";
import { Root, About } from "./pages";

const Router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
  },
  {
    path: "/about",
    element: <About />,
  },
]);

export default Router;
