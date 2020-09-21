import React, { useEffect, useState } from "react";
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiIcon,
  EuiPage,
  EuiPageBody,
  EuiPageContent,
  EuiPageHeader,
  EuiPageHeaderSection,
  EuiSideNav,
  EuiTitle
} from "@elastic/eui";
import { slugify } from "@gojek/mlp-ui/src/utils";
import { navigation, pages } from "./index";
import { Accounts } from "./Accounts";
import queryString from "query-string";

export const Settings = ({ section, ...props }) => {
  const [isSideNavOpenOnMobile, setSideNavOpenOnMobile] = useState(true);
  const [items, setItems] = useState(navigation);
  const [page, setPage] = useState(null);

  const toggleOpenOnMobile = () =>
    setSideNavOpenOnMobile(!isSideNavOpenOnMobile);

  const search = queryString.parse(props.location.search);

  useEffect(() => {
    let activePage = null;

    const items = navigation.map(({ items, ...rest }) => ({
      ...rest,
      items: items.map(({ id: itemId, name: itemName, ...rest }) => {
        const isSelected = slugify(itemName) === section;

        if (isSelected) {
          activePage = pages[itemId];
        }

        return {
          id: itemId,
          name: itemName,
          isSelected: isSelected,
          ...rest
        };
      })
    }));
    if (activePage) {
      setPage(activePage);
    }
    setItems(items);
  }, [section, setItems, setPage]);
  return (
    <div className="guideBody">
      <EuiPage className="guidePage">
        <EuiPageBody>
          <EuiPageHeader>
            <EuiPageHeaderSection>
              <EuiTitle size="l">
                <h1>
                  <EuiIcon type="gear" size="xl" /> Settings
                </h1>
              </EuiTitle>
            </EuiPageHeaderSection>
          </EuiPageHeader>
          <EuiFlexGroup>
            <EuiFlexItem grow={6}>
              <EuiPageContent>
                {!section && !page && (
                  <Accounts search={search} />
                ) /* default */}
                {section === "connected-accounts" && (
                  <Accounts search={search} />
                )}
              </EuiPageContent>
            </EuiFlexItem>
            <EuiFlexItem grow={2}>
              <EuiSideNav
                mobileTitle="Settings Menu"
                toggleOpenOnMobile={toggleOpenOnMobile}
                isOpenOnMobile={isSideNavOpenOnMobile}
                items={items}
              />
            </EuiFlexItem>
            <EuiFlexItem grow={2} />
          </EuiFlexGroup>
        </EuiPageBody>
      </EuiPage>
    </div>
  );
};
