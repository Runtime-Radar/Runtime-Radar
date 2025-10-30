import { Update } from '@ngrx/entity';
import { createAction, props } from '@ngrx/store';

import { GetTokenResponse } from '@cs/domains/auth';

import { UserState } from '../interfaces/state/user-state.interface';
import { CreateUserRequest, User, UserEditRequest } from '../interfaces';

export const LOAD_USER_ENTITIES_TODO_ACTION = createAction('[User] Load');

export const POLLING_LOAD_USER_ENTITIES_TODO_ACTION = createAction('[User] Polling Load');

export const CREATE_USER_ENTITY_TODO_ACTION = createAction('[User] Create', props<{ item: CreateUserRequest }>());

export const UPDATE_USER_ENTITY_TODO_ACTION = createAction(
    '[User] Update',
    props<{ id: string; item: UserEditRequest }>()
);

export const UPDATE_USER_PASSWORD_TODO_ACTION = createAction('[User] Update Password', props<{ password: string }>());

export const DELETE_USER_ENTITY_TODO_ACTION = createAction('[User] Delete', props<{ id: string }>());

export const UPDATE_USER_STATE_DOC_ACTION = createAction(
    '[User] (Doc) Update State',
    props<Partial<Omit<UserState, 'list'>>>()
);

export const SET_ALL_USER_ENTITIES_DOC_ACTION = createAction('[User] (Doc) Set All', props<{ list: User[] }>());

export const SET_USER_ENTITY_DOC_ACTION = createAction('[User] (Doc) Set One', props<{ item: User }>());

export const UPDATE_USER_ENTITY_DOC_ACTION = createAction('[User] (Doc) Update', props<{ item: Update<User> }>());

export const DELETE_USER_ENTITY_DOC_ACTION = createAction('[User] (Doc) Delete', props<{ id: string }>());

export const DELETE_USER_ENTITIES_DOC_ACTION = createAction('[User] (Doc) Delete Many', props<{ ids: string[] }>());

export const UPDATE_USER_PASSWORD_EVENT_ACTION = createAction(
    '[User] {Event} Update Password',
    props<{ tokens: GetTokenResponse }>()
);
