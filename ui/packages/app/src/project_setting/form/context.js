import React, { useCallback, useState } from "react";

export const ProjectFormContext = React.createContext({});

export const ProjectFormContextProvider = ({
  project: initProject,
  ...props
}) => {
  const [project, setProject] = useState(initProject);

  const setName = useCallback(
    name => {
      setProject(p => ({
        ...p,
        name: name
      }));
    },
    [setProject]
  );

  const setTeam = useCallback(
    team => {
      setProject(p => ({
        ...p,
        team: team
      }));
    },
    [setProject]
  );

  const setStream = useCallback(
    stream => {
      setProject(p => ({
        ...p,
        stream: stream,
        team: undefined
      }));
    },
    [setProject]
  );

  const setAdmin = useCallback(
    admin => {
      setProject(p => ({
        ...p,
        administrators: admin
      }));
    },
    [setProject]
  );

  const setReader = useCallback(
    reader => {
      setProject(p => ({
        ...p,
        readers: reader
      }));
    },
    [setProject]
  );

  const setLabels = useCallback(
    labels => {
      setProject(p => ({
        ...p,
        labels: labels
      }));
    },
    [setProject]
  );

  return (
    <ProjectFormContext.Provider
      value={{
        project,
        setName,
        setTeam,
        setStream,
        setAdmin,
        setReader,
        setLabels
      }}>
      {props.children}
    </ProjectFormContext.Provider>
  );
};
