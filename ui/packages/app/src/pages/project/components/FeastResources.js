import React, { useEffect, useState } from "react";
import { EuiListGroup, EuiText } from "@elastic/eui";

import "./ListGroup.scss";

export const FeastResources = ({ project, entities, featureTables }) => {
  const [items, setItems] = useState([]);

  useEffect(() => {
    if (
      project &&
      entities &&
      entities.entities &&
      featureTables &&
      featureTables.tables
    ) {
      setItems([
        {
          className: "listGroupItem",
          label: (
            <EuiText size="s">{entities.entities.length} entities</EuiText>
          ),
          onClick: () => {
            window.location.href = `/feast/projects/${project.id}/entities`;
          },
          size: "s"
        },
        {
          className: "listGroupItem",
          label: (
            <EuiText size="s">
              {featureTables.tables.length} feature tables
            </EuiText>
          ),
          onClick: () => {
            window.location.href = `/feast/projects/${project.id}/featuretables`;
          },
          size: "s"
        }
      ]);
    }
  }, [project, entities, featureTables]);

  return items.length > 0 ? (
    <EuiListGroup
      color="primary"
      flush={true}
      gutterSize="none"
      listItems={items}
    />
  ) : (
    <EuiText size="s">-</EuiText>
  );
};
