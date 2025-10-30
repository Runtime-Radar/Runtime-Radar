import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';

import { Role } from '../interfaces';
import { RoleState } from '../interfaces/state/role-state.interface';
import { getRole, getRoles } from '../stores/role-selector.store';

@Injectable({
    providedIn: 'root'
})
export class RoleStoreService {
    readonly roles$: Observable<Role[]> = this.store.select(getRoles);

    readonly role$ = (id: string): Observable<Role | undefined> => this.store.select(getRole(id));

    constructor(private readonly store: Store<RoleState>) {}
}
