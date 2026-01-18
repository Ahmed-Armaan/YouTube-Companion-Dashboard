import { createBrowserRouter } from "react-router";
import { RouterProvider } from "react-router-dom"
import Login from "./pages/login";
import Home from "./pages/videoList";
import VideoDashboard from "./pages/videoDashboard";

const Router = createBrowserRouter([
  {
    path: "/",
    element: <Login />,
  },
  {
    path: "/Home/",
    element: <Home />,
  },
  {
    path: "/video/",
    element: <VideoDashboard />,
  },
  {
    path: "*",
    element: <Home />
  }
])

function App() {
  return <RouterProvider router={Router} />
}

export default App
