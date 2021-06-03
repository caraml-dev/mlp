import React, { useContext, useMemo } from "react";
import ProjectsContext from "./projectsContext";

const CurrentProjectContext = React.createContext({});
const projectIdKey = "lastSelectedProjectId";

const getSelectedProjectId = (projectId, projects) => {
  if (projectId !== undefined) {
    localStorage.setItem(projectIdKey, projectId);
    return projectId;
  }

  let lastSelectedProjectId = localStorage.getItem(projectIdKey);

  if (localStorage.getItem(projectIdKey) === null) {
    if (projects.length > 0) {
      const selectedProject = projects.sort((a, b) =>
        a.name.localeCompare(b.name)
      )[0];
      lastSelectedProjectId = selectedProject.id;
    }
  }
  if (lastSelectedProjectId !== null) {
    localStorage.setItem(projectIdKey, lastSelectedProjectId);
  }
  return lastSelectedProjectId;
};

export const CurrentProjectContextProvider = ({ projectId, children }) => {
  const { projects, refresh } = useContext(ProjectsContext);
  projectId = getSelectedProjectId(projectId, projects);
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
