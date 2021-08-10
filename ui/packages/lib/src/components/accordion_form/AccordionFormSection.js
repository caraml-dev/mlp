import React from "react";
import { Element } from "react-scroll";
import { EuiAccordion } from "@elastic/eui";
import { slugify } from "../../utils";
import { FormValidationContext } from "../form/validation";
import { isSectionInvalid } from "./functions";

export const AccordionFormSection = ({ section, errors, renderTitle }) => (
  <Element name={slugify(section.title)}>
    <FormValidationContext.Provider value={{ errors }}>
      <EuiAccordion
        id={`accordion-form-${slugify(section.title)}`}
        initialIsOpen={true}
        forceState={isSectionInvalid(errors) ? "open" : undefined}
        buttonClassName="euiAccordionForm__button"
        buttonContent={
          renderTitle
            ? renderTitle(section.title, section.iconType)
            : section.title
        }
        extraAction={section.extraAction}
        paddingSize="s">
        {section.children}
      </EuiAccordion>
    </FormValidationContext.Provider>
  </Element>
);
