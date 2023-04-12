import React from "react";
import {
  ApplicationsContext,
  ApplicationsContextProvider,
  Header,
  PrivateRoute,
  ProjectsContextProvider
} from "@caraml-dev/ui-lib";
import config from "./config";
import urlJoin from "proper-url-join";
import { Outlet, useNavigate } from "react-router-dom";

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
                  navigate(urlJoin(currentApp?.homepage, "projects", pId))
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
