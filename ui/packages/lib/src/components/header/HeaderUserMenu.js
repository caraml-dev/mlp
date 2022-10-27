import React, { Fragment, useCallback, useState } from "react";

import {
  EuiAvatar,
  EuiFlexGroup,
  EuiFlexItem,
  EuiHeaderSectionItemButton,
  EuiHorizontalRule,
  EuiContextMenuPanel,
  EuiContextMenuItem,
  EuiText,
  EuiTitle,
  EuiPopover
} from "@elastic/eui";
import PropTypes from "prop-types";
import "./HeaderUserMenu.scss";

export const HeaderUserMenu = ({ profileObj, logout, children }) => {
  const [isOpen, setOpen] = useState(false);

  const togglePopover = useCallback(() => setOpen(isOpen => !isOpen), []);

  const button = (
    <EuiHeaderSectionItemButton
      aria-controls="headerUserMenu"
      aria-expanded={isOpen}
      aria-haspopup="true"
      aria-label="Account menu"
      onClick={togglePopover}>
      <EuiAvatar
        imageUrl={profileObj.picture}
        name={profileObj.name}
        size="s"
      />
    </EuiHeaderSectionItemButton>
  );

  const horizontalSeparator = (
    <EuiHorizontalRule
      margin="none"
      style={{ width: "95%", marginLeft: "auto", marginRight: "auto" }}
    />
  );

  return (
    <EuiPopover
      id="headerUserMenu"
      button={button}
      isOpen={isOpen}
      anchorPosition="downRight"
      closePopover={togglePopover}
      panelPaddingSize="s">
      <EuiContextMenuPanel className="euiContextMenuPanel--headerUserMenu">
        <EuiFlexGroup
          gutterSize="m"
          className="euiHeaderProfile"
          responsive={false}>
          <EuiFlexItem grow={false}>
            <EuiAvatar
              imageUrl={profileObj.picture}
              name={profileObj.name}
              size="xl"
            />
          </EuiFlexItem>

          <EuiFlexItem>
            <EuiTitle size="xs">
              <span>{profileObj.name}</span>
            </EuiTitle>
            <EuiText>
              <p>{profileObj.email}</p>
            </EuiText>
          </EuiFlexItem>
        </EuiFlexGroup>

        {children ? (
          <Fragment>
            {horizontalSeparator}

            {children}
          </Fragment>
        ) : null}

        {horizontalSeparator}

        <EuiContextMenuItem
          key="log_out"
          icon="exit"
          toolTipPosition="right"
          onClick={logout}>
          Log out
        </EuiContextMenuItem>
      </EuiContextMenuPanel>
    </EuiPopover>
  );
};

HeaderUserMenu.propTypes = {
  profileObj: PropTypes.shape({
    email: PropTypes.string,
    name: PropTypes.string,
    picture: PropTypes.string
  }).isRequired,
  logout: PropTypes.func.isRequired
};
