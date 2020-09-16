import React, { Fragment } from "react";
import { EuiButton, EuiEmptyPrompt } from "@elastic/eui";

import * as Sentry from "@sentry/browser";

export class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { error: null, errorInfo: null, eventId: null };
  }

  componentDidCatch(error, errorInfo) {
    this.setState({
      error: error,
      errorInfo: errorInfo
    });

    Sentry.withScope(scope => {
      scope.setExtras(errorInfo);
      const eventId = Sentry.captureException(error);
      this.setState({ eventId });
    });
  }

  render() {
    return !!this.state.errorInfo ? (
      <EuiEmptyPrompt
        iconType="editorStrike"
        title={<h2>Something Went Wrong</h2>}
        body={
          <Fragment>
            <p>
              Something wen't wrong with the UI, no worries, we will fix this as
              soon possible
            </p>
            <p>Please report this to maintainer my friend</p>
          </Fragment>
        }
        actions={
          <EuiButton
            color="primary"
            fill
            onClick={() =>
              Sentry.showReportDialog({ eventId: this.state.eventId })
            }>
            Report Feedback
          </EuiButton>
        }
      />
    ) : (
      this.props.children
    );
  }
}
