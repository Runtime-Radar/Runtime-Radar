import { createAction, props } from '@ngrx/store';

import { CreateTokenRequest, Token, TokenState } from '../interfaces';

export const LOAD_TOKEN_ENTITIES_TODO_ACTION = createAction('[Token] Load');

export const POLLING_LOAD_TOKEN_ENTITIES_TODO_ACTION = createAction('[Token] Polling Load');

export const CREATE_TOKEN_ENTITY_TODO_ACTION = createAction('[Token] Create', props<{ item: CreateTokenRequest }>());

export const DELETE_TOKEN_ENTITY_TODO_ACTION = createAction('[Token] Delete', props<{ id: string }>());

export const REVOKE_TOKEN_ENTITIES_TODO_ACTION = createAction('[Token] Revoke');

export const UPDATE_TOKEN_STATE_DOC_ACTION = createAction(
    '[Token] (Doc) Update State',
    props<Partial<Omit<TokenState, 'list'>>>()
);

export const SET_ALL_TOKEN_ENTITIES_DOC_ACTION = createAction('[Token] (Doc) Set All', props<{ list: Token[] }>());

export const DELETE_ALL_TOKEN_ENTITIES_DOC_ACTION = createAction('[Token] (Doc) Delete All');

export const SET_TOKEN_ENTITY_DOC_ACTION = createAction('[Token] (Doc) Set One', props<{ item: Token }>());

export const UPSERT_TOKEN_ENTITIES_DOC_ACTION = createAction('[Token] (Doc) Upsert', props<{ list: Token[] }>());

export const DELETE_TOKEN_ENTITY_DOC_ACTION = createAction('[Token] (Doc) Delete', props<{ id: string }>());
