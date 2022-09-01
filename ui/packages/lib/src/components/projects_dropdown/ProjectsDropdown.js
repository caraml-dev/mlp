import React, { Fragment, useContext, useEffect, useState } from "react";
import {
  EuiButton,
  EuiIcon,
  EuiPopover,
  EuiPopoverFooter,
  EuiPopoverTitle,
  EuiSelectable
} from "@elastic/eui";
import { useToggle } from "../../hooks";
import { ProjectsContext } from "../../providers";
import "./ProjectDropdown.scss";

const NUM_PANELS_TO_ADD_SEARCH = 7;

export const ProjectsDropdown = ({ onProjectSelect }) => {
  const { projects, currentProject } = useContext(ProjectsContext);

  const [isPopoverOpen, togglePopover] = useToggle();
  const [panels, setPanels] = useState([]);

  useEffect(() => {
    const panels = projects
      .sort((a, b) => (a.name > b.name ? 1 : -1))
      .map(p => {
        return {
          label: p.name,
          checked: p.id === currentProject?.id ? "on" : undefined,
          id: p.id,
          prepend: <EuiIcon type="folderClosed" />
        };
      });
    setPanels(panels);
  }, [projects, currentProject]);

  const onChange = panels => {
    const selectedProject = panels.filter(panel => panel.checked)[0];
    togglePopover();
    onProjectSelect(selectedProject.id);
  };

  return (
    <EuiPopover
      id="projectSelector"
      initialFocus=".euiFieldSearch"
      button={
        <EuiButton
          iconSide="right"
          iconType="arrowDown"
          size="s"
          fill
          onClick={togglePopover}>
          {currentProject?.name || "Projects"}
        </EuiButton>
      }
      isOpen={isPopoverOpen}
      closePopover={togglePopover}
      panelPaddingSize="s"
      anchorPosition="downLeft">
      <EuiSelectable
        searchable={panels.length > NUM_PANELS_TO_ADD_SEARCH}
        searchProps={{
          placeholder: "Find a project",
          compressed: true
        }}
        options={panels}
        singleSelection="always"
        style={{ width: 256 }}
        onChange={onChange}
        listProps={{
          className: "euiSelectableList-projectList",
          rowHeight: 40,
          showIcons: false
        }}>
        {(list, search) => (
          <Fragment>
            <EuiPopoverTitle>{search || "Your projects"}</EuiPopoverTitle>
            {list}
            <EuiPopoverFooter>
              <EuiButton fullWidth size="s" href="/projects/create">
                Create Project
              </EuiButton>
            </EuiPopoverFooter>
          </Fragment>
        )}
      </EuiSelectable>
    </EuiPopover>
  );
};
