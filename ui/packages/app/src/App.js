import React from "react";
import { Router } from "@reach/router";
import {
  AuthProvider,
  Empty,
  ErrorBoundary,
  Login,
  MlpApiContextProvider,
  PrivateRoute,
  Toast
} from "@gojek/mlp-ui";
import { Home, Project, Settings } from "./pages";
import config from "./config";
import { PrivateLayout } from "./PrivateLayout";
import { ProjectCreation } from "./project_setting/ProjectCreation";
import ProjectSetting from "./project_setting/ProjectSetting";
import { EuiProvider } from "@elastic/eui";

const App = () => (
  <EuiProvider>
    <ErrorBoundary>
      <MlpApiContextProvider
        mlpApiUrl={config.API}
        timeout={config.TIMEOUT}
        useMockData={config.USE_MOCK_DATA}>
        <AuthProvider clientId={config.OAUTH_CLIENT_ID}>
          <Router role="group">
            <Login path="/login" />

            {/* PROJECT LANDING PAGE */}
            <PrivateRoute
              path="/projects/:projectId"
              render={PrivateLayout(Project)}
            />

            {/* PROJECT SETTING */}
            <PrivateRoute
              path="/projects/:projectId/settings/*"
              render={PrivateLayout(ProjectSetting)}
            />

            {/* LANDING */}
            <PrivateRoute path="/" render={PrivateLayout(Home)} />

            {/* New Project */}
            <PrivateRoute
              path="/projects/create"
              render={PrivateLayout(ProjectCreation)}
            />

            {/* SETTINGS */}
            <PrivateRoute path="/settings" render={PrivateLayout(Settings)} />
            <PrivateRoute
              path="/settings/:section"
              render={PrivateLayout(Settings)}
            />

            {/* DEFAULT */}
            <Empty default />
          </Router>
          <Toast />
        </AuthProvider>
      </MlpApiContextProvider>
    </ErrorBoundary>
  </EuiProvider>
);

export default App;
