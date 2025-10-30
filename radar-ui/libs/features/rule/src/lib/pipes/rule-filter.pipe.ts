import { Pipe, PipeTransform } from '@angular/core';

import { Rule, RuleSeverity } from '@cs/domains/rule';

import { RuleFilters } from '../interfaces/rule-form.interface';

@Pipe({
    name: 'ruleFilter'
})
export class RuleFeatureFilterPipe implements PipeTransform {
    transform(values?: Rule[] | null, filters?: RuleFilters): Rule[] {
        if (!values) {
            return [];
        }

        if (!filters) {
            return values;
        }

        const keys = Object.keys(filters).reduce((acc, key: string) => {
            const value = filters[key as keyof RuleFilters];

            return value !== null && !!value.length ? [...acc, key] : acc;
        }, [] as string[]);

        return values.filter((item) =>
            keys.every((key) => {
                if (key === 'name') {
                    return item.name.toLowerCase().indexOf(filters.name.toLowerCase()) === 0;
                } else if (key === 'type') {
                    return filters.type.includes(item.type);
                }

                const notify = item.rule.notify?.severity || RuleSeverity.NONE;

                if (key === 'notifySeverity') {
                    return filters.notifySeverity.includes(notify);
                }

                return true;
            })
        );
    }
}
