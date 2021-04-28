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
  EuiShowFor
} from "@elastic/eui";
import ApplicationsContext from "../../providers/application/context";
import { CurrentProjectContext } from "../../providers/project";
import { useToggle } from "../../hooks/useToggle";
import "./NavDrawer.scss";

export const NavDrawer = ({ homeUrl = "/", docLinks }) => {
  const { projectId } = useContext(CurrentProjectContext);
  const { apps } = useContext(ApplicationsContext);

  const mlpLinks = useMemo(
    () =>
      apps.map(a => ({
        label: a.name,
        iconType: a.icon,
        href: projectId ? `${a.href}/projects/${projectId}` : a.href,
        isActive: a.href === homeUrl
      })),
    [apps, homeUrl, projectId]
  );

  const adminLinks = [
    {
      label: "Project Settings",
      iconType: "managementApp",
      href: `/projects/${projectId}/settings`
    }
  ];

  const [navIsOpen, setNavIsOpen] = useState(JSON.parse(String(localStorage.getItem('mlp-navIsDocked'))) || false);
  const [navIsDocked, setNavIsDocked] = useState(JSON.parse(String(localStorage.getItem('mlp-navIsDocked'))) || false);

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
              isCollapsible={false}
              initialIsOpen={appExpanded}
              onToggle={toggleAppExpanded}
              iconType="apps"
              title="Products">
              <EuiListGroup
                aria-label="MLP apps"
                listItems={mlpLinks}
                maxWidth="none"
                color="text"
                gutterSize="s"
                size="s"
              />
            </EuiCollapsibleNavGroup>
          )}

          {docLinks.length > 0 && (
            <EuiCollapsibleNavGroup
              title="Learn"
              iconType="training"
              isCollapsible={false}
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
          <EuiShowFor sizes={["l", "xl"]}>
            <EuiCollapsibleNavGroup>
              <EuiListGroupItem
                size="s"
                color="subdued"
                label={`${navIsDocked ? "Undock" : "Dock"} navigation`}
                onClick={() => {
                  setNavIsDocked(!navIsDocked);
                  localStorage.setItem(
                    'mlp-navIsDocked',
                    JSON.stringify(!navIsDocked)
                  );
                }}
                iconType={navIsDocked ? "lock" : "lockOpen"}
              />
            </EuiCollapsibleNavGroup>
          </EuiShowFor>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiCollapsibleNav>
  );
};
