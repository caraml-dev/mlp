import React, {
  Fragment,
  useContext,
  useEffect,
  useMemo,
  useState
} from "react";
import {
  EuiButton,
  EuiIcon,
  EuiPopover,
  EuiPopoverFooter,
  EuiPopoverTitle,
  EuiSelectable
} from "@elastic/eui";
import { CurrentProjectContext } from "../../providers/project";
import { useToggle } from "../../hooks";
import "./ProjectDropdown.scss";

const NUM_PANELS_TO_ADD_SEARCH = 7;

export const ProjectsDropdown = ({ projects, onProjectSelect }) => {
  const { project = {} } = useContext(CurrentProjectContext);

  const popoverLabel = useMemo(() => {
    return project.name || "Projects";
  }, [project.name]);

  const [isPopoverOpen, togglePopover] = useToggle();

  const [panels, setPanels] = useState([]);

  useEffect(() => {
    const panels = projects
      .sort((a, b) => (a.name > b.name ? 1 : -1))
      .map(p => {
        return {
          label: p.name,
          checked: p.name === project.name ? "on" : undefined,
          id: p.id,
          prepend: <EuiIcon type="folderClosed" />
        };
      });
    setPanels(panels);
  }, [projects, project.name]);

  const onChange = panels => {
    const selectedProject = panels.filter(panel => panel.checked)[0];
    togglePopover();
    onProjectSelect(selectedProject.id);
  };

  return (
    <EuiPopover
      id="projectSelector"
      button={
        <EuiButton
          iconSide="right"
          iconType="arrowDown"
          size="s"
          fill
          onClick={togglePopover}>
          {popoverLabel}
        </EuiButton>
      }
      isOpen={isPopoverOpen}
      closePopover={togglePopover}
      panelPaddingSize="none"
      withTitle
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
