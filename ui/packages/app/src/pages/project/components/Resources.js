import React from "react";
import {
  EuiContextMenuItem,
  EuiContextMenuPanel,
  EuiIcon,
  EuiLink,
  EuiPopover,
  EuiText
} from "@elastic/eui";
import config from "../../../config";
import { Panel } from "./Panel";
import { FeastResources } from "./FeastResources";
import { MerlinModels } from "./MerlinModels";
import { TuringRouters } from "./TuringRouters";
import { useToggle } from "@caraml-dev/ui-lib";

export const Resources = ({ apps, project, models, routers }) => {
  const [isPopoverOpen, togglePopover] = useToggle(false);

  const newResourceMenus = [
    <EuiContextMenuItem
      href={config.KUBEFLOW_UI_HOMEPAGE}
      key="notebook"
      size="s">
      Create notebook
    </EuiContextMenuItem>,
    ...(apps.some(a => a.name === "Feast")
      ? [
          <EuiContextMenuItem
            href={`${apps.find(a => a.name === "Feast")?.homepage}/projects/${
              project.id
            }/entities/create`}
            key="feature-table"
            size="s">
            Create Entities
          </EuiContextMenuItem>,
          <EuiContextMenuItem
            href={`${apps.find(a => a.name === "Feast")?.homepage}/projects/${
              project.id
            }/featuretables/create`}
            key="feature-table"
            size="s">
            Create FeatureTable
          </EuiContextMenuItem>
        ]
      : []),
    ...(apps.some(a => a.name === "Merlin")
      ? [
          <EuiContextMenuItem
            href={`${apps.find(a => a.name === "Merlin")?.homepage}/projects/${
              project.id
            }`}
            key="model"
            size="s">
            Deploy model
          </EuiContextMenuItem>
        ]
      : []),
    ...(apps.some(a => a.name === "Turing")
      ? [
          <EuiContextMenuItem
            href={`${apps.find(a => a.name === "Turing")?.homepage}/projects/${
              project.id
            }/routers/create`}
            key="experiment"
            size="s">
            Set up experiment
          </EuiContextMenuItem>
        ]
      : []),
    <EuiContextMenuItem href={config.CLOCKWORK_UI_HOMEPAGE} key="job" size="s">
      Schedule a job
    </EuiContextMenuItem>
  ];

  const items = [
    ...(apps.some(a => a.name === "Feast")
      ? [
          {
            title: "Features",
            description: (
              <FeastResources
                project={project}
                homepage={apps.find(a => a.name === "Feast")?.homepage}
              />
            )
          }
        ]
      : []),
    ...(apps.some(a => a.name === "Merlin")
      ? [
          {
            title: "Models",
            description: (
              <MerlinModels
                project={project}
                models={models}
                homepage={apps.find(a => a.name === "Merlin")?.homepage}
              />
            )
          }
        ]
      : []),
    ...(apps.some(a => a.name === "Turing")
      ? [
          {
            title: "Experiments",
            description: (
              <TuringRouters
                project={project}
                routers={routers}
                homepage={apps.find(a => a.name === "Turing")?.homepage}
              />
            )
          }
        ]
      : []),
    {
      title: "Notebooks",
      description: (
        <EuiLink href={config.KUBEFLOW_UI_HOMEPAGE}>Kubeflow Notebooks</EuiLink>
      )
    },
    {
      title: "Pipelines",
      description: (
        <EuiLink href={config.CLOCKWORK_UI_HOMEPAGE}>
          Clockwork Pipelines
        </EuiLink>
      )
    }
  ];

  const actions = (
    <EuiPopover
      button={
        <EuiLink onClick={togglePopover}>
          <EuiText size="s">
            <EuiIcon type="arrowRight" /> Create a new resource
          </EuiText>
        </EuiLink>
      }
      isOpen={isPopoverOpen}
      closePopover={togglePopover}
      anchorPosition="rightUp"
      offset={10}
      panelPaddingSize="s">
      <EuiContextMenuPanel size="s" items={newResourceMenus} />
    </EuiPopover>
  );

  return (
    <Panel
      title="Resources"
      items={items}
      actions={actions}
      iconType="beaker"
    />
  );
};
