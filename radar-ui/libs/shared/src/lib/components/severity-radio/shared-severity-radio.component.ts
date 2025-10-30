import {
    AbstractControl,
    ControlValueAccessor,
    NG_VALIDATORS,
    NG_VALUE_ACCESSOR,
    ValidationErrors,
    Validators
} from '@angular/forms';
import { Component, Input } from '@angular/core';

import { RULE_SEVERITIES, RuleSeverity } from '@cs/domains/rule';

@Component({
    selector: 'cs-severity-radio-component',
    templateUrl: './shared-severity-radio.component.html',
    styleUrl: './shared-severity-radio.component.scss',
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: SharedSeverityRadioComponent,
            multi: true
        },
        {
            provide: NG_VALIDATORS,
            useExisting: SharedSeverityRadioComponent,
            multi: true
        }
    ]
})
export class SharedSeverityRadioComponent implements ControlValueAccessor {
    @Input() id?: string;

    @Input() testLocator?: string;

    @Input() noneLabelLocalizationKey?: string;

    isTouched = false;

    isDisabled = false;

    severity = RuleSeverity.NONE;

    readonly ruleSeverityOptions = RULE_SEVERITIES;

    /* eslint @typescript-eslint/no-empty-function: "off" */
    onChange = (severity: RuleSeverity) => {};

    /* eslint @typescript-eslint/no-empty-function: "off" */
    onTouched = () => {};

    registerOnChange(fn: any) {
        this.onChange = fn;
    }

    registerOnTouched(fn: any) {
        this.onTouched = fn;
    }

    markAsTouched() {
        if (!this.isTouched) {
            this.isTouched = true;
            this.onTouched();
        }
    }

    setDisabledState(isDisabled: boolean) {
        this.isDisabled = isDisabled;
    }

    writeValue(severity?: RuleSeverity | null) {
        if (severity) {
            this.severity = severity;
        }
    }

    validate(control: AbstractControl): ValidationErrors | null {
        return control.hasValidator(Validators.required) && !control.value ? { required: true } : null;
    }

    changeSeverity(severity: RuleSeverity) {
        if (!this.isDisabled) {
            this.severity = severity;
            this.onChange(this.severity);
        }

        this.markAsTouched();
    }
}
