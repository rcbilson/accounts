import React from "react";
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";

import ErrorPage from "./ErrorPage.js";
import QueryPage from "./QueryPage.js";
import NewPage from "./NewPage.js";
import HomePage from "./HomePage.js";

const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
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
