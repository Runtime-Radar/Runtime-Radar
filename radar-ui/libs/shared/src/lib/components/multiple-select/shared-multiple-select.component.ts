import {
    AbstractControl,
    ControlValueAccessor,
    NG_VALIDATORS,
    NG_VALUE_ACCESSOR,
    ValidationErrors,
    Validators
} from '@angular/forms';
import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';

const MULTIPLE_SELECT_ALL_OPTION_KEY = 'MULTIPLE_SELECT_ALL_OPTION_KEY';

interface AbstractMultipleSelectOption {
    id: string;
    localizationKey: string;
}

@Component({
    selector: 'cs-multiple-select-component',
    templateUrl: './shared-multiple-select.component.html',
    styleUrl: './shared-multiple-select.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: SharedMultipleSelectComponent,
            multi: true
        },
        {
            provide: NG_VALIDATORS,
            useExisting: SharedMultipleSelectComponent,
            multi: true
        }
    ]
})
export class SharedMultipleSelectComponent<T extends AbstractMultipleSelectOption>
    implements ControlValueAccessor, OnInit
{
    @Input({ required: true }) options!: T[];

    @Input({ required: true }) localizationFn!: any;

    @Input() placeholder?: string;

    @Input() id?: string;

    @Input() testId?: string;

    optionItems = 0;

    selectedItems = 0;

    selected: string[] = [];

    isTouched = false;

    isDisabled = false;

    readonly allOptionKey = MULTIPLE_SELECT_ALL_OPTION_KEY;

    ngOnInit() {
        this.optionItems = this.options.length;
    }

    /* eslint @typescript-eslint/no-empty-function: "off" */
    onChange = (selected: string[]) => {};

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

    writeValue(selected?: string[] | null) {
        if (selected) {
            this.selected = selected;
            this.selectedItems = selected.length;
        }
    }

    validate(control: AbstractControl): ValidationErrors | null {
        return control.hasValidator(Validators.required) && !control.value ? { required: true } : null;
    }

    changeSelected(selected: string[]) {
        if (!this.isDisabled) {
            const value = this.getSelectedOption(selected);
            if (value === MULTIPLE_SELECT_ALL_OPTION_KEY) {
                this.selected = this.selectedItems !== this.optionItems ? this.options.map((item) => item.id) : [];
            } else {
                this.selected = selected;
            }

            this.selectedItems = this.selected.length;
            this.onChange(this.selected);
        }

        this.markAsTouched();
    }

    private getSelectedOption(selected: string[]): string | undefined {
        let options = selected.filter((item) => !this.selected.includes(item));
        if (!options.length) {
            options = this.selected.filter((item) => !selected.includes(item));
        }

        return options.at(0);
    }
}
