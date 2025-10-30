import { RuleSeverity } from '@cs/domains/rule';
import { RuntimeEventProcessorHistoryControl } from '@cs/domains/runtime';

export interface RuntimeSettingPolicyForm {
    isEnabled: boolean;
    name: string;
    description: string;
    yaml: string;
}

export type RuntimeSettingPolicyRecord = {
    [key: string]: RuntimeSettingPolicyForm;
};

export interface RuntimeSettingPermissionForm {
    isAllowedType: boolean;
    namespaces: string[];
    pods: string[];
    labels: string[];
}

export type RuntimeSettingPermissionRecord = {
    [key: string]: RuntimeSettingPermissionForm;
};

export interface RuntimeSettingForm {
    policies: RuntimeSettingPolicyRecord;
    permissions: RuntimeSettingPermissionRecord;
    historyControl: RuntimeEventProcessorHistoryControl;
}

export interface RuntimeSettingPermissionCounter {
    allow: number;
    deny: number;
}

export interface RuntimeRuleForm {
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

export interface RuntimeExpertModeForm {
    isExpert: boolean;
}
