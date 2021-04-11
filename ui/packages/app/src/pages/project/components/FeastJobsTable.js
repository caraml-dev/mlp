import React, { useEffect, useState } from "react";
import { navigate } from "@reach/router";
import { EuiInMemoryTable } from "@elastic/eui";

export const FeastJobsTable = ({ project, feastIngestionJobs }) => {
  const [items, setItems] = useState([]);

  useEffect(() => {
    if (project && feastIngestionJobs) {
      let items = [];
      feastIngestionJobs.jobs
        .filter(job => job.status === "JOB_STATUS_RUNNING")
        .sort((a, b) => (a.startTime > b.startTime ? 1 : -1))
        .forEach(job => {
          if (items.length >= 5) {
            return;
          }
          let tableName = "";
          if (job.type === "STREAM_INGESTION_JOB") {
            tableName = job.streamIngestion.tableName;
          } else if (job.type === "BATCH_INGESTION_JOB") {
            tableName = job.batchIngestion.tableName;
          }

          items.push({
            id: job.id,
            table_name: tableName,
            type: job.type === "STREAM_INGESTION_JOB" ? "stream" : "batch",
            status: "running"
          });
        });
      setItems(items);
    }
  }, [project, feastIngestionJobs]);

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

  const cellProps = () => ({
    style: { cursor: "pointer" },
    onClick: () => navigate(`/feast/projects/${project.id}`),
  });

  return <EuiInMemoryTable items={items} columns={columns} cellProps={cellProps} />;
};
