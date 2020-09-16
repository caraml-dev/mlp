import React from "react";
import { useMlpApi } from "../../hooks";

const ProjectsContext = React.createContext({
  projects: [],
  refresh: () => {}
});

export const ProjectsContextProvider = ({ children }) => {
  const [{ data: projects }, refresh] = useMlpApi(`/projects`, {}, []);

  return (
    <ProjectsContext.Provider value={{ projects, refresh }}>
      {children}
    </ProjectsContext.Provider>
  );
};

export default ProjectsContext;
