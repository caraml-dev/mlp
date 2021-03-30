import React from "react";
import { navigate } from "@reach/router";
import { EuiPage } from "@elastic/eui";
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
          <EuiPage style={{ paddingTop: "49px" }}>
            <Component {...props} />
          </EuiPage>
        </CurrentProjectContextProvider>
      </ProjectsContextProvider>
    </ApplicationsContextProvider>
  );
};
