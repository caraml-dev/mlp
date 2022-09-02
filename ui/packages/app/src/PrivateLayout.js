import React from "react";
import {
  ApplicationsContext,
  ApplicationsContextProvider,
  Header,
  PrivateRoute,
  ProjectsContextProvider
} from "@gojek/mlp-ui";
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
                appIcon="graphApp"
                onProjectSelect={pId =>
                  navigate(urlJoin(currentApp?.href, "projects", pId))
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
