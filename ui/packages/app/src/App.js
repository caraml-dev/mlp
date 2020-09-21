import React from "react";
import { Redirect, Router } from "@reach/router";
import {
  AuthProvider,
  Empty,
  ErrorBoundary,
  Login,
  MlpApiContextProvider,
  PrivateRoute,
  Toast
} from "@gojek/mlp-ui";
import config from "./config";
import Home from "./views/Home";
import { PrivateLayout } from "./PrivateLayout";
import ProjectSetting from "./project_setting/ProjectSetting";
import { ProjectCreation } from "./project_setting/ProjectCreation";
import { Settings } from "./views/settings/Settings";

export default () => (
  <ErrorBoundary>
    <MlpApiContextProvider
      mlpApiUrl={config.API}
      timeout={config.TIMEOUT}
      useMockData={config.USE_MOCK_DATA}>
      <AuthProvider clientId={config.OAUTH_CLIENT_ID}>
        <Router role="group">
          <Login path="/login" />

          <Redirect
            from="/projects/:projectId"
            to="/projects/:projectId/settings"
            noThrow
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
);
