import React, { useEffect, useState, useMemo } from "react";
import {
  EuiFieldText,
  EuiFormRow,
  EuiButton,
  EuiButtonEmpty,
  EuiTextArea,
  EuiFlexItem,
  EuiFlexGroup,
  EuiForm,
  EuiText,
  EuiTitle,
  EuiIcon,
  EuiToolTip,
  EuiSpacer,
  EuiPanel,
  EuiSelect
} from "@elastic/eui";
import { addToast, useMlpApi } from "@caraml-dev/ui-lib";
import {
  validateSecretData,
  validateSecretName
} from "../../validation/validation";

const SubmitSecretForm = ({ projectId, fetchUpdates, secret, toggleAdd }) => {
  const DEFAULT_SECRET_STORAGE_ID = 2;
  const [request, setRequest] = useState({
    name: secret ? secret.name : "",
    // only authenticated user with proper access to the project can get the secret value,
    // so it's safe to show the value here
    data: secret ? secret.data : "",
    secret_storage_id: secret
      ? secret.secret_storage_id
      : DEFAULT_SECRET_STORAGE_ID
  });

  const [listSecretStorageResponse] = useMlpApi(
    `/v1/projects/${projectId}/secret_storages`,
    {
      method: "GET"
    }
  );

  const [submissionResponse, submitForm] = useMlpApi(
    secret
      ? `/v1/projects/${projectId}/secrets/${secret.id}`
      : `/v1/projects/${projectId}/secrets`,
    {
      method: secret ? "PATCH" : "POST",
      headers: { "Content-Type": "application/json" }
    },
    {},
    false
  );

  const secretStorageOptions = useMemo(() => {
    if (
      listSecretStorageResponse.isLoaded &&
      !listSecretStorageResponse.error
    ) {
      return listSecretStorageResponse.data.map(item => {
        return {
          value: item.id,
          text: item.name
        };
      });
    }
    return [];
  }, [listSecretStorageResponse]);

  useEffect(() => {
    if (submissionResponse.isLoaded && !submissionResponse.error) {
      addToast({
        id: "submit-success-create",
        title: secret ? "Secret key changed!" : "New secret is created!",
        color: "success",
        iconType: "check"
      });
      fetchUpdates();

      if (!secret) {
        toggleAdd();
      } else {
        onChange("data")("");
      }
    }
  }, [submissionResponse, fetchUpdates, toggleAdd, secret]);

  const onChange = field => {
    return value => {
      setRequest(r => ({ ...r, [field]: value }));
    };
  };

  const onNameChanges = e => {
    setValidName(validateSecretName(e.target.value));
    onChange("name")(e.target.value);
  };

  const onDataChanges = e => {
    setValidData(validateSecretData(e.target.value));
    onChange("data")(e.target.value);
  };

  const onSecretStorageChanges = e => {
    onChange("secret_storage_id")(parseInt(e.target.value));
  };

  const saveAction = () => {
    submitForm({ body: JSON.stringify(request) });
  };

  const [isValidName, setValidName] = useState(
    validateSecretName(request.name)
  );
  const [isValidData, setValidData] = useState(
    validateSecretData(request.data)
  );

  return (
    <EuiPanel paddingSize="m">
      <EuiFlexGroup direction="column">
        <EuiFlexItem>
          <EuiTitle size="xs">
            <h5> {secret ? "Edit Secret" : "Add a Secret"} </h5>
          </EuiTitle>
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiForm
            isInvalid={submissionResponse.error}
            error={
              submissionResponse.error ? [submissionResponse.error.message] : ""
            }>
            {!secret && (
              <EuiFormRow
                fullWidth
                label={
                  <EuiToolTip content="Specify name of secret">
                    <span>
                      Secret Name{" "}
                      <EuiIcon type="questionInCircle" color="subdued" />
                    </span>
                  </EuiToolTip>
                }
                display="columnCompressed">
                <EuiFieldText
                  placeholder="e.g MERLIN SECRET"
                  value={request.name}
                  onChange={e => onNameChanges(e)}
                  name="name"
                  isInvalid={!isValidName}
                />
              </EuiFormRow>
            )}
            <EuiSpacer size="s" />
            <EuiFormRow
              fullWidth
              label={
                <EuiToolTip content="Specify the secret storage to store the secret">
                  <span>
                    Secret Storage{" "}
                    <EuiIcon type="questionInCircle" color="subdued" />
                  </span>
                </EuiToolTip>
              }
              display="columnCompressed">
              <EuiSelect
                id="select-secret-storage"
                options={secretStorageOptions}
                value={request.secret_storage_id}
                onChange={e => onSecretStorageChanges(e)}
              />
            </EuiFormRow>
            <EuiFormRow
              fullWidth
              label={
                <EuiToolTip content="Specify content of secret">
                  <span>
                    Data <EuiIcon type="questionInCircle" color="subdued" />
                  </span>
                </EuiToolTip>
              }
              display="columnCompressed">
              <EuiTextArea
                fullWidth
                placeholder="e.g p@ass-w0rD"
                value={request.data}
                onChange={e => onDataChanges(e)}
                name="data"
                isInvalid={!isValidData}
              />
            </EuiFormRow>
            <EuiFormRow>
              <EuiFlexGroup direction="row">
                <EuiFlexItem grow={false}>
                  <EuiButton
                    fill
                    size="s"
                    disabled={
                      secret ? !isValidData : !isValidData || !isValidName
                    }
                    onClick={() => saveAction()}>
                    <EuiText size="s"> {secret ? "Save" : "Add"}</EuiText>
                  </EuiButton>
                </EuiFlexItem>
                {!secret && (
                  <EuiFlexItem grow={false}>
                    <EuiButtonEmpty size="s" onClick={() => toggleAdd()}>
                      <EuiText size="s">Cancel</EuiText>
                    </EuiButtonEmpty>
                  </EuiFlexItem>
                )}
              </EuiFlexGroup>
            </EuiFormRow>
          </EuiForm>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPanel>
  );
};

export default SubmitSecretForm;
