import { Update } from '@ngrx/entity';
import { createAction, props } from '@ngrx/store';

import { RuleState } from '../interfaces/state/rule-state.interface';
import { CreateRuleRequest, Rule, UpdateRuleRequest } from '../interfaces';

export const LOAD_RULE_ENTITIES_TODO_ACTION = createAction('[Rule] Load');

export const POLLING_LOAD_RULE_ENTITIES_TODO_ACTION = createAction('[Rule] Polling Load');

export const CREATE_RULE_ENTITY_TODO_ACTION = createAction('[Rule] Create', props<{ item: CreateRuleRequest }>());

export const UPDATE_RULE_ENTITY_TODO_ACTION = createAction(
    '[Rule] Update',
    props<{ id: string; item: UpdateRuleRequest }>()
);

export const DELETE_RULE_ENTITY_TODO_ACTION = createAction('[Rule] Delete', props<{ id: string }>());

export const UPDATE_RULE_STATE_DOC_ACTION = createAction(
    '[Rule] (Doc) Update State',
    props<Partial<Omit<RuleState, 'list'>>>()
);

export const SET_ALL_RULE_ENTITIES_DOC_ACTION = createAction('[Rule] (Doc) Set All', props<{ list: Rule[] }>());

export const SET_RULE_ENTITY_DOC_ACTION = createAction('[Rule] (Doc) Set One', props<{ item: Rule }>());

export const UPDATE_RULE_ENTITY_DOC_ACTION = createAction('[Rule] (Doc) Update', props<{ item: Update<Rule> }>());

export const DELETE_RULE_ENTITY_DOC_ACTION = createAction('[Rule] (Doc) Delete', props<{ id: string }>());
