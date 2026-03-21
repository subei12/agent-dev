import { createBrowserRouter } from "react-router-dom";

import { AppLayout } from "./layout";
import { HomePage } from "../pages/home";
import { MissionPage } from "../pages/mission";
import { RuntimeSessionPage } from "../pages/runtime-session";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <AppLayout />,
    children: [
      {
        index: true,
        element: <HomePage />
      },
      {
        path: "missions/:missionId",
        element: <MissionPage />
      },
      {
        path: "runtime/sessions/:sessionId",
        element: <RuntimeSessionPage />
      }
    ]
  }
]);
