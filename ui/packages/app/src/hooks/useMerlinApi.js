import { useContext } from "react";
import { AuthContext, useApi } from "@gojek/mlp-ui";
import config from "../config";

export const useMerlinApi = (
  endpoint,
  options,
  result,
  callImmediately = true
) => {
  const authCtx = useContext(AuthContext);

  return useApi(
    endpoint,
    {
      baseApiUrl: config.MERLIN_API,
      timeout: config.TIMEOUT,
      useMockData: config.USE_MOCK_DATA,
      ...options
    },
    authCtx,
    result,
    callImmediately
  );
};
