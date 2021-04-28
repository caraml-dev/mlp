import React, { useEffect, useState } from "react";
import { EuiInMemoryTable } from "@elastic/eui";

export const MerlinDeploymentsTable = ({ project, models }) => {
  const [items, setItems] = useState([]);

  useEffect(() => {
    if (project && models) {
      let items = [];
      models
        .filter(
          model =>
            model.endpoints &&
            model.endpoints.length > 0 &&
            model.endpoints.find(endpoint => endpoint.status === "serving")
        )
        .sort((a, b) => (a.name > b.name ? 1 : -1))
        .forEach(model => {
          model.endpoints.forEach(endpoint => {
            if (items.length >= 5) {
              return;
            }

            endpoint.rule.destinations.forEach(destination => {
              items.push({
                id: model.id,
                name: model.name,
                environment_name: endpoint.environment_name,
                endpoint: endpoint.url,
                version_id: destination.version_endpoint.version_id,
                version_endpoint_id: destination.version_endpoint.id,
              });
            });
          });
        });
      setItems(items);
    }
  }, [project, models]);

  const columns = [
    {
      field: "name",
      name: "Name",
      sortable: true,
      width: "25%"
    },
    {
      field: "environment_name",
      name: "Environment",
      width: "15%"
    },
    {
      field: "endpoint",
      name: "Endpoint",
      width: "60%"
    }
  ];

  const cellProps = item => ({
    style: { cursor: "pointer" },
    onClick: () => window.location.href = `/merlin/projects/${project.id}/models/${item.id}/versions/${item.version_id}/endpoints/${item.version_endpoint_id}/details`,
  });

  return <EuiInMemoryTable items={items} columns={columns} cellProps={cellProps} />;
};
