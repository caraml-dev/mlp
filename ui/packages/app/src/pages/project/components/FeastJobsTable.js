import React, { useEffect, useState } from "react";
import { EuiInMemoryTable } from "@elastic/eui";

export const FeastJobsTable = ({
  project,
  feastStreamIngestionJobs,
  feastBatchIngestionJobs
}) => {
  const [items, setItems] = useState([]);

  useEffect(() => {
    let items = [];

    if (feastStreamIngestionJobs) {
      feastStreamIngestionJobs
        .filter(job => job.status === "IN PROGRESS")
        .sort((a, b) => (a.startTime > b.startTime ? 1 : -1))
        .forEach(job => {
          if (items.length >= 5) {
            return;
          }

          items.push({
            id: job.job_id,
            table_name: job.feature_table,
            type: "stream",
            status: "running"
          });
        });
    }

    if (feastBatchIngestionJobs) {
      console.log("feastBatchIngestionJobs", feastBatchIngestionJobs);
      feastBatchIngestionJobs
        .filter(job => job.status === "IN PROGRESS")
        .sort((a, b) => (a.startTime > b.startTime ? 1 : -1))
        .forEach(job => {
          if (items.length >= 5) {
            return;
          }

          items.push({
            id: job.job_id,
            table_name: job.feature_table,
            type: "batch",
            status: "running"
          });
        });
    }
    setItems(items);
  }, [feastStreamIngestionJobs, feastBatchIngestionJobs]);

  const columns = [
    {
      field: "id",
      name: "Job ID",
      sortable: true,
      width: "25%"
    },
    {
      field: "table_name",
      name: "Feature Table",
      width: "25%"
    },
    {
      field: "type",
      name: "Ingestion Type",
      width: "25%"
    },
    {
      field: "status",
      name: "Status",
      width: "25%"
    }
  ];

  const cellProps = item => ({
    style: { cursor: "pointer" },
    onClick: () =>
      (window.location.href = `/feast/projects/${project.id}/jobs/${item.type}`)
  });

  return (
    <EuiInMemoryTable items={items} columns={columns} cellProps={cellProps} />
  );
};
