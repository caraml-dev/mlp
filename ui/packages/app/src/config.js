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
  TEAMS: (getEnv("REACT_APP_TEAMS") || []).map(team => team.trim()),
  STREAMS: (getEnv("REACT_APP_STREAMS") || []).map(stream => stream.trim()),
  DOC_LINKS: getEnv("REACT_APP_DOC_LINKS") || [
    {
      href:
        "https://github.com/gojek/merlin/blob/main/docs/getting-started/README.md",
      label: "Merlin User Guide"
    },
    { href: "https://github.com/gojek/turing", label: "Turing User Guide" },
    {
      href: "https://docs.feast.dev/user-guide/overview",
      label: "Feast User Guide"
    }
  ],

  FEAST_CORE_API: getEnv("REACT_APP_FEAST_CORE_API"), // ${MLP_HOST}/feast/api
  MERLIN_API: getEnv("REACT_APP_MERLIN_API"), // ${MLP_HOST}/api/merlin/v1
  TURING_API: getEnv("REACT_APP_TURING_API"), // ${MLP_HOST}/api/turing/v1

  CLOCKWORK_UI_HOMEPAGE: getEnv("REACT_APP_CLOCKWORK_UI_HOMEPAGE"),
  FEAST_UI_HOMEPAGE: getEnv("REACT_APP_FEAST_UI_HOMEPAGE") || "/feast",
  KUBEFLOW_UI_HOMEPAGE: getEnv("REACT_APP_KUBEFLOW_UI_HOMEPAGE"),
  MERLIN_UI_HOMEPAGE: getEnv("REACT_APP_MERLIN_UI_HOMEPAGE") || "/merlin",
  TURING_UI_HOMEPAGE: getEnv("REACT_APP_TURING_UI_HOMEPAGE") || "/turing"
};

export default config;
