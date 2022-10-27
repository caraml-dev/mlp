import React, { useEffect, useState } from "react";
import {
  EuiFlexItem,
  EuiText,
  EuiFlexGroup,
  EuiIcon,
  EuiButtonIcon,
  EuiInMemoryTable
} from "@elastic/eui";
import UserRoleForm from "./UserRoleForm";
import DeleteUserRoleModal from "./DeleteUserRoleModal";
import SubmitUserRoleForm from "./SubmitUserRoleForm";

export const ADMINISTRATOR_ROLE = "Administrator";
export const READER_ROLE = "Reader";
export const ROLE_OPTIONS = [ADMINISTRATOR_ROLE, READER_ROLE];
const defaultTextSize = "s";

const UserRoleListTable = ({ project, fetchUpdates }) => {
  const [userRoles, setUserRoles] = useState([]);
  useEffect(() => {
    if (project) {
      setUserRoles(formatUserRoles(project));
    }
  }, [project]);

  const [userRoleIdToExpandedRowMap, setUserRoleIdToExpandedRowMap] = useState(
    {}
  );
  const toggleEdit = item => {
    const userRoleIdToExpandedRowMapValues = { ...userRoleIdToExpandedRowMap };
    if (userRoleIdToExpandedRowMapValues[item.id]) {
      delete userRoleIdToExpandedRowMapValues[item.id];
    } else {
      userRoleIdToExpandedRowMapValues[item.id] = (
        <SubmitUserRoleForm
          userRole={item}
          project={project}
          fetchUpdates={fetchUpdates}
        />
      );
    }
    setUserRoleIdToExpandedRowMap(userRoleIdToExpandedRowMapValues);
  };

  const [userRoleToDelete, setUserRoleToDelete] = useState({});
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const toggleDelete = userRole => {
    setUserRoleToDelete(userRole);
    setShowDeleteModal(true);
  };

  const columns = [
    {
      width: "24px",
      render: () => (
        <EuiIcon type="user" size="m" style={{ verticalAlign: "bottom" }} />
      )
    },
    {
      field: "user",
      name: "User",
      sortable: true,
      width: "55%",
      render: name => <EuiText size={defaultTextSize}>{name}</EuiText>
    },
    {
      field: "roles",
      name: "Roles",
      width: "30%",
      render: roles => {
        return (
          <EuiFlexGroup direction="column" gutterSize="none">
            {roles.map(role => (
              <EuiFlexItem key={`${role}`}>
                <EuiText size={defaultTextSize}>{role}</EuiText>
              </EuiFlexItem>
            ))}
          </EuiFlexGroup>
        );
      }
    },
    {
      width: "40px",
      isExpander: true,
      render: item => (
        <EuiButtonIcon
          onClick={() => toggleEdit(item)}
          aria-label={
            userRoleIdToExpandedRowMap[item.id] ? "Collapse" : "Expand"
          }
          iconType={userRoleIdToExpandedRowMap[item.id] ? "arrowUp" : "pencil"}
        />
      )
    },
    {
      width: "40px",
      render: item => (
        <EuiButtonIcon
          onClick={() => toggleDelete(item)}
          aria-label="Delete"
          iconType="trash"
          color="danger"
        />
      )
    }
  ];

  return (
    <>
      {showDeleteModal && (
        <DeleteUserRoleModal
          project={project}
          userRole={userRoleToDelete}
          closeModal={() => setShowDeleteModal(false)}
          fetchUpdates={fetchUpdates}
        />
      )}
      <EuiFlexGroup direction="column">
        <EuiFlexItem>
          {userRoles.length === 0 ? (
            <EuiText size="m" textAlign="center">
              This project has no active User Roles.
            </EuiText>
          ) : (
            <EuiInMemoryTable
              columns={columns}
              itemId="id"
              items={userRoles}
              itemIdToExpandedRowMap={userRoleIdToExpandedRowMap}
              sorting={{ sort: { field: "User", direction: "asc" } }}
            />
          )}
        </EuiFlexItem>
        <EuiFlexItem>
          <UserRoleForm project={project} fetchUpdates={fetchUpdates} />
        </EuiFlexItem>
      </EuiFlexGroup>
    </>
  );
};

export default UserRoleListTable;

const formatUserRoles = project => {
  const administrators = project.administrators?.map(user => ({
    user: user,
    role: ADMINISTRATOR_ROLE
  }));

  const readers = project.readers?.map(user => ({
    user: user,
    role: READER_ROLE
  }));

  const groupByMap = key => array => mapFunc =>
    array.reduce((objectsByKeyValue, obj) => {
      const value = obj[key];
      objectsByKeyValue[value] = (objectsByKeyValue[value] || []).concat(
        mapFunc(obj)
      );
      return objectsByKeyValue;
    }, {});

  const rolesPerUser = groupByMap("user")([...administrators, ...readers])(
    obj => obj.role
  );
  return Object.keys(rolesPerUser).map((user, index) => ({
    id: index,
    user: user,
    roles: rolesPerUser[user]
  }));
};
