import React from "react";
import {
  AuthProvider,
  Page404,
  ErrorBoundary,
  Login,
  MlpApiContextProvider,
  Toast
} from "@gojek/mlp-ui";
import config from "./config";
import { PrivateLayout } from "./PrivateLayout";
import { EuiProvider } from "@elastic/eui";
import { Route, Routes } from "react-router-dom";
import AppRoutes from "./AppRoutes";

const App = () => (
  <EuiProvider>
    <ErrorBoundary>
      <MlpApiContextProvider
        mlpApiUrl={config.API}
        timeout={config.TIMEOUT}
        useMockData={config.USE_MOCK_DATA}>
        <AuthProvider clientId={config.OAUTH_CLIENT_ID}>
          <Routes>
            <Route path="/login" element={<Login />} />

            <Route element={<PrivateLayout />}>
              <Route path="/*" element={<AppRoutes />} />
            </Route>

            <Route path="/pages/404" element={<Page404 />} />
          </Routes>
          <Toast />
        </AuthProvider>
      </MlpApiContextProvider>
    </ErrorBoundary>
  </EuiProvider>
);

export default App;
