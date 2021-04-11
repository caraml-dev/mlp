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
  DOC_LINKS: getEnv("REACT_APP_DOC_LINKS") || [],

  FEAST_CORE_API: getEnv("REACT_APP_FEAST_CORE_API"),
  MERLIN_API: getEnv("REACT_APP_MERLIN_API"),
  TURING_API: getEnv("REACT_APP_TURING_API"),

  CLOCKWORK_UI_HOMEPAGE: getEnv("REACT_APP_CLOCKWORK_UI_HOMEPAGE"),
  FEAST_UI_HOMEPAGE: getEnv("REACT_APP_FEAST_UI_HOMEPAGE"),
  KUBEFLOW_UI_HOMEPAGE: getEnv("REACT_APP_KUBEFLOW_UI_HOMEPAGE"),
  MERLIN_UI_HOMEPAGE: getEnv("REACT_APP_MERLIN_UI_HOMEPAGE"),
  TURING_UI_HOMEPAGE: getEnv("REACT_APP_TURING_UI_HOMEPAGE")
};

export default {
  // Add common config values here
  ...config
};
