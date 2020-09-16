import {
  ApplicationsContextProvider,
  CurrentProjectContext,
  CurrentProjectContextProvider,
  Header,
  NavDrawer,
  ProjectsContextProvider
} from "@mlp/ui";
import { navigate } from "@reach/router";
import React from "react";

export const PrivateLayout = Component => {
  return props => (
    <ApplicationsContextProvider>
      <ProjectsContextProvider>
        <CurrentProjectContextProvider {...props}>
          <Header
            appIcon="graphApp"
            onProjectSelect={projectId =>
              navigate(`/projects/${projectId}/settings`)
            }
          />
          <CurrentProjectContext.Consumer>
            {({ projectId }) => projectId && <NavDrawer />}
          </CurrentProjectContext.Consumer>
          <Component {...props} />
        </CurrentProjectContextProvider>
      </ProjectsContextProvider>
    </ApplicationsContextProvider>
  );
};
