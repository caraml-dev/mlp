import React, { useState } from "react";
import { EuiSideNav, EuiIcon, EuiPageTemplate } from "@elastic/eui";
import { slugify } from "@gojek/mlp-ui/src/utils";
import { Redirect, Router } from "@reach/router";
import UserRoleSetting from "./UserRoleSetting";
import SecretSetting from "./SecretSetting";

const sections = {
  "user-roles": {
    iconType: "user",
    name: "User Roles"
  },
  "secrets-management": {
    iconType: "lock",
    name: "Secrets Management"
  }
};

const ProjectSetting = ({ "*": section, navigate }) => {
  const [isSideNavOpenOnMobile, setSideNavOpenOnMobile] = useState(true);
  const toggleOpenOnMobile = () =>
    setSideNavOpenOnMobile(!isSideNavOpenOnMobile);

  const nav = [
    {
      name: "Project Settings",
      icon: <EuiIcon type="gear" />,
      id: "settings",
      items: Object.entries(sections).map(([id, settingsItem]) => ({
        id: slugify(id),
        name: settingsItem.name,
        onClick: () => navigate(`./${id}`),
        isSelected: slugify(id) === section
      }))
    }
  ];

  return (
    <EuiPageTemplate>
      <EuiPageTemplate.Header
        restrictWidth={false}
        iconType={(sections[section] || {}).iconType}
        pageTitle={(sections[section] || {}).name}
        bottomBorder={false}
      />
      <EuiPageTemplate.Sidebar>
        <EuiSideNav
          mobileTitle="Project Settings Menu"
          toggleOpenOnMobile={toggleOpenOnMobile}
          isOpenOnMobile={isSideNavOpenOnMobile}
          items={nav}
        />
      </EuiPageTemplate.Sidebar>
      <EuiPageTemplate.Section restrictWidth="95%">
        <Router primary={false}>
          <Redirect from="/" to="user-roles" noThrow />
          <UserRoleSetting path="/user-roles" />
          <SecretSetting path="/secrets-management" />
          <Redirect default from="any" to="/errors/404" noThrow />
        </Router>
      </EuiPageTemplate.Section>
    </EuiPageTemplate>
  );
};

export default ProjectSetting;
