import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';
import { inject } from '@angular/core';
import { Router, UrlTree } from '@angular/router';
import { filter, map } from 'rxjs/operators';

import { DEFAULT_ROUTER_NAME, LoadStatus } from '@cs/core';

import { AuthState } from '../interfaces/state/auth-state.interface';
import { getAuthLoadStatus } from '../stores/auth-selector.store';

const authErrorRouteActivate = (): Observable<boolean | UrlTree> => {
    const store = inject<Store<AuthState>>(Store);
    const router = inject(Router);

    return store.select(getAuthLoadStatus).pipe(
        filter((status) => [LoadStatus.LOADED, LoadStatus.ERROR].includes(status)),
        map((status) => status === LoadStatus.ERROR),
        map((isNotAuthorized) => isNotAuthorized || router.createUrlTree([DEFAULT_ROUTER_NAME]))
    );
};

export const authErrorRouteActivateGuard = () => authErrorRouteActivate();

export const authErrorRouteActivateChildGuard = () => authErrorRouteActivate();
