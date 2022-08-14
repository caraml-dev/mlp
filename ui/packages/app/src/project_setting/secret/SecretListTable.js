import React, { Fragment, useState } from "react";
import {
  EuiFlexItem,
  EuiFlexGroup,
  EuiInMemoryTable,
  EuiIcon,
  EuiButtonIcon,
  EuiToolTip,
  EuiPanel,
  EuiText
} from "@elastic/eui";
import SecretForm from "./SecretForm";
import DeleteSecretModal from "./DeleteSecretModal";
import SubmitSecretForm from "./SubmitSecretForm";

const moment = require("moment");
const defaultTextSize = "s";

const SecretListTable = ({ secrets, projectId, fetchUpdates }) => {
  const [secretIdToExpandedRowMap, setSecretIdToExpandedRowMap] = useState({});
  const toggleEdit = item => {
    const secretIdToExpandedRowMapValues = { ...secretIdToExpandedRowMap };
    if (secretIdToExpandedRowMapValues[item.id]) {
      delete secretIdToExpandedRowMapValues[item.id];
    } else {
      secretIdToExpandedRowMapValues[item.id] = (
        <SubmitSecretForm
          secret={item}
          projectId={projectId}
          fetchUpdates={fetchUpdates}
        />
      );
    }
    setSecretIdToExpandedRowMap(secretIdToExpandedRowMapValues);
  };

  const [secretToDelete, setSecretToDelete] = useState({});
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const toggleDelete = secret => {
    setSecretToDelete(secret);
    setShowDeleteModal(true);
  };

  const columns = [
    {
      width: "24px",
      render: () => (
        <EuiIcon type="lock" size="m" style={{ verticalAlign: "bottom" }} />
      )
    },
    {
      field: "name",
      name: "Name",
      sortable: true,
      width: "45%",
      render: name => <EuiText size={defaultTextSize}>{name}</EuiText>
    },
    {
      field: "created_at",
      name: "Created",
      width: "20%",
      render: date => (
        <EuiToolTip
          position="top"
          content={moment(date, "YYYY-MM-DDTHH:mm.SSZ").toLocaleString()}>
          <EuiText size={defaultTextSize}>
            {moment(date, "YYYY-MM-DDTHH:mm.SSZ").fromNow()}
          </EuiText>
        </EuiToolTip>
      )
    },
    {
      field: "updated_at",
      name: "Updated",
      width: "20%",
      render: date => (
        <EuiToolTip
          position="top"
          content={moment(date, "YYYY-MM-DDTHH:mm.SSZ").toLocaleString()}>
          <EuiText size={defaultTextSize}>
            {moment(date, "YYYY-MM-DDTHH:mm.SSZ").fromNow()}
          </EuiText>
        </EuiToolTip>
      )
    },
    {
      width: "40px",
      isExpander: true,
      render: item => (
        <EuiButtonIcon
          onClick={() => toggleEdit(item)}
          aria-label={secretIdToExpandedRowMap[item.id] ? "Collapse" : "Expand"}
          iconType={secretIdToExpandedRowMap[item.id] ? "arrowUp" : "pencil"}
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

  if (secrets.length === 0) {
    return (
      <Fragment>
        <EuiFlexGroup direction="column">
          <EuiFlexItem>
            <EuiText size="s">
              <EuiFlexItem alignitems="center">
                <EuiPanel>
                  <h4 align="center">This project has no active Secrets.</h4>
                </EuiPanel>
              </EuiFlexItem>
            </EuiText>
          </EuiFlexItem>
          <EuiFlexItem>
            <SecretForm projectId={projectId} fetchUpdates={fetchUpdates} />
          </EuiFlexItem>
        </EuiFlexGroup>
      </Fragment>
    );
  } else {
    return (
      <Fragment>
        {showDeleteModal && (
          <DeleteSecretModal
            projectId={projectId}
            secret={secretToDelete}
            closeModal={() => setShowDeleteModal(false)}
            fetchUpdates={fetchUpdates}
          />
        )}
        <EuiFlexGroup direction="column">
          <EuiFlexItem>
            <EuiInMemoryTable
              columns={columns}
              itemId="id"
              items={secrets}
              itemIdToExpandedRowMap={secretIdToExpandedRowMap}
              sorting={{ sort: { field: "Name", direction: "asc" } }}
            />
          </EuiFlexItem>
          <EuiFlexItem>
            <SecretForm projectId={projectId} fetchUpdates={fetchUpdates} />
          </EuiFlexItem>
        </EuiFlexGroup>
      </Fragment>
    );
  }
};
export default SecretListTable;
