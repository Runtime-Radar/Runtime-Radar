import { Rule, RuleSeverity } from '@cs/domains/rule';

export interface SharedRuleSidepanelFormProps {
    rule: Partial<Rule>;
    isEdit: boolean;
}

export interface RuleForm {
    name: string;
    namespaces: string[];
    notifySeverity: RuleSeverity;
    mailIds: string[];
    detectors: string[];
    pods: string[];
    containers: string[];
    nodes: string[];
    binaries: string[];
    imageNames: string[];
    registries: string[];
}
