import React, { useContext, useMemo } from "react";
import { ApplicationsContext } from "../application";
import { useMlpApi } from "../../hooks/useMlpApi";
import { useMatch } from "react-router-dom";
import urlJoin from "proper-url-join";

const projectIdKey = "lastSelectedProjectId";

const isInt = str => !isNaN(parseInt(str));

const getSelectedProjectId = (projectId, projects) => {
  // If projectId is not a valid ID, read last selected projectId from the local storage
  if (!isInt(projectId)) {
    projectId = localStorage.getItem(projectIdKey);
    // If local storage doesn't contain valid projectId, sort projects alphabetically
    // and select ID of the first project from the list
    if (!isInt(projectId)) {
      projectId = projects.sort((a, b) => a.name.localeCompare(b.name))?.[0]
        ?.id;
    }
  }

  localStorage.setItem(projectIdKey, projectId);
  return projectId;
};

const Context = React.createContext({
  projects: [],
  currentProject: undefined,
  refresh: () => {}
});

export const ProjectsContextProvider = ({ children }) => {
  const [{ data: projects }, refresh] = useMlpApi(`/v1/projects`, {}, []);

  const { currentApp = {} } = useContext(ApplicationsContext);

  const projectIdMatch = useMatch({
    path: urlJoin(currentApp.href, "/projects/:projectId"),
    caseSensitive: true,
    end: false
  });

  const currentProject = useMemo(() => {
    const selectedProjectId = getSelectedProjectId(
      projectIdMatch?.params?.projectId,
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
