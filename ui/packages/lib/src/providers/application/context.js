import React, { useMemo } from "react";
import { useMlpApi } from "../../hooks/useMlpApi";
import { useLocation } from "react-router-dom";

const ApplicationsContext = React.createContext({
  apps: [],
  currentApp: undefined
});

export const ApplicationsContextProvider = ({ children }) => {
  const location = useLocation();
  const [{ data: apps }] = useMlpApi("/v1/applications", {}, []);

  const currentApp = useMemo(
    () => apps.find(a => location.pathname.startsWith(a.href)),
    [apps, location.pathname]
  );

  return (
    <ApplicationsContext.Provider value={{ currentApp, apps }}>
      {children}
    </ApplicationsContext.Provider>
  );
};

export default ApplicationsContext;
