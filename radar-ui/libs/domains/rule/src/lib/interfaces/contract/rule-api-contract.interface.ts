import { Rule, RuleEntity, RuleScope, RuleType } from './rule-contract.interface';

export interface GetRulesResponse {
    total: number;
    rules: Rule[];
}

export interface GetRuleResponse {
    rule: Rule;
    deleted: boolean;
}

export interface CreateRuleRequest {
    name: string;
    type: RuleType;
    rule: RuleEntity;
    scope: RuleScope;
}

export interface CreateRuleResponse {
    id: string;
}

export type UpdateRuleRequest = CreateRuleRequest;

export type EmptyRuleResponse = Record<string, unknown>;
