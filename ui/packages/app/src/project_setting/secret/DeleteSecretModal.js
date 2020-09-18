import React, { useEffect } from "react";
import { useMlpApi } from "@gojek/mlp-ui";
import { EuiConfirmModal, EuiOverlayMask } from "@elastic/eui";

const DeleteSecretModal = ({ projectId, secret, closeModal, fetchUpdates }) => {
  const [deleteResponse, deleteSecret] = useMlpApi(
    `/projects/${projectId}/secrets/${secret.id}`,
    {
      method: "DELETE"
    },
    {},
    false
  );

  useEffect(() => {
    if (deleteResponse.isLoaded && !deleteResponse.error) {
      closeModal();
      fetchUpdates();
    }
  }, [deleteResponse, projectId, fetchUpdates, closeModal]);

  return (
    <EuiOverlayMask>
      <EuiConfirmModal
        title="Delete secret"
        onCancel={() => closeModal()}
        onConfirm={() => deleteSecret()}
        cancelButtonText="Cancel"
        confirmButtonText="Delete"
        buttonColor="danger"
        defaultFocusedButton="confirm">
        <p>
          You are about to remove secret <b>{secret.name}</b>, are you sure?
        </p>
      </EuiConfirmModal>
    </EuiOverlayMask>
  );
};

export default DeleteSecretModal;
