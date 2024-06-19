import {
  EuiCollapsibleNav,
  EuiCollapsibleNavGroup,
  EuiFlexGroup,
  EuiFlexItem,
  EuiHeaderSectionItemButton,
  EuiHorizontalRule,
  EuiIcon,
  EuiListGroup,
  EuiListGroupItem,
  EuiTreeView
} from "@elastic/eui";
import urlJoin from "proper-url-join";
import React, { useContext, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useToggle } from "../../hooks";
import { ApplicationsContext, ProjectsContext } from "../../providers";
import { slugify } from "../../utils";

import "./NavDrawer.scss";

export const NavDrawer = ({ docLinks }) => {
  const navigate = useNavigate();
  const { apps, currentApp } = useContext(ApplicationsContext);
  const { currentProject } = useContext(ProjectsContext);

  const mlpLinks = useMemo(() => {
    const isRootApplication = !currentApp;

    return apps.map(a => {
      const isAppActive = a === currentApp;
      const children = !!currentProject
        ? a?.config?.navigation?.map(s => ({
          id: slugify(`${a.name}.${s.label}`),
          label: s.label,
          callback: () => {
            const dest = urlJoin(
              a.homepage,
              "projects",
              currentProject.id,
              s.destination
            );

            isAppActive ? navigate(dest) : (window.location.href = dest);
          },
          className: "euiTreeView__node---small---subsection"
        }))
        : undefined;

      return {
        id: slugify(a.name),
        label: a.name,
        icon: <EuiIcon type={a.config.icon} />,
        isExpanded: isAppActive || isRootApplication,
        className: isAppActive
          ? "euiTreeView__node---small---active"
          : "euiTreeView__node---small",

        callback: () =>
          !children || !currentProject
            ? (window.location.href = !!currentProject
              ? urlJoin(a.homepage, "projects", currentProject.id)
              : a.homepage)
            : {},
        children: children
      };
    });
  }, [apps, currentApp, navigate, currentProject]);

  const adminLinks = [
    {
      label: "Project Settings",
      iconType: "managementApp",
      href: `/projects/${currentProject?.id}/settings`
    }
  ];

  const [navIsOpen, setNavIsOpen] = useState(
    JSON.parse(String(localStorage.getItem("mlp-navIsDocked"))) || false
  );
  const [navIsDocked, setNavIsDocked] = useState(
    JSON.parse(String(localStorage.getItem("mlp-navIsDocked"))) || false
  );

  const [appExpanded, toggleAppExpanded] = useToggle(true);
  const [docsExpanded, toggleDocsExpanded] = useToggle(true);

  return (
    <EuiCollapsibleNav
      aria-label="Main navigation"
      isOpen={navIsOpen}
      isDocked={navIsDocked}
      showButtonIfDocked={true}
      closeButtonProps={{ style: { display: "none" } }}
      button={
        <EuiHeaderSectionItemButton
          aria-label="Toggle main navigation"
          onClick={() => setNavIsOpen(!navIsOpen)}>
          <EuiIcon type={"menu"} size="m" aria-hidden="true" />
        </EuiHeaderSectionItemButton>
      }
      onClose={() => setNavIsOpen(false)}>
      <EuiFlexGroup
        gutterSize="none"
        direction="column"
        style={{ height: "100%" }}>
        <EuiFlexItem grow={true}>
          {mlpLinks.length > 0 && (
            <EuiCollapsibleNavGroup
              isCollapsible={true}
              initialIsOpen={appExpanded}
              onToggle={toggleAppExpanded}
              iconType="apps"
              title="Products">
              <EuiTreeView
                aria-label="MLP apps"
                items={mlpLinks}
                className="mlpApplications"
              />
            </EuiCollapsibleNavGroup>
          )}

          {docLinks.length > 0 && (
            <EuiCollapsibleNavGroup
              title="Learn"
              iconType="training"
              isCollapsible={true}
              initialIsOpen={docsExpanded}
              onToggle={toggleDocsExpanded}>
              <EuiListGroup
                aria-label="Learn"
                listItems={docLinks}
                maxWidth="none"
                color="subdued"
                gutterSize="s"
                size="s"
              />
            </EuiCollapsibleNavGroup>
          )}
        </EuiFlexItem>

        <EuiFlexItem grow={false}>
          {!!currentProject && (
            <EuiCollapsibleNavGroup>
              <EuiListGroup
                listItems={adminLinks}
                color="subdued"
                size="s"
                gutterSize="none"
              />
            </EuiCollapsibleNavGroup>
          )}
        </EuiFlexItem>

        <EuiHorizontalRule margin="none" />

        <EuiFlexItem grow={false}>
          <EuiCollapsibleNavGroup>
            <EuiListGroupItem
              size="s"
              color="subdued"
              label={`${navIsDocked ? "Undock" : "Dock"} navigation`}
              onClick={() => {
                setNavIsDocked(!navIsDocked);
                localStorage.setItem(
                  "mlp-navIsDocked",
                  JSON.stringify(!navIsDocked)
                );
              }}
              iconType={navIsDocked ? "lock" : "lockOpen"}
            />
          </EuiCollapsibleNavGroup>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiCollapsibleNav>
  );
};
