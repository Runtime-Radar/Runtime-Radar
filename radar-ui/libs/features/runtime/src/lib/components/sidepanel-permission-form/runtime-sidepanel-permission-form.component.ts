import { KbqTagInputEvent } from '@koobiq/components/tags';
import { AfterViewInit, ChangeDetectionStrategy, Component, Inject } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { KBQ_SIDEPANEL_DATA, KbqSidepanelRef } from '@koobiq/components/sidepanel';
import { Observable, debounceTime, distinctUntilChanged, map, startWith } from 'rxjs';

import { FORM_SEPARATOR_KEY_CODES, FORM_VALIDATION_REG_EXP, FormScheme, CoreUtilsService as utils } from '@cs/core';

import { RuntimeSettingPermissionForm } from '../../interfaces/runtime-form.interface';
import { RuntimeSidepanelPermissionFormProps } from '../../interfaces/runtime-sidepanel.interface';

@Component({
    templateUrl: './runtime-sidepanel-permission-form.component.html',
    styleUrl: './runtime-sidepanel-permission-form.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuntimeFeatureSidepanelPermissionFormComponent implements AfterViewInit {
    readonly form: FormGroup<FormScheme<RuntimeSettingPermissionForm, never, 'namespaces' | 'pods' | 'labels'>> =
        this.formBuilder.group({
            isAllowedType: [true],
            namespaces: this.formBuilder.array<string>([]),
            pods: this.formBuilder.array<string>([]),
            labels: this.formBuilder.array<string>([])
        });

    readonly isFormValid$: Observable<boolean> = this.form.valueChanges.pipe(
        startWith(this.form.value),
        /* eslint @typescript-eslint/no-magic-numbers: "off" */
        debounceTime(250),
        distinctUntilChanged(),
        map(() => {
            const values = utils.getFormValues<RuntimeSettingPermissionForm>(this.form.controls);

            return (
                utils.isFormValid(this.form.controls) &&
                !!(values.namespaces.length || values.pods.length || values.labels.length)
            );
        })
    );

    get namespacesControl(): FormArray {
        return this.form.get('namespaces') as FormArray;
    }

    get podsControl(): FormArray {
        return this.form.get('pods') as FormArray;
    }

    get labelsControl(): FormArray {
        return this.form.get('labels') as FormArray;
    }

    readonly separatorKeyCodes = FORM_SEPARATOR_KEY_CODES;

    readonly formValidationRegExp = FORM_VALIDATION_REG_EXP;

    constructor(
        private readonly formBuilder: FormBuilder,
        private readonly sidepanelRef: KbqSidepanelRef,
        @Inject(KBQ_SIDEPANEL_DATA) public readonly props: Partial<RuntimeSidepanelPermissionFormProps>
    ) {}

    ngAfterViewInit() {
        utils.setArrayControlValue(
            this.namespacesControl,
            this.props.namespaces,
            this.formBuilder,
            FORM_VALIDATION_REG_EXP.OBJECT_NAME
        );

        utils.setArrayControlValue(
            this.podsControl,
            this.props.pods,
            this.formBuilder,
            FORM_VALIDATION_REG_EXP.REG_EXP
        );

        utils.setArrayControlValue(
            this.labelsControl,
            this.props.labels,
            this.formBuilder,
            FORM_VALIDATION_REG_EXP.KEY_VALUE_PAIR
        );

        if (this.props.isEdit) {
            this.form.patchValue({
                isAllowedType: this.props.isAllowedType
            });
        }
    }

    addEntity(event: KbqTagInputEvent, control: FormArray, validationRegExp: RegExp) {
        const value = event.value.trim();

        if (value) {
            control.push(this.formBuilder.control(value, Validators.pattern(validationRegExp)));
            event.input.value = '';
        }
    }

    removeEntity(control: FormArray, id: number) {
        control.removeAt(id);
    }

    confirm() {
        if (this.form.valid) {
            const formValues = utils.getFormValues<RuntimeSettingPermissionForm>(this.form.controls);
            this.sidepanelRef.close(utils.getTrimmedFormValues<RuntimeSettingPermissionForm>(formValues));
        }
    }

    cancel() {
        this.sidepanelRef.close(undefined);
    }
}
