import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import * as Sentry from "@sentry/browser";
import * as serviceWorker from "./serviceWorker";
import { sentryConfig } from "./config";

// Styles
import "@elastic/eui/dist/eui_theme_light.css";
import "@gojek/mlp-ui/dist/index.css";

Sentry.init(sentryConfig);

ReactDOM.render(<App />, document.getElementById("root"));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
