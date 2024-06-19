import {
  AuthProvider,
  ErrorBoundary,
  Login,
  MlpApiContextProvider,
  Page404,
  Toast
} from "@caraml-dev/ui-lib";
import { EuiProvider } from "@elastic/eui";
import React from "react";
import { Route, Routes } from "react-router-dom";
import AppRoutes from "./AppRoutes";
import { PrivateLayout } from "./PrivateLayout";
import config from "./config";

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
