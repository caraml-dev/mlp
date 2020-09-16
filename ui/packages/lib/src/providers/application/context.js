import React from "react";
import { useMlpApi } from "../../hooks";

const ApplicationsContext = React.createContext({
  apps: []
});

export const ApplicationsContextProvider = ({ children }) => {
  const [{ data: apps }] = useMlpApi("/applications", {}, []);

  return (
    <ApplicationsContext.Provider value={{ apps }}>
      {children}
    </ApplicationsContext.Provider>
  );
};

export default ApplicationsContext;
