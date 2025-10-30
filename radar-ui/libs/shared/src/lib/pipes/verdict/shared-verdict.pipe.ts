import { Pipe, PipeTransform } from '@angular/core';

import { RULE_VERDICTS, RuleVerdict, SEVERITY_UNDEFINED_LOCALIZATION_KEY } from '@cs/domains/rule';

@Pipe({
    name: 'verdictLocalization',
    pure: false
})
export class SharedVerdictLocalizationPipe implements PipeTransform {
    transform(verdict?: RuleVerdict | null, noneLabelLocalizationKey?: string): string {
        const value = RULE_VERDICTS.find((item) => item.id === verdict);

        if (verdict === RuleVerdict.NONE && noneLabelLocalizationKey) {
            return noneLabelLocalizationKey;
        }

        return value ? value.localizationKey : SEVERITY_UNDEFINED_LOCALIZATION_KEY;
    }
}
