import React from "react";

const MlpApiContext = React.createContext({
  mlpApiUrl: undefined,
  timeout: undefined,
  useMockData: false
});

export const MlpApiContextProvider = ({
  mlpApiUrl,
  timeout,
  useMockData,
  children
}) => (
  <MlpApiContext.Provider value={{ mlpApiUrl, timeout, useMockData }}>
    {children}
  </MlpApiContext.Provider>
);

export default MlpApiContext;
