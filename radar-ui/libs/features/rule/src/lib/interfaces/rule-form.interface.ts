import { RuleSeverity, RuleType } from '@cs/domains/rule';

export interface RuleFilters {
    name: string;
    type: RuleType[];
    notifySeverity: RuleSeverity[];
    blockSeverity: RuleSeverity[];
}
