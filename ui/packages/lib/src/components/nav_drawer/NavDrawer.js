import React, { useContext, useMemo, useState } from "react";
import {
  EuiHorizontalRule,
  EuiCollapsibleNav,
  EuiCollapsibleNavGroup,
  EuiListGroup,
  EuiListGroupItem,
  EuiFlexGroup,
  EuiHeaderSectionItemButton,
  EuiIcon,
  EuiFlexItem,
  EuiSpacer,
  EuiTreeView
} from "@elastic/eui";
import { useToggle } from "../../hooks";
import { slugify } from "../../utils";
import urlJoin from "proper-url-join";
import { useNavigate } from "react-router-dom";
import { ApplicationsContext, ProjectsContext } from "../../providers";

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
        ? a?.config?.sections.map(s => ({
            id: slugify(`${a.name}.${s.name}`),
            label: s.name,
            callback: () => {
              const dest = urlJoin(
                a.href,
                "projects",
                currentProject.id,
                s.href
              );

              isAppActive ? navigate(dest) : (window.location.href = dest);
            },
            className: "euiTreeView__node---small---subsection"
          }))
        : undefined;

      return {
        id: slugify(a.name),
        label: a.name,
        icon: <EuiIcon type={a.icon} />,
        isExpanded: isAppActive || isRootApplication,
        className: isAppActive
          ? "euiTreeView__node---small---active"
          : "euiTreeView__node---small",

        callback: () =>
          !children || !currentProject
            ? (window.location.href = !!currentProject
                ? urlJoin(a.href, "projects", currentProject.id)
                : a.href)
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
          <EuiSpacer size="s" />
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
