import * as Sentry from "@sentry/browser";
import React from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import { sentryConfig } from "./config";
import * as serviceWorker from "./serviceWorker";

import "@caraml-dev/ui-lib/dist/index.css";
import "@elastic/eui/dist/eui_theme_light.css";

Sentry.init(sentryConfig);

const container = document.getElementById("root");
const root = createRoot(container);

root.render(
  <BrowserRouter>
    <App />
  </BrowserRouter>,
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
