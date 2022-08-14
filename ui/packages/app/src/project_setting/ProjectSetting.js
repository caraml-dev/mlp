import React, { Fragment, useState } from "react";
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiIcon,
  EuiPage,
  EuiPageBody,
  EuiPageContent,
  EuiPageHeader,
  EuiPageHeaderSection,
  EuiSideNav,
  EuiTitle
} from "@elastic/eui";
import { slugify } from "@gojek/mlp-ui/src/utils";
import { Redirect, Router } from "@reach/router";
import UserRoleSetting from "./UserRoleSetting";
import SecretSetting from "./SecretSetting";

const sections = [
  {
    id: "user-roles",
    name: "User Roles"
  },
  {
    id: "secrets-management",
    name: "Secrets Management"
  }
];

const ProjectSetting = ({ "*": section, navigate }) => {
  const [isSideNavOpenOnMobile, setSideNavOpenOnMobile] = useState(true);
  const toggleOpenOnMobile = () =>
    setSideNavOpenOnMobile(!isSideNavOpenOnMobile);

  const nav = [
    {
      name: "Settings",
      id: "settings",
      items: sections.map(settingsItem => ({
        id: settingsItem.id,
        name: settingsItem.name,
        onClick: () => navigate(`./${settingsItem.id}`),
        isSelected: slugify(settingsItem.id) === section
      }))
    }
  ];

  return (
    <div />
    // <EuiPage>
    //   <EuiPageBody>
    //     <EuiPageHeader>
    //       <EuiPageHeaderSection>
    //         <EuiTitle size="l">
    //           <h1>
    //             <EuiIcon type="gear" size="xl" /> Project Settings
    //           </h1>
    //         </EuiTitle>
    //       </EuiPageHeaderSection>
    //     </EuiPageHeader>
    //     <EuiFlexGroup>
    //       <EuiFlexItem grow={1}>
    //         <EuiSideNav
    //           mobileTitle="Project Settings Menu"
    //           toggleOpenOnMobile={toggleOpenOnMobile}
    //           isOpenOnMobile={isSideNavOpenOnMobile}
    //           items={nav}
    //         />
    //       </EuiFlexItem>
    //       <EuiFlexItem grow={6}>
    //         <EuiPageContent>
    //           <Router primary={false}>
    //             <Redirect from="/" to="user-roles" noThrow />
    //             <UserRoleSetting path="/user-roles" />
    //             <SecretSetting path="/secrets-management" />
    //             <Redirect default from="any" to="/errors/404" noThrow />
    //           </Router>
    //         </EuiPageContent>
    //       </EuiFlexItem>
    //     </EuiFlexGroup>
    //   </EuiPageBody>
    // </EuiPage>
  );
};

export default ProjectSetting;

export const SettingsSection = ({ title, children }) => (
  <Fragment>
    <EuiTitle size="m">
      <h3>{title}</h3>
    </EuiTitle>
    {children}
  </Fragment>
);
