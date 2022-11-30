import { useContext } from "react";
import { ApplicationsContext, AuthContext, useApi } from "@gojek/mlp-ui";
import config from "../config";

export const useTuringApi = (
  endpoint,
  options,
  result,
  callImmediately = true
) => {
  const authCtx = useContext(AuthContext);
  const { apps } = useContext(ApplicationsContext);

  return useApi(
    endpoint,
    {
      baseApiUrl: apps.find(app => app.name === "Turing")?.config?.api,
      timeout: config.TIMEOUT,
      useMockData: config.USE_MOCK_DATA,
      ...options
    },
    authCtx,
    result,
    apps.some(app => app.name === "Turing") && callImmediately
  );
};
