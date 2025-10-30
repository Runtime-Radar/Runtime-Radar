import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ChangeDetectionStrategy, Component, DestroyRef, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Observable, debounceTime, distinctUntilChanged, map, tap } from 'rxjs';

import { FormScheme, CoreUtilsService as utils } from '@cs/core';
import { RULE_SEVERITIES, RULE_TYPE, RuleSeverity, RuleType } from '@cs/domains/rule';

import { RuleFilters } from '../../interfaces/rule-form.interface';

@Component({
    selector: 'cs-rule-feature-panel-filter-component',
    templateUrl: './rule-panel-filter.component.html',
    styleUrl: './rule-panel-filter.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuleFeaturePanelFilterComponent implements OnInit {
    @Output() filterChange = new EventEmitter<RuleFilters>();

    readonly filterForm: FormGroup<FormScheme<RuleFilters>> = this.formBuilder.group({
        name: [''],
        type: [[] as RuleType[]],
        notifySeverity: [[] as RuleSeverity[]],
        blockSeverity: [[] as RuleSeverity[]]
    });

    private readonly onFormValidChanges$: Observable<RuleFilters> = this.filterForm.valueChanges.pipe(
        /* eslint @typescript-eslint/no-magic-numbers: "off" */
        debounceTime(500),
        distinctUntilChanged(),
        map(() => utils.getFormValues<RuleFilters>(this.filterForm.controls)),
        tap((values) => {
            this.filterChange.emit(values);
        })
    );

    readonly ruleTypeOptions = RULE_TYPE;

    readonly ruleSeverityOptions = RULE_SEVERITIES;

    constructor(
        private readonly destroyRef: DestroyRef,
        private readonly formBuilder: FormBuilder
    ) {}

    ngOnInit() {
        this.onFormValidChanges$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe();
    }
}
