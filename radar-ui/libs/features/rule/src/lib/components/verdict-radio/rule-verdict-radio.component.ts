import {
    AbstractControl,
    ControlValueAccessor,
    NG_VALIDATORS,
    NG_VALUE_ACCESSOR,
    ValidationErrors,
    Validators
} from '@angular/forms';
import { Component, Input } from '@angular/core';

import { RULE_VERDICTS, RuleVerdict } from '@cs/domains/rule';

@Component({
    selector: 'cs-rule-feature-verdict-radio-component',
    templateUrl: './rule-verdict-radio.component.html',
    // @todo: make radio styles as shared, styles from severity are duplicated now
    styleUrl: './rule-verdict-radio.component.scss',
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: RuleFeatureVerdictRadioComponent,
            multi: true
        },
        {
            provide: NG_VALIDATORS,
            useExisting: RuleFeatureVerdictRadioComponent,
            multi: true
        }
    ]
})
export class RuleFeatureVerdictRadioComponent implements ControlValueAccessor {
    @Input() id?: string;

    @Input() testLocator?: string;

    @Input() noneLabelLocalizationKey?: string;

    isTouched = false;

    isDisabled = false;

    verdict = RuleVerdict.NONE;

    readonly ruleVerdictOptions = RULE_VERDICTS;

    /* eslint @typescript-eslint/no-empty-function: "off" */
    onChange = (verdict: RuleVerdict) => {};

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

    writeValue(verdict?: RuleVerdict | null) {
        if (verdict) {
            this.verdict = verdict;
        }
    }

    validate(control: AbstractControl): ValidationErrors | null {
        return control.hasValidator(Validators.required) && !control.value ? { required: true } : null;
    }

    changeVerdict(verdict: RuleVerdict) {
        if (!this.isDisabled) {
            this.verdict = verdict;
            this.onChange(this.verdict);
        }

        this.markAsTouched();
    }
}
