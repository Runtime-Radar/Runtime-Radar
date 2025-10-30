import { ChangeDetectionStrategy, Component, Input, booleanAttribute } from '@angular/core';

import { RuleSeverity, RuleVerdict, SEVERITY_UNDEFINED_LOCALIZATION_KEY } from '@cs/domains/rule';

@Component({
    selector: 'cs-severity-label-component',
    templateUrl: './shared-severity-label.component.html',
    styleUrl: './shared-severity-label.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class SharedSeverityLabelComponent {
    @Input() severity?: RuleSeverity | null;

    @Input() verdict?: RuleVerdict | null;

    @Input() noneLabelLocalizationKey?: string;

    @Input({ transform: booleanAttribute }) isWidthAuto = false;

    @Input({ transform: booleanAttribute }) isLabelColored = false;

    readonly undefinedLocalizationKey = SEVERITY_UNDEFINED_LOCALIZATION_KEY;
}
