import { Store } from '@ngrx/store';
import { inject } from '@angular/core';
import { ActivatedRouteSnapshot, ResolveFn } from '@angular/router';
import { distinctUntilChanged, map, switchMap } from 'rxjs';

import { AuthStoreService } from '@cs/domains/auth';

import { RoleState } from '../interfaces/state/role-state.interface';
import { getRole } from '../stores/role-selector.store';
import { PermissionName, PermissionType, Role, RolePermissionMap } from '../interfaces';

export const rolePermissionsResolver: ResolveFn<RolePermissionMap> = (route: ActivatedRouteSnapshot) => {
    const authStoreService = inject(AuthStoreService);
    const store = inject<Store<RoleState>>(Store);

    return authStoreService.credentials$.pipe(
        map((credentials) => credentials.roleId),
        distinctUntilChanged(),
        switchMap((roleId) => store.select(getRole(roleId))),
        map((role: Role | undefined) => {
            /* eslint @typescript-eslint/dot-notation: "off" */
            const permissions: PermissionName[] | undefined = route.data['permissions'];
            if (!role || !permissions) {
                console.warn('roleId must be provided');

                return {} as RolePermissionMap;
            }

            return permissions.reduce((acc, key) => {
                const actions = role.role_permissions[key].actions;
                acc[key] = actions.reduce((m, type) => m.set(type, true), new Map<PermissionType, boolean>());

                return acc;
            }, {} as RolePermissionMap);
        })
    );
};
