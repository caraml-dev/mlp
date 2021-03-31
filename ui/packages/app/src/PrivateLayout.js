import React from "react";
import { navigate } from "@reach/router";
import {
  ApplicationsContextProvider,
  CurrentProjectContext,
  CurrentProjectContextProvider,
  Header,
  ProjectsContextProvider
} from "@gojek/mlp-ui";

export const PrivateLayout = Component => {
  return props => (
    <ApplicationsContextProvider>
      <ProjectsContextProvider>
        <CurrentProjectContextProvider {...props}>
          <CurrentProjectContext.Consumer>
            {({ _ }) => (
              <Header
                appIcon="graphApp"
                onProjectSelect={projectId =>
                  navigate(`/projects/${projectId}/settings`)
                }
              />
            )}
          </CurrentProjectContext.Consumer>

          <div style={{ paddingTop: "49px" }}>
            <Component {...props} />
          </div>
        </CurrentProjectContextProvider>
      </ProjectsContextProvider>
    </ApplicationsContextProvider>
  );
};
