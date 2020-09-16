import { useContext } from "react";
import { useApi } from "./useApi";
import AuthContext from "../auth/context";
import MlpApiContext from "../providers/api/context";

export const useMlpApi = (
  endpoint,
  options,
  result,
  callImmediately = true
) => {
  const authCtx = useContext(AuthContext);
  const apiCtx = useContext(MlpApiContext);

  return useApi(
    endpoint,
    {
      baseApiUrl: apiCtx.mlpApiUrl,
      timeout: apiCtx.timeout,
      useMockData: apiCtx.useMockData,
      ...options
    },
    authCtx,
    result,
    callImmediately
  );
};
