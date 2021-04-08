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
  DOC_LINKS: getEnv("REACT_APP_DOC_LINKS") || []
};

export default {
  // Add common config values here
  ...config
};
