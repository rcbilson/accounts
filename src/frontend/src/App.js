import React from "react";
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";

import ErrorPage from "./ErrorPage.js";
import QueryPage from "./QueryPage.js";
import NewPage from "./NewPage.js";

const router = createBrowserRouter([
  {
    path: "/",
    element: <div>Hello world!</div>,
    errorElement: <ErrorPage />,
  },
  {
    path: "/query",
    element: <QueryPage />
  },
  {
    path: "/new",
    element: <NewPage />
  },
]);

export default function App() {
  return (
    <RouterProvider router={router} />
  )
}
