import React, { useEffect, useState } from "react";
import { navigate } from "@reach/router";
import { EuiLink, EuiListGroup, EuiText } from "@elastic/eui";

export const MerlinDeployments = ({ project, models }) => {
  const [deployedModels, setDeployedModels] = useState([]);

  useEffect(() => {
    if (project && models) {
      const items = models
        .filter(
          model =>
            model.endpoints &&
            model.endpoints.length > 0 &&
            model.endpoints.find(endpoint => endpoint.status === "serving")
        )
        .sort((a, b) => (a.name > b.name ? 1 : -1))
        .map(model => {
          let totalStandardTransformer = 0;
          let totalCustomTransformer = 0;
          let totalServing = 0;
          model.endpoints.forEach(endpoint => {
            totalServing++;

            endpoint.rule.destinations.forEach(destination => {
              if (destination.version_endpoint.transformer) {
                if (
                  destination.version_endpoint.transformer.transformer_type ===
                  "standard"
                ) {
                  totalStandardTransformer++;
                } else {
                  totalCustomTransformer++;
                }
              }
            });
          });
          return {
            label: (
              <>
                <EuiLink
                  href={`/merlin/projects/${project.id}/models/${model.id}`}
                  size="s">
                  {model.name}
                </EuiLink>
                {totalStandardTransformer > 0 && (
                  <EuiText size="s">
                    {totalStandardTransformer} active Standard Transformer
                  </EuiText>
                )}
                {totalCustomTransformer > 0 && (
                  <EuiText size="s">
                    {totalCustomTransformer} active Custom Transformer
                  </EuiText>
                )}
                <EuiText size="s">{totalServing} active serving</EuiText>
              </>
            ),
            onClick: () => {
              navigate(`/merlin/projects/${project.id}/models/${model.id}`);
            },
            size: "s"
          };
        });
      setDeployedModels(items);
    }
  }, [project, models]);

  return (
    <EuiListGroup
      color="text"
      flush={true}
      gutterSize="none"
      listItems={deployedModels}
    />
  );
};
