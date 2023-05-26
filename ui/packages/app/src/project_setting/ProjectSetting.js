import React, { useState } from "react";
import { EuiSideNav, EuiIcon, EuiPageTemplate } from "@elastic/eui";
import { slugify } from "@caraml-dev/ui-lib/src/utils";
import UserRoleSetting from "./UserRoleSetting";
import SecretSetting from "./SecretSetting";
import {
  Navigate,
  Route,
  Routes,
  useNavigate,
  useParams
} from "react-router-dom";

const sections = {
  "user-roles": {
    iconType: "user",
    name: "User Roles"
  },
  "secrets-management": {
    iconType: "lock",
    name: "Secrets"
  }
};

const ProjectSetting = () => {
  const { "*": section } = useParams();

  const navigate = useNavigate();
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
        onClick: () => navigate(`${id}`),
        isSelected: slugify(id) === section
      }))
    }
  ];

  return (
    <EuiPageTemplate>
      <EuiPageTemplate.Header
        restrictWidth={false}
        iconType={sections[section]?.iconType}
        pageTitle={sections[section]?.name}
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
        <Routes>
          <Route index element={<Navigate to="user-roles" replace={true} />} />
          <Route path="user-roles" element={<UserRoleSetting />} />
          <Route path="secrets-management" element={<SecretSetting />} />
          <Route
            path="*"
            element={<Navigate to="/pages/404" replace={true} />}
          />
        </Routes>
      </EuiPageTemplate.Section>
    </EuiPageTemplate>
  );
};

export default ProjectSetting;
