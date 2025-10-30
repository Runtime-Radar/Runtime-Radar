import { RuleSeverity, RuleType, RuleVerdict } from './contract/rule-contract.interface';

export type RuleSeverityOrder = {
    [key in RuleSeverity]: number;
};

export interface RuleSeverityOption {
    id: RuleSeverity;
    localizationKey: string;
    testId: string;
}

export interface RuleVerdictOption {
    id: RuleVerdict;
    localizationKey: string;
    testId: string;
}

export interface RuleTypeOption {
    id: RuleType;
    localizationKey: string;
}
