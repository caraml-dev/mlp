import React, { useContext, useMemo } from "react";
import ProjectsContext from "./projectsContext";

const CurrentProjectContext = React.createContext({});

export const CurrentProjectContextProvider = ({ projectId, children }) => {
  const { projects, refresh } = useContext(ProjectsContext);

  const project = useMemo(() => {
    return !!projects
      ? projects.find(p => String(p.id) === projectId)
      : undefined;
  }, [projects, projectId]);

  return (
    <CurrentProjectContext.Provider value={{ projectId, project, refresh }}>
      {children}
    </CurrentProjectContext.Provider>
  );
};

export default CurrentProjectContext;
