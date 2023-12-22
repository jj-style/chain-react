import { createBrowserRouter } from "react-router-dom";
import { Root, Support, Cheat, Daily } from "./pages";

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
  { path: "/daily", element: <Daily /> },
]);

export default Router;
