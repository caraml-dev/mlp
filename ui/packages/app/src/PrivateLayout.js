import React from "react";
import { navigate } from "@reach/router";
import {
  ApplicationsContextProvider,
  CurrentProjectContextProvider,
  Header,
  ProjectsContextProvider
} from "@gojek/mlp-ui";
import config from "./config";
import "./PrivateLayout.scss";

export const PrivateLayout = Component => {
  return props => (
    <ApplicationsContextProvider>
      <ProjectsContextProvider>
        <CurrentProjectContextProvider {...props}>
          <Header
            appIcon="graphApp"
            onProjectSelect={projectId => navigate(`/projects/${projectId}`)}
            docLinks={config.DOC_LINKS}
          />
          <div className="main-component-layout">
            <Component {...props} />
          </div>
        </CurrentProjectContextProvider>
      </ProjectsContextProvider>
    </ApplicationsContextProvider>
  );
};
