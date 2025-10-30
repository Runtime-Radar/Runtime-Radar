import { createAction, props } from '@ngrx/store';

import {
    RuntimeEventProcessorHistoryControl,
    RuntimeMonitorConfig,
    RuntimeMonitorConfigExtended,
    RuntimeState
} from '../interfaces';

export const LOAD_RUNTIME_CONFIG_TODO_ACTION = createAction('[Runtime] Load Config');

export const DEACTIVATE_RUNTIME_CONFIG_TODO_ACTION = createAction('[Runtime] Deactivate Config');

export const CREATE_RUNTIME_CONFIG_TODO_ACTION = createAction(
    '[Runtime] Create Config',
    props<{ config: RuntimeMonitorConfig; historyControl: RuntimeEventProcessorHistoryControl }>()
);

export const CHECK_RUNTIME_CHANGES_TODO_ACTION = createAction(
    '[Runtime] Check Changes',
    props<{ config: RuntimeMonitorConfigExtended }>()
);

export const SWITCH_RUNTIME_EXPERT_MODE_TODO_ACTION = createAction('[Runtime] Switch Expert Mode');

export const HIDE_RUNTIME_OVERLAY_TODO_ACTION = createAction('[Runtime] Hide Overlay');

export const UPDATE_RUNTIME_STATE_DOC_ACTION = createAction(
    '[Runtime] (Doc) Update State',
    props<Partial<RuntimeState>>()
);
