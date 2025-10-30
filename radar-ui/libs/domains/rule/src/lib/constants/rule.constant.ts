import {
    RuleSeverity,
    RuleSeverityOption,
    RuleSeverityOrder,
    RuleType,
    RuleTypeOption,
    RuleVerdict,
    RuleVerdictOption
} from '../interfaces';

export const SEVERITY_UNDEFINED_LOCALIZATION_KEY = 'Common.Pseudo.Severity.Undefined';

export const RULE_SEVERITY_ORDER: RuleSeverityOrder = {
    [RuleSeverity.NONE]: 4,
    [RuleSeverity.LOW]: 3,
    [RuleSeverity.MEDIUM]: 2,
    [RuleSeverity.HIGH]: 1,
    [RuleSeverity.CRITICAL]: 0
};

export const RULE_SEVERITIES: RuleSeverityOption[] = [
    {
        id: RuleSeverity.LOW,
        localizationKey: 'Common.Pseudo.Severity.Low',
        testId: 'low-radio'
    },
    {
        id: RuleSeverity.MEDIUM,
        localizationKey: 'Common.Pseudo.Severity.Medium',
        testId: 'medium-radio'
    },
    {
        id: RuleSeverity.HIGH,
        localizationKey: 'Common.Pseudo.Severity.High',
        testId: 'high-radio'
    },
    {
        id: RuleSeverity.CRITICAL,
        localizationKey: 'Common.Pseudo.Severity.Critical',
        testId: 'critical-radio'
    },
    {
        id: RuleSeverity.NONE,
        localizationKey: 'Common.Pseudo.Severity.None',
        testId: 'none-radio'
    }
];

export const RULE_VERDICTS: RuleVerdictOption[] = [
    {
        id: RuleVerdict.CLEAN,
        localizationKey: 'Common.Pseudo.Verdict.Clean',
        testId: 'clean-radio'
    },
    {
        id: RuleVerdict.UNWANTED,
        localizationKey: 'Common.Pseudo.Verdict.Unwanted',
        testId: 'unwanted-radio'
    },
    {
        id: RuleVerdict.DANGEROUS,
        localizationKey: 'Common.Pseudo.Verdict.Dangerous',
        testId: 'dangerous-radio'
    },
    {
        id: RuleVerdict.NONE,
        localizationKey: 'Common.Pseudo.Verdict.None',
        testId: 'none-radio'
    }
];

export const RULE_TYPE: RuleTypeOption[] = [
    {
        id: RuleType.TYPE_RUNTIME,
        localizationKey: 'Common.Pseudo.ScanType.Runtime'
    }
];
