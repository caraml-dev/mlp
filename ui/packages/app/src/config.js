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
  USE_MOCK_DATA: false
};

export default {
  // Add common config values here
  ...config
};
