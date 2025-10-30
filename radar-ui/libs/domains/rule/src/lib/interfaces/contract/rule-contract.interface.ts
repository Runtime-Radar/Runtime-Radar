import { OneOf } from '@cs/core';

export enum RuleType {
    TYPE_RUNTIME = 'TYPE_RUNTIME'
}

export enum RuleSeverity {
    NONE = 'none',
    LOW = 'low',
    MEDIUM = 'medium',
    HIGH = 'high',
    CRITICAL = 'critical'
}

export enum RuleVerdict {
    NONE = 'none',
    CLEAN = 'clean',
    UNWANTED = 'unwanted',
    DANGEROUS = 'dangerous'
}

// fields are described into RuleFeatureHelperService#19
export interface RuleWhiteList {
    threats: string[];
    binaries: string[];
}

export type RuleNotifyEntity = {
    targets: string[]; // targetIds - ids for notification
} & OneOf<{
    severity: RuleSeverity;
    verdict: RuleVerdict;
}>;

export interface RuleEntity {
    version: string;
    notify?: RuleNotifyEntity | null;
    whitelist: RuleWhiteList;
}

export interface RuleScope {
    version: string;
    image_names: string[];
    registries: string[];
    clusters: string[];
    namespaces: string[];
    pods: string[];
    containers: string[];
    nodes: string[];
}

export interface Rule {
    id: string;
    type: RuleType;
    name: string;
    rule: RuleEntity;
    scope: RuleScope;
}
