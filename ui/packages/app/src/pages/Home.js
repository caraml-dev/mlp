import React, { Fragment } from "react";
import {
  EuiButton,
  EuiCard,
  EuiFlexGroup,
  EuiFlexItem,
  EuiIcon,
  EuiFlexGrid,
  EuiEmptyPrompt
} from "@elastic/eui";
import { ApplicationsContext } from "@gojek/mlp-ui";

const Home = () => {
  return (
    <Fragment>
      <EuiEmptyPrompt
        iconType="graphApp"
        title={<h2>Machine Learning Platform</h2>}
      />
      <EuiFlexGroup style={{ paddingLeft: 50, paddingRight: 50 }}>
        <EuiFlexItem grow={false} style={{ margin: "auto" }}>
          <EuiFlexGrid columns={3}>
            <ApplicationsContext.Consumer>
              {({ apps }) =>
                apps.map((app, idx) => {
                  return (
                    <EuiFlexItem key={idx}>
                      <EuiCard
                        icon={<EuiIcon size="xxl" type={app.icon} />}
                        title={app.name}
                        description={app.description}
                        footer={
                          <EuiButton href={app.href}>
                            Go to {app.name}
                          </EuiButton>
                        }
                        betaBadgeLabel={app.is_in_beta ? "Beta" : undefined}
                      />
                    </EuiFlexItem>
                  );
                })
              }
            </ApplicationsContext.Consumer>
          </EuiFlexGrid>
        </EuiFlexItem>
      </EuiFlexGroup>
    </Fragment>
  );
};

export default Home;
