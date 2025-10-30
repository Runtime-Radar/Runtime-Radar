import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';

import { UserState } from '../interfaces/state/user-state.interface';
import { getUsers } from '../stores/user-selector.store';
import {
    CREATE_USER_ENTITY_TODO_ACTION,
    DELETE_USER_ENTITY_TODO_ACTION,
    UPDATE_USER_ENTITY_TODO_ACTION,
    UPDATE_USER_PASSWORD_TODO_ACTION
} from '../stores/user-action.store';
import { CreateUserRequest, User, UserEditRequest } from '../interfaces';

@Injectable({
    providedIn: 'root'
})
export class UserStoreService {
    readonly users$: Observable<User[]> = this.store.select(getUsers);

    constructor(private readonly store: Store<UserState>) {}

    createUser(item: CreateUserRequest) {
        this.store.dispatch(CREATE_USER_ENTITY_TODO_ACTION({ item }));
    }

    updateUser(id: string, item: UserEditRequest) {
        this.store.dispatch(UPDATE_USER_ENTITY_TODO_ACTION({ id, item }));
    }

    deleteUser(id: string) {
        this.store.dispatch(DELETE_USER_ENTITY_TODO_ACTION({ id }));
    }

    changePassword(password: string) {
        this.store.dispatch(UPDATE_USER_PASSWORD_TODO_ACTION({ password }));
    }
}
