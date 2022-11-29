import React, { useEffect, useState } from "react";
import { EuiListGroup, EuiText } from "@elastic/eui";
import { EXPERIMENT_TYPE_NAME_MAP } from "../../../services/turing/Turing";

import "./ListGroup.scss";

export const TuringRouters = ({ project, routers, homepage }) => {
  const [experiments, setExperiments] = useState([]);

  useEffect(() => {
    if (project && routers) {
      let expMap = new Map();
      routers
        .filter(router => router.config.experiment_engine.type !== "nop")
        .sort((a, b) =>
          a.config.experiment_engine.type > b.config.experiment_engine.type
            ? 1
            : -1
        )
        .forEach(router => {
          const expType = router.config.experiment_engine.type;
          if (!expMap.has(expType)) {
            expMap.set(expType, 1);
          } else {
            expMap.set(expType, expMap.get(expType) + 1);
          }
        });

      let exps = [];
      expMap.forEach((expCount, expType) => {
        exps.push({
          className: "listGroupItem",
          label: (
            <EuiText size="s">
              {expCount} {EXPERIMENT_TYPE_NAME_MAP[expType]}
            </EuiText>
          ),
          onClick: () => {
            window.location.href = `${homepage}/projects/${project.id}/routers?experiment_type=${expType}`;
          },
          size: "s"
        });
      });
      setExperiments(exps);
    }
  }, [project, routers]);

  return experiments.length > 0 ? (
    <EuiListGroup
      color="primary"
      flush={true}
      gutterSize="none"
      listItems={experiments}
    />
  ) : (
    <EuiText size="s">-</EuiText>
  );
};
