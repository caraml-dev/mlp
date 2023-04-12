import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import * as Sentry from "@sentry/browser";
import * as serviceWorker from "./serviceWorker";
import { sentryConfig } from "./config";
import { BrowserRouter } from "react-router-dom";

// Styles
import "@elastic/eui/dist/eui_theme_light.css";
import "@caraml-dev/ui-lib/dist/index.css";

Sentry.init(sentryConfig);

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
