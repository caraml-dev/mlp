import React, { useContext, useEffect, useMemo, useState } from "react";
import {
  EuiHorizontalRule,
  EuiNavDrawer,
  EuiNavDrawerGroup,
  EuiFlexGroup,
  EuiFlexItem,
  EuiSpacer
} from "@elastic/eui";
import ApplicationsContext from "../../providers/application/context";
import { CurrentProjectContext } from "../../providers/project";
import "./NavDrawer.scss";

const HEADER_HEIGHT = 49;

export const NavDrawer = ({ homeUrl = "/", appLinks }) => {
  const { projectId } = useContext(CurrentProjectContext);
  const { apps } = useContext(ApplicationsContext);
  const [topPosition, setTopPosition] = useState(HEADER_HEIGHT);

  const handleScroll = () => {
    const position = Math.max(HEADER_HEIGHT - window.pageYOffset, 0);
    setTopPosition(position);
  };

  useEffect(() => {
    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, []);

  const mlpLinks = useMemo(
    () =>
      apps
        ? apps
            .filter(a => a.href !== homeUrl)
            .map(a => ({
              label: a.name,
              iconType: a.icon,
              href: projectId ? `${a.href}/projects/${projectId}` : a.href
            }))
        : [],
    [apps, homeUrl, projectId]
  );

  const adminLinks = [
    {
      label: "Project Settings",
      iconType: "managementApp",
      href: `/projects/${projectId}/settings`
    }
  ];

  return (
    <EuiNavDrawer>
      <EuiFlexGroup
        gutterSize="none"
        direction="column"
        style={{ height: "100%" }}>
        <EuiFlexItem grow={true} style={{ marginTop: `${topPosition}px` }}>
          {appLinks && (
            <>
              <EuiNavDrawerGroup showToolTips={true} listItems={appLinks} />
              <EuiHorizontalRule margin="none" />
            </>
          )}
          <EuiNavDrawerGroup showToolTips={true} listItems={mlpLinks} />
          <EuiHorizontalRule margin="none" />
        </EuiFlexItem>
        <EuiFlexItem grow={false}>
          {projectId && (
            <EuiNavDrawerGroup showToolTips={true} listItems={adminLinks} />
          )}
          <EuiSpacer size="s" />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiNavDrawer>
  );
};
