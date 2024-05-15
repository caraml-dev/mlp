const getEnv = env => {
  return window.env && window.env[env] ? window.env[env] : process.env[env];
};

export const sentryConfig = {
  dsn: getEnv("REACT_APP_SENTRY_DSN"),
  environment: getEnv("REACT_APP_ENVIRONMENT")
};

const config = {
  TIMEOUT: 5000,
  API: getEnv("REACT_APP_API_URL"),
  OAUTH_CLIENT_ID: getEnv("REACT_APP_OAUTH_CLIENT_ID"),
  USE_MOCK_DATA: false,
  STREAMS: getEnv("REACT_APP_STREAMS") || {},
  DOC_LINKS: getEnv("REACT_APP_DOC_LINKS") || [
    {
      href:
        "https://github.com/caraml-dev/merlin/blob/main/docs/getting-started/README.md",
      label: "Merlin User Guide"
    },
    {
      href: "https://github.com/caraml-dev/turing",
      label: "Turing User Guide"
    },
    {
      href: "https://docs.feast.dev/user-guide/overview",
      label: "Feast User Guide"
    }
  ],
  MAX_AUTHZ_CACHE_EXPIRY_MINUTES: parseInt(
    getEnv("REACT_APP_MAX_AUTHZ_CACHE_EXPIRY_MINUTES") || "0"
  ),

  CLOCKWORK_UI_HOMEPAGE: getEnv("REACT_APP_CLOCKWORK_UI_HOMEPAGE"),
  KUBEFLOW_UI_HOMEPAGE: getEnv("REACT_APP_KUBEFLOW_UI_HOMEPAGE"),
  ALLOW_CUSTOM_STREAM: getEnv("REACT_APP_ALLOW_CUSTOM_STREAM") || true,
  ALLOW_CUSTOM_TEAM: getEnv("REACT_APP_ALLOW_CUSTOM_TEAM") || true
};

export default config;
