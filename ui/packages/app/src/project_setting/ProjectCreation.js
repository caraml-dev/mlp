import React from "react";
import ProjectForm from "./form/ProjectForm";
import { ProjectFormContextProvider } from "./form/context";
import { Project } from "./form/project";
import { EuiPageTemplate } from "@elastic/eui";

export const ProjectCreation = () => {
  return (
    <EuiPageTemplate restrictWidth={720} panelled={false}>
      <EuiPageTemplate.Header
        bottomBorder={false}
        iconType="folderClosed"
        pageTitle="Create New Project"
      />
      <EuiPageTemplate.Section paddingSize="none">
        <ProjectFormContextProvider project={new Project()}>
          <ProjectForm />
        </ProjectFormContextProvider>
      </EuiPageTemplate.Section>
    </EuiPageTemplate>
  );
};
