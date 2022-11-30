import React, { useEffect, useState } from "react";
import { EuiListGroup, EuiText } from "@elastic/eui";
import { MODEL_TYPE_NAME_MAP } from "../../../services/merlin/Model";

import "./ListGroup.scss";

export const MerlinModels = ({ project, models, homepage }) => {
  const [modelItems, setModelItems] = useState([]);

  useEffect(() => {
    if (project && models) {
      let modelMap = new Map();
      models
        .sort((a, b) => (a.type > b.type ? 1 : -1))
        .forEach(model => {
          const modelType = model.type === "pyfunc_v2" ? "pyfunc" : model.type;
          if (!modelMap.has(modelType)) {
            modelMap.set(modelType, 1);
          } else {
            modelMap.set(modelType, modelMap.get(modelType) + 1);
          }
        });

      let items = [];
      modelMap.forEach((modelCount, modelType) => {
        items.push({
          className: "listGroupItem",
          label: (
            <EuiText size="s">
              {modelCount} {MODEL_TYPE_NAME_MAP[modelType]}
            </EuiText>
          ),
          onClick: () => {
            window.location.href = `${homepage}/projects/${project.id}/models?type=${modelType}`;
          },
          size: "s"
        });
      });
      setModelItems(items);
    }
  }, [project, models, homepage]);

  return modelItems.length > 0 ? (
    <EuiListGroup
      color="primary"
      flush={true}
      gutterSize="none"
      listItems={modelItems}
    />
  ) : (
    <EuiText size="s">-</EuiText>
  );
};
