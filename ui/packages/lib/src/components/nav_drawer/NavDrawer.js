import React, { useContext, useEffect, useMemo, useState } from "react";
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

export const NavDrawer = ({ homeUrl = "/", appLinks }) => {
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

  const docLinks = [
    {
      label: "Merlin User Guide",
      href:
        "https://go-jek.atlassian.net/wiki/spaces/DSP/pages/1450639714/Merlin+User+Guide"
    },
    {
      label: "Turing User Guide",
      href:
        "https://go-jek.atlassian.net/wiki/spaces/DSP/pages/1757883181/Turing+User+Documentation"
    },
    {
      label: "Feast User Guide",
      hrefh: "https://docs.feast.dev/user-guide/overview"
    }
  ];

  const [navIsOpen, setNavIsOpen] = useState(true);
  const [navIsDocked, setNavIsDocked] = useState(true);

  useEffect(() => {
    setNavIsDocked(navIsOpen);
  }, [navIsOpen]);

  const [appExpanded, toggleAppExpanded] = useToggle(true);
  const [docsExpanded, toggleDocsExpanded] = useToggle(true);
  return (
    <EuiCollapsibleNav
      aria-label="Main navigation"
      isOpen={navIsOpen}
      showCloseButton={false}
      isDocked={navIsDocked}
      showButtonIfDocked={true}
      style={{ top: "49px", height: "calc(100% - 49px)", width: "240px" }}
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
          {appLinks && (
            <>
              <EuiCollapsibleNavGroup listItems={appLinks} />
              <EuiHorizontalRule margin="none" />
            </>
          )}
          <EuiCollapsibleNavGroup
            isCollapsible={true}
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

          <EuiCollapsibleNavGroup
            title="Learn"
            iconType="training"
            isCollapsible={true}
            initialIsOpen={docsExpanded}
            onToggle={toggleDocsExpanded}>
            <EuiListGroup
              aria-label="Learn" // A11y : EuiCollapsibleNavGroup can't correctly pass the `title` as the `aria-label` to the right HTML element, so it must be added manually
              listItems={docLinks}
              maxWidth="none"
              color="subdued"
              gutterSize="s"
              size="s"
            />
          </EuiCollapsibleNavGroup>
        </EuiFlexItem>

        <EuiHorizontalRule margin="none" />

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
        <EuiFlexItem grow={false}>
          <EuiShowFor sizes={["l", "xl"]}>
            <EuiCollapsibleNavGroup>
              <EuiListGroupItem
                size="s"
                color="subdued"
                label={`${navIsDocked ? "Undock" : "Dock"} navigation`}
                onClick={() => {
                  setNavIsDocked(!navIsDocked);
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
