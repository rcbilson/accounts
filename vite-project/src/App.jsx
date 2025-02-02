import React from "react";
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";

import ErrorPage from "./ErrorPage.jsx";
import QueryPage from "./QueryPage.jsx";
import NewPage from "./NewPage.jsx";
import HomePage from "./HomePage.jsx";

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
