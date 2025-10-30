import { createAction, props } from '@ngrx/store';

import { AuthCredentials, AuthState } from '../interfaces/state/auth-state.interface';

export const GET_LOCATION_PATH_TODO_ACTION = createAction('[Auth] Get Location Path', props<{ location: Location }>());

export const REDIRECT_TO_SWITCH_ROUTE_TODO_ACTION = createAction('[Auth] Redirect To Switch Route');

export const APPLY_AUTH_TOKENS_TODO_ACTION = createAction('[Auth] Apply Tokens');

export const EXPIRE_AUTH_TOKENS_TODO_ACTION = createAction('[Auth] Expire Tokens');

export const EXPIRE_PASSWORD_TODO_ACTION = createAction('[Auth] Password Expire');

export const SIGN_IN_TODO_ACTION = createAction('[Auth] Sign In', props<{ username: string; password: string }>());

export const SUCCESS_SIGN_IN_TODO_ACTION = createAction('[Auth] Success Sign In', props<AuthCredentials>());

export const SIGN_OUT_TODO_ACTION = createAction('[Auth] Sign Out');

export const RESET_AUTH_CREDENTIALS_DOC_ACTION = createAction('[Auth] (Doc) Reset Credentials');

export const UPDATE_AUTH_STATE_DOC_ACTION = createAction('[Auth] (Doc) Update State', props<Partial<AuthState>>());

export const ALLOW_AUTH_EVENT_ACTION = createAction('[Auth] {Event} Allow Auth');

export const SIGN_OUT_EVENT_ACTION = createAction('[Auth] {Event} Sign Out');
