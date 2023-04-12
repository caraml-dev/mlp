import React, { useContext } from "react";
import UserRoleListTable from "./user_role/UserRoleListTable";
import { ProjectsContext } from "@caraml-dev/ui-lib";
import { EuiLoadingChart, EuiTextAlign } from "@elastic/eui";

const UserRoleSetting = () => {
  const { currentProject, refresh } = useContext(ProjectsContext);

  return (
    <>
      {!currentProject ? (
        <EuiTextAlign textAlign="center">
          <EuiLoadingChart size="xl" mono />
        </EuiTextAlign>
      ) : (
        <UserRoleListTable project={currentProject} fetchUpdates={refresh} />
      )}
    </>
  );
};

export default UserRoleSetting;
