import React from "react";
import FormValidationContext from "./context";
import { MultiSectionFormValidationContextProvider } from "./multi_section_provider";

export const FormValidationContextProvider = ({
  schema,
  context,
  onSubmit,
  children
}) => (
  <MultiSectionFormValidationContextProvider
    onSubmit={onSubmit}
    schemas={[schema]}
    contexts={[context]}>
    <FormValidationContext.Consumer>
      {({ onSubmit, isSubmitting, errors }) => (
        <FormValidationContext.Provider
          value={{ onSubmit, isSubmitting, errors: errors[0] }}>
          {children}
        </FormValidationContext.Provider>
      )}
    </FormValidationContext.Consumer>
  </MultiSectionFormValidationContextProvider>
);
