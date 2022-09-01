import React from "react";
import { EuiButton, EuiPageTemplate } from "@elastic/eui";
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
      <EuiPageTemplate>
        <EuiPageTemplate.EmptyPrompt
          color="subdue"
          iconType="editorStrike"
          title={<h2>Something Went Wrong</h2>}
          body={
            <>
              <p>
                Something went wrong with the UI, no worries, we will fix
                <br />
                this as soon as possible
              </p>
              <p>Please report this to maintainer my friend</p>
            </>
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
      </EuiPageTemplate>
    ) : (
      this.props.children
    );
  }
}

export const ErrorBox = ({ error }) => {
  const Error = () => {
    throw error;
  };
  return (
    <ErrorBoundary>
      <Error />
    </ErrorBoundary>
  );
};
