import React, { useEffect, useState } from "react";
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
  EuiPanel
} from "@elastic/eui";
import { addToast, useMlpApi } from "@mlp/ui";
import {
  validateSecretKey,
  validateSecretName
} from "../../validation/validateSecret";

const SubmitSecretForm = ({ projectId, fetchUpdates, secret, toggleAdd }) => {
  const [request, setRequest] = useState({
    name: secret ? secret.name : "",
    data: ""
  });

  const [submissionResponse, submitForm] = useMlpApi(
    secret
      ? `/projects/${projectId}/secrets/${secret.id}`
      : `/projects/${projectId}/secrets`,
    {
      method: secret ? "PATCH" : "POST",
      headers: { "Content-Type": "application/json" }
    },
    {},
    false
  );

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
    setValidKey(validateSecretKey(e.target.value));
    onChange("data")(e.target.value);
  };

  const saveAction = () => {
    submitForm({ body: JSON.stringify(request) });
  };

  const [isValidName, setValidName] = useState(false);
  const [isValidKey, setValidKey] = useState(false);

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
                <EuiToolTip content="Specify content of secret">
                  <span>
                    Key <EuiIcon type="questionInCircle" color="subdued" />
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
                isInvalid={!isValidKey}
              />
            </EuiFormRow>
            <EuiFormRow>
              <EuiFlexGroup direction="row">
                <EuiFlexItem grow={false}>
                  <EuiButton
                    fill
                    size="s"
                    disabled={
                      secret ? !isValidKey : !isValidKey || !isValidName
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
