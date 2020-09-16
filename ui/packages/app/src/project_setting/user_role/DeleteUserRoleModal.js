import React, { useEffect } from "react";
import { useMlpApi } from "@mlp/ui";
import { EuiConfirmModal, EuiOverlayMask } from "@elastic/eui";

const DeleteUserRoleModal = ({
  project,
  userRole,
  closeModal,
  fetchUpdates
}) => {
  const [deleteResponse, deleteUserRole] = useMlpApi(
    `/projects/${project.id}`,
    {
      method: "PUT"
    },
    {},
    false
  );

  const removeUserRole = () => {
    const proj = { ...project };
    userRole.roles.forEach(role => {
      role === "Administrator"
        ? (proj.administrators = proj.administrators.filter(
            a => a !== userRole.user
          ))
        : (proj.readers = proj.readers.filter(r => r !== userRole.user));
    });
    return proj;
  };

  useEffect(() => {
    if (deleteResponse.isLoaded && !deleteResponse.error) {
      closeModal();
      fetchUpdates();
    }
  }, [deleteResponse, project, fetchUpdates, closeModal]);

  return (
    <EuiOverlayMask>
      <EuiConfirmModal
        title="Delete user role"
        onCancel={() => closeModal()}
        onConfirm={() =>
          deleteUserRole({ body: JSON.stringify(removeUserRole()) })
        }
        cancelButtonText="Cancel"
        confirmButtonText="Delete"
        buttonColor="danger"
        defaultFocusedButton="confirm">
        <p>
          You are about to remove user <b>{userRole.user}</b>, are you sure?
        </p>
      </EuiConfirmModal>
    </EuiOverlayMask>
  );
};

export default DeleteUserRoleModal;
