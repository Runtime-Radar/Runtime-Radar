import { KBQ_SIDEPANEL_DATA } from '@koobiq/components/sidepanel';
import { ChangeDetectionStrategy, Component, Inject } from '@angular/core';

import { RULE_SEVERITIES, RULE_SEVERITY_ORDER, RuleSeverityOption } from '@cs/domains/rule';

import { RuntimeSidepanelIncidentProps } from '../../interfaces/runtime-sidepanel.interface';

@Component({
    templateUrl: './runtime-sidepanel-incident.component.html',
    styleUrl: './runtime-sidepanel-incident.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuntimeFeatureSidepanelIncidentComponent {
    readonly ruleSeverityOptions: RuleSeverityOption[] = [...RULE_SEVERITIES].sort(
        (a, b) => RULE_SEVERITY_ORDER[a.id] - RULE_SEVERITY_ORDER[b.id]
    );

    constructor(@Inject(KBQ_SIDEPANEL_DATA) public readonly props: RuntimeSidepanelIncidentProps) {}
}
