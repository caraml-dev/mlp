import React, { useState } from "react";
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

export const Resources = ({
  project,
  entities,
  featureTables,
  models,
  routers
}) => {
  const [isPopoverOpen, setPopover] = useState(false);

  const onButtonClick = () => {
    setPopover(!isPopoverOpen);
  };

  const closePopover = () => {
    setPopover(false);
  };

  const newResourceMenus = [
    <EuiContextMenuItem
      href={config.KUBEFLOW_UI_HOMEPAGE}
      key="notebook"
      size="s">
      Create notebook
    </EuiContextMenuItem>,
    <EuiContextMenuItem
      href={`${config.FEAST_UI_HOMEPAGE}/projects/${project.id}/entities/create`}
      key="feature-table"
      size="s">
      Create Entities
    </EuiContextMenuItem>,
    <EuiContextMenuItem
      href={`${config.FEAST_UI_HOMEPAGE}/projects/${project.id}/featuretables/create`}
      key="feature-table"
      size="s">
      Create FeatureTable
    </EuiContextMenuItem>,
    <EuiContextMenuItem href={`${config.MERLIN_UI_HOMEPAGE}/projects/${project.id}`} key="model" size="s">
      Deploy model
    </EuiContextMenuItem>,
    <EuiContextMenuItem
      href={`${config.TURING_UI_HOMEPAGE}/projects/${project.id}/routers/create`}
      key="experiment"
      size="s">
      Set up experiment
    </EuiContextMenuItem>,
    <EuiContextMenuItem href={config.CLOCKWORK_UI_HOMEPAGE} key="job" size="s">
      Schedule a job
    </EuiContextMenuItem>
  ];

  const items = [
    {
      title: "Features",
      description: (
        <FeastResources
          project={project}
          entities={entities}
          featureTables={featureTables}
        />
      )
    },
    {
      title: "Models",
      description: <MerlinModels project={project} models={models} />
    },
    {
      title: "Experiments",
      description: <TuringRouters project={project} routers={routers} />
    },
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
    },
  ];

  const actions = (
    <EuiPopover
      button={
        <EuiLink onClick={onButtonClick}>
          <EuiText size="s"><EuiIcon type="arrowRight" /> Create a new resource</EuiText>
        </EuiLink>
      }
      isOpen={isPopoverOpen}
      closePopover={closePopover}
      anchorPosition="rightUp"
      offset={10}
      panelPaddingSize="s">
      <EuiContextMenuPanel size="s" items={newResourceMenus} />
    </EuiPopover>
  );

  return <Panel title="Resources" items={items} actions={actions} iconType="beaker" />;
};
