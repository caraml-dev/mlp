import React, { Fragment } from "react";
import { EuiFlexGroup, EuiFlexItem, EuiSpacer } from "@elastic/eui";
import ProjectForm from "./form/ProjectForm";
import { ProjectFormContextProvider } from "./form/context";
import { Project } from "./form/project";

export const ProjectCreation = () => {
  return (
    <Fragment>
      <EuiFlexGroup justifyContent="spaceAround">
        <EuiFlexItem style={{ maxWidth: 720 }}>
          <EuiSpacer />
          <ProjectFormContextProvider project={new Project()}>
            <ProjectForm />
          </ProjectFormContextProvider>
        </EuiFlexItem>
      </EuiFlexGroup>
    </Fragment>
  );
};
