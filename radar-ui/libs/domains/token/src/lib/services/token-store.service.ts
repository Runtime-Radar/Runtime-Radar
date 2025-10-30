import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';

import { LoadStatus } from '@cs/core';

import {
    CREATE_TOKEN_ENTITY_TODO_ACTION,
    DELETE_TOKEN_ENTITY_TODO_ACTION,
    REVOKE_TOKEN_ENTITIES_TODO_ACTION
} from '../stores/token-action.store';
import { CreateTokenRequest, Token, TokenState } from '../interfaces';
import { getTokenLoadStatus, getTokens } from '../stores/token-selector.store';

@Injectable({
    providedIn: 'root'
})
export class TokenStoreService {
    readonly tokens$: Observable<Token[]> = this.store.select(getTokens);

    readonly loadStatus$: Observable<LoadStatus> = this.store.select(getTokenLoadStatus);

    constructor(private readonly store: Store<TokenState>) {}

    createToken(item: CreateTokenRequest) {
        this.store.dispatch(CREATE_TOKEN_ENTITY_TODO_ACTION({ item }));
    }

    deleteToken(id: string) {
        this.store.dispatch(DELETE_TOKEN_ENTITY_TODO_ACTION({ id }));
    }

    revokeTokens() {
        this.store.dispatch(REVOKE_TOKEN_ENTITIES_TODO_ACTION());
    }
}
