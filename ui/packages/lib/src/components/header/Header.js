import React, { useContext } from "react";
import {
  EuiContextMenuItem,
  EuiHeader,
  EuiHeaderLink,
  EuiHeaderSection,
  EuiHeaderSectionItem,
  EuiText
} from "@elastic/eui";
import AuthContext from "../../auth/context";
import { ProjectsDropdown } from "../projects_dropdown";
import { Breadcrumbs } from "../breadcrumbs";
import { HeaderUserMenu } from "./HeaderUserMenu";
import { slugify } from "../../utils";
import { NavDrawer } from "../nav_drawer";
import "./Header.scss";

export const Header = ({
  onProjectSelect,
  userMenuItems,
  helpLink,
  docLinks,
  homepage = "/"
}) => {
  const { state, onLogout } = useContext(AuthContext);

  return (
    <EuiHeader position="fixed">
      <EuiHeaderSection grow={false}>
        <EuiHeaderSectionItem>
          <NavDrawer docLinks={docLinks} />
        </EuiHeaderSectionItem>

        <EuiHeaderSectionItem border="none">
          <a href={homepage} className="euiHeaderLogo euiHeaderLogo__title">
            <EuiText>
              <h4>Machine Learning Platform</h4>
            </EuiText>
          </a>
        </EuiHeaderSectionItem>
      </EuiHeaderSection>

      <EuiHeaderSection>
        <EuiHeaderSectionItem>
          <ProjectsDropdown onProjectSelect={onProjectSelect} />
        </EuiHeaderSectionItem>
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
