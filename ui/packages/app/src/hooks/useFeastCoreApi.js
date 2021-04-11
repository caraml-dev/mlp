import { useContext } from "react";
import { AuthContext, useApi } from "@gojek/mlp-ui";
import config from "../config";

export const useFeastCoreApi = (
  endpoint,
  options,
  result,
  callImmediately = true
) => {
  const authCtx = useContext(AuthContext);

  /* Use undefined for authCtx so that the authorization header passed in the options
   * will be used instead of being overwritten. Ref: https://github.com/gojek/mlp/blob/main/ui/packages/lib/src/utils/fetchJson.js#L39*/

  return useApi(
    endpoint,
    {
      baseApiUrl: config.FEAST_CORE_API,
      timeout: config.TIMEOUT,
      useMockData: config.USE_MOCK_DATA,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authCtx.state.idToken}`
      },
      ...options
    },
    undefined,
    result,
    callImmediately
  );
};
