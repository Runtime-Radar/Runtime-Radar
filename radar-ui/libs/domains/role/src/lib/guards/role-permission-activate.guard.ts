import { Store } from '@ngrx/store';
import { inject } from '@angular/core';
import { ActivatedRouteSnapshot, Router, UrlTree } from '@angular/router';
import { Observable, filter, map, switchMap, tap } from 'rxjs';

import { AuthStoreService } from '@cs/domains/auth';
import { LoadStatus, RouterName } from '@cs/core';

import { RoleState } from '../interfaces/state/role-state.interface';
import { PermissionType, Role } from '../interfaces';
import { getRole, getRoleLoadStatus } from '../stores/role-selector.store';

const rolePermissionActivate = (route: ActivatedRouteSnapshot): Observable<boolean | UrlTree> => {
    const authStoreService = inject<AuthStoreService>(AuthStoreService);
    const router = inject(Router);
    const store = inject<Store<RoleState>>(Store);

    return store.select(getRoleLoadStatus).pipe(
        tap((status) => {
            if (status === LoadStatus.ERROR) {
                console.warn('roles must be provided');
                authStoreService.destroyTokens();
            }
        }),
        filter((status) => status === LoadStatus.LOADED),
        switchMap(() => authStoreService.credentials$.pipe(map((credentials) => credentials.roleId))),
        switchMap((roleId) => store.select(getRole(roleId))),
        map((role: Role | undefined) => {
            /* eslint @typescript-eslint/dot-notation: "off" */
            const permissions: PermissionName[] | undefined = route.data['guards'];
            if (!role || !permissions) {
                console.warn('roleId must be provided');
                authStoreService.destroyTokens();

                return false;
            }

            return permissions.every((key) => {
                const prm = role.role_permissions[key.toString() as keyof typeof role.role_permissions];

                return prm.actions.includes(PermissionType.READ);
            });
        }),
        map((hasAccess) => hasAccess || router.createUrlTree([RouterName.FORBIDDEN]))
    );
};

export const rolePermissionActivateGuard = (route: ActivatedRouteSnapshot) => rolePermissionActivate(route);

export const rolePermissionActivateChildGuard = (route: ActivatedRouteSnapshot) => rolePermissionActivate(route);
