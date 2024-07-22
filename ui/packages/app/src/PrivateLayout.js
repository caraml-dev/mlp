import {
  ApplicationsContext,
  ApplicationsContextProvider,
  Header,
  PrivateRoute,
  ProjectsContextProvider
} from "@caraml-dev/ui-lib";
import urlJoin from "proper-url-join";
import React from "react";
import { Outlet, useNavigate } from "react-router-dom";
import config from "./config";

export const PrivateLayout = () => {
  const navigate = useNavigate();
  return (
    <PrivateRoute>
      <ApplicationsContextProvider>
        <ProjectsContextProvider>
          <ApplicationsContext.Consumer>
            {({ currentApp }) => (
              <Header
                onProjectSelect={pId =>
                  // Return user to /projects/:projectId (MLP project landing page) if app is project agnostic
                  navigate(urlJoin(currentApp?.is_project_agnostic ? "" : currentApp?.homepage, "projects", pId))
                }
                docLinks={config.DOC_LINKS}
              />
            )}
          </ApplicationsContext.Consumer>
          <Outlet />
        </ProjectsContextProvider>
      </ApplicationsContextProvider>
    </PrivateRoute>
  );
};
