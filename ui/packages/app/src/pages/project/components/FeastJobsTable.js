import React, { useEffect, useState } from "react";
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
          let jobType = "";
          if (job.type === "STREAM_INGESTION_JOB") {
            tableName = job.streamIngestion.tableName;
            jobType = "stream";
          } else if (job.type === "BATCH_INGESTION_JOB") {
            tableName = job.batchIngestion.tableName;
            jobType = "batch";
          }

          items.push({
            id: job.id,
            table_name: tableName,
            type: jobType,
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

  const cellProps = item => ({
    style: { cursor: "pointer" },
    onClick: () => window.location.href = `/feast/projects/${project.id}/jobs/${item.type}`,
  });

  return <EuiInMemoryTable items={items} columns={columns} cellProps={cellProps} />;
};
