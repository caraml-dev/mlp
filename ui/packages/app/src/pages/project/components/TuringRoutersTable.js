import React, { useEffect, useState } from "react";
import { navigate } from "@reach/router";
import { EuiInMemoryTable } from "@elastic/eui";

export const TuringRoutersTable = ({ project, routers }) => {
  const [items, setItems] = useState([]);

  useEffect(() => {
    if (project && routers) {
      let items = [];
      routers
        .filter(router => router.status === "deployed")
        .sort((a, b) => (a.name > b.name ? 1 : -1))
        .forEach(router => {
          if (items.length >= 5) {
            return;
          }
          const { hostname } = new URL(router.endpoint);
          items.push({
            id: router.id,
            name: router.name,
            environment_name: router.environment_name,
            endpoint: hostname
          });
        });
      setItems(items);
    }
  }, [project, routers]);

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
    onClick: () => navigate(`/turing/projects/${project.id}/routers/${item.id}/details`),
  });

  return <EuiInMemoryTable items={items} columns={columns} cellProps={cellProps} />;
};
