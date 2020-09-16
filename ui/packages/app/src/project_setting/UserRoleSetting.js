import React, { useContext, useEffect } from "react";
import UserRoleListTable from "./user_role/UserRoleListTable";
import { CurrentProjectContext, replaceBreadcrumbs } from "@mlp/ui";
import { SettingsSection } from "./ProjectSetting";
import { EuiLoadingChart, EuiTextAlign } from "@elastic/eui";

const UserRoleSetting = () => {
  const { project, refresh } = useContext(CurrentProjectContext);

  useEffect(() => {
    if (project) {
      replaceBreadcrumbs([
        {
          text: "User Roles"
        }
      ]);
    }
  }, [project]);

  return (
    <SettingsSection title="User Roles">
      {!project ? (
        <EuiTextAlign textAlign="center">
          <EuiLoadingChart size="xl" mono />
        </EuiTextAlign>
      ) : (
        <UserRoleListTable project={project} fetchUpdates={() => refresh()} />
      )}
    </SettingsSection>
  );
};

export default UserRoleSetting;
