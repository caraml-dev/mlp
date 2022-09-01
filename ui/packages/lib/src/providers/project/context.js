import React, { useContext, useMemo } from "react";
import { ApplicationsContext } from "../application";
import { useMlpApi } from "../../hooks/useMlpApi";
import { useMatch } from "react-router-dom";
import urlJoin from "proper-url-join";
import { get } from "../../utils";

const projectIdKey = "lastSelectedProjectId";
const getSelectedProjectId = (projectId, projects) => {
  if (!!projectId) {
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

const Context = React.createContext({
  projects: [],
  currentProject: undefined,
  refresh: () => {}
});

export const ProjectsContextProvider = ({ children }) => {
  const [{ data: projects }, refresh] = useMlpApi(`/projects`, {}, []);

  const { currentApp = {} } = useContext(ApplicationsContext);

  const projectIdMatch = useMatch({
    path: urlJoin(currentApp.href, "/projects/:projectId"),
    caseSensitive: true,
    end: false
  });

  const currentProject = useMemo(() => {
    const selectedProjectId = getSelectedProjectId(
      get(projectIdMatch, "params.projectId"),
      projects
    );
    return (projects || []).find(p => String(p.id) === selectedProjectId);
  }, [projects, projectIdMatch]);

  return (
    <Context.Provider value={{ projects, currentProject, refresh }}>
      {children}
    </Context.Provider>
  );
};

export default Context;
