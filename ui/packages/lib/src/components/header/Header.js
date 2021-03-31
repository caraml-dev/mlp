import React, { useContext } from "react";
import {
  EuiContextMenuItem,
  EuiHeader,
  EuiHeaderLogo,
  EuiHeaderLink,
  EuiHeaderSection,
  EuiHeaderSectionItem
} from "@elastic/eui";
import AuthContext from "../../auth/context";
import ProjectsContext from "../../providers/project/projectsContext";
import CurrentProjectContext from "../../providers/project/currentProjectContext";
import { ProjectsDropdown } from "../projects_dropdown";
import { Breadcrumbs } from "../breadcrumbs";
import { HeaderUserMenu } from "./HeaderUserMenu";
import { slugify } from "../../utils";
import { NavDrawer } from "../nav_drawer";
import { navigate } from "@reach/router";

export const Header = ({
  homeUrl = "/",
  onProjectSelect,
  userMenuItems,
  helpLink,
  appLinks
}) => {
  const { state, onLogout } = useContext(AuthContext);
  const { projects } = useContext(ProjectsContext);
  const { projectId } = useContext(CurrentProjectContext);

  return (
    <EuiHeader position="fixed">
      <EuiHeaderSection grow={false}>
        {projectId && (
          <EuiHeaderSectionItem>
            <NavDrawer homeUrl={homeUrl} />
          </EuiHeaderSectionItem>
        )}
        <EuiHeaderSectionItem border="right">
          <EuiHeaderLogo
            iconType="graphApp"
            size="m"
            onClick={() =>
              projectId ? onProjectSelect(projectId) : navigate(homeUrl)
            }
            aria-label="Machine Learning Platform"
            title="Machine Learning Platform"
          />
        </EuiHeaderSectionItem>
        <ProjectsDropdown
          projects={projects}
          onProjectSelect={onProjectSelect}
        />
      </EuiHeaderSection>

      <EuiHeaderSection
        grow={true}
        className="euiBreadcrumbs euiHeaderBreadcrumbs"
        style={{ marginRight: 0 }}>
        <Breadcrumbs />
      </EuiHeaderSection>

      <EuiHeaderSection side="right">
        <EuiHeaderSectionItem>
          <HeaderUserMenu profileObj={state.profileObj} logout={onLogout}>
            <EuiContextMenuItem key="settings" icon="gear" href="/settings">
              Settings
            </EuiContextMenuItem>
            {userMenuItems &&
              userMenuItems.map(item => (
                <EuiContextMenuItem key={slugify(item.label)} {...item}>
                  {item.label}
                </EuiContextMenuItem>
              ))}
          </HeaderUserMenu>
        </EuiHeaderSectionItem>

        {helpLink && (
          <EuiHeaderSectionItem>
            <EuiHeaderLink iconType="help" href={helpLink} target="_blank">
              Help
            </EuiHeaderLink>
          </EuiHeaderSectionItem>
        )}
      </EuiHeaderSection>
    </EuiHeader>
  );
};
