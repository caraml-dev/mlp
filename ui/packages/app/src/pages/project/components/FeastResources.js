import React, { useEffect, useState } from "react";
import { navigate } from "@reach/router";
import { EuiListGroup, EuiText } from "@elastic/eui";

import "./ListGroup.scss";

export const FeastResources = ({ project, entities, featureTables }) => {
  const [items, setItems] = useState([]);

  useEffect(() => {
    if (project && entities && featureTables) {
      setItems([
        {
          className: "listGroupItem",
          label: (
            <EuiText size="s">{entities.entities.length} entities</EuiText>
          ),
          onClick: () => {
            navigate(`/feast/entities`);
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
            navigate(`/feast/featuretables`);
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
