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
import ApplicationsContext from "../../providers/application/context";
import { CurrentProjectContext } from "../../providers/project";
import { useToggle } from "../../hooks";
import { navigate } from "@reach/router";
import { get, slugify } from "../../utils";
import urlJoin from "proper-url-join";

import "./NavDrawer.scss";

export const NavDrawer = ({ homeUrl = "/", docLinks }) => {
  const { projectId } = useContext(CurrentProjectContext);
  const { apps } = useContext(ApplicationsContext);

  const isRootApplication = homeUrl === "/";

  const mlpLinks = useMemo(
    () =>
      apps.map(a => {
        const isAppActive = a.href === homeUrl;

        const children =
          !!projectId && get(a, "config.sections")
            ? a.config.sections.map(s => ({
                id: slugify(`${a.name}.${s.name}`),
                label: s.name,
                callback: () =>
                  navigate(urlJoin(a.href, "projects", projectId, s.href)),
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
            !children || !projectId
              ? navigate(projectId ? `${a.href}/projects/${projectId}` : a.href)
              : {},
          children: children
        };
      }),
    [apps, homeUrl, projectId, isRootApplication]
  );

  const adminLinks = [
    {
      label: "Project Settings",
      iconType: "managementApp",
      href: `/projects/${projectId}/settings`
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
      showCloseButton={false}
      showButtonIfDocked={true}
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
          {projectId && (
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
