import React, { useEffect, useState } from "react";
import { EuiListGroup, EuiText } from "@elastic/eui";
import { useFeastCoreApi } from "../../../hooks/useFeastCoreApi";

import "./ListGroup.scss";

export const FeastResources = ({ project, homepage }) => {
  const [{ data: entities }] = useFeastCoreApi(
    `/entities?project=${project?.name}`,
    { method: "GET" },
    { entities: [] },
    !!project
  );
  const [{ data: featureTables }] = useFeastCoreApi(
    `/tables?project=${project?.name}`,
    { method: "GET" },
    { tables: [] },
    !!project
  );

  const [items, setItems] = useState([]);

  useEffect(() => {
    setItems([
      {
        className: "listGroupItem",
        label: <EuiText size="s">{entities.entities.length} entities</EuiText>,
        onClick: () => {
          window.location.href = `${homepage}/projects/${project.id}/entities`;
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
          window.location.href = `${homepage}/projects/${project.id}/featuretables`;
        },
        size: "s"
      }
    ]);
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
