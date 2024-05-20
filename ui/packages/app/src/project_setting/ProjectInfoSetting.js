import React, { useContext } from "react";
import ProjectInfoForm from "./project_info/ProjectInfoForm";
import { ProjectFormContextProvider } from "./form/context";
import { ProjectsContext } from "@caraml-dev/ui-lib";
import { Project } from "./form/project";
import { EuiLoadingChart, EuiTextAlign } from "@elastic/eui";

const ProjectInfoSetting = () => {
  const { currentProject, refresh } = useContext(ProjectsContext);

  return (
    <>
      {!currentProject ? (
        <EuiTextAlign textAlign="center">
          <EuiLoadingChart size="xl" mono />
        </EuiTextAlign>
      ) : (
        <ProjectFormContextProvider project={Project.from(currentProject)}>
          <ProjectInfoForm
            originalProject={currentProject}
            fetchUpdates={refresh}
          />
        </ProjectFormContextProvider>
      )}
    </>
  );
};

export default ProjectInfoSetting;
