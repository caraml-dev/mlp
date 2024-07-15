import React, { useRef } from "react";
import { EuiFlexGroup, EuiFlexItem, EuiSpacer } from "@elastic/eui";
import Sticky from "react-stickynode";
import { StepActions } from "../multi_steps_form";
import FormValidationContext from "../form/validation/context";
import { MultiSectionFormValidationContextProvider } from "../form";
import { AccordionFormSideNav } from "./AccordionFormSideNav";
import { AccordionFormScrollController } from "./AccordionFormScrollController";
import { AccordionFormSection } from "./AccordionFormSection";
import { useDimension } from "../../hooks";

import "./AccordionForm.scss";

export const AccordionForm = ({
  name,
  sections,
  onCancel,
  onSubmit,
  submitLabel,
  renderTitle
}) => {
  const lastSectionRef = useRef(null);
  const { height: lastSectionHeight } = useDimension(lastSectionRef);

  return (
    <EuiFlexGroup direction="row" gutterSize="none" className="accordionForm">
      <EuiFlexItem grow={false} className="accordionForm--sideNavContainer">
        <Sticky enabled={true}>
          <AccordionFormSideNav name={name} sections={sections} />
        </Sticky>
      </EuiFlexItem>
      <EuiFlexItem grow={true} className="accordionForm--content">
        <MultiSectionFormValidationContextProvider
          onSubmit={onSubmit}
          schemas={sections.map(s => s.validationSchema)}
          contexts={sections.map(s => s.validationContext)}>
          <FormValidationContext.Consumer>
            {({ errors, onSubmit, isSubmitting }) => (
              <EuiFlexGroup
                direction="column"
                gutterSize="none"
                alignItems="center">
                <AccordionFormScrollController sections={sections} />

                {sections.map((section, idx) => (
                  <EuiFlexItem key={idx}>
                    <span
                      ref={
                        idx === sections.length - 1
                          ? lastSectionRef
                          : undefined
                      }>
                      <AccordionFormSection
                        section={section}
                        errors={errors[idx]}
                        renderTitle={renderTitle}
                      />
                    </span>
                  </EuiFlexItem>
                ))}

                <EuiSpacer size="l" />

                <EuiFlexItem
                  // set the minHeight dynamically, based on the height of the last section
                  style={{
                    minHeight: `calc(100vh - ${lastSectionHeight +
                      24 +
                      16}px)`
                  }}>
                  <StepActions
                    submitLabel={submitLabel}
                    onCancel={onCancel}
                    onSubmit={onSubmit}
                    isSubmitting={isSubmitting}
                  />
                </EuiFlexItem>
              </EuiFlexGroup>
            )}
          </FormValidationContext.Consumer>
        </MultiSectionFormValidationContextProvider>
      </EuiFlexItem>
    </EuiFlexGroup>
  );
};
