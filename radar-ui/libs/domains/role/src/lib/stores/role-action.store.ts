import { createAction, props } from '@ngrx/store';

import { LoadStatus } from '@cs/core';

import { Role } from '../interfaces';

export const SET_ALL_ROLE_ENTITIES_DOC_ACTION = createAction('[Role] (Doc) Set All', props<{ list: Role[] }>());

export const UPDATE_ROLE_LOAD_STATUS_DOC_ACTION = createAction(
    '[Role] (Doc) Update Status',
    props<{ loadStatus: LoadStatus }>()
);

export const ROLE_LOAD_DONE_EVENT_ACTION = createAction(
    '[Role] {Event} Load Done',
    props<{ loadStatus: LoadStatus }>()
);
