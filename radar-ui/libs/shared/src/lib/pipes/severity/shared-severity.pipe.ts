import { Pipe, PipeTransform } from '@angular/core';

import { RULE_SEVERITIES, RuleSeverity, SEVERITY_UNDEFINED_LOCALIZATION_KEY } from '@cs/domains/rule';

@Pipe({
    name: 'severityLocalization',
    pure: false
})
export class SharedSeverityLocalizationPipe implements PipeTransform {
    transform(severity?: RuleSeverity | null, noneLabelLocalizationKey?: string): string {
        const value = RULE_SEVERITIES.find((item) => item.id === severity);

        if (severity === RuleSeverity.NONE && noneLabelLocalizationKey) {
            return noneLabelLocalizationKey;
        }

        return value ? value.localizationKey : SEVERITY_UNDEFINED_LOCALIZATION_KEY;
    }
}
