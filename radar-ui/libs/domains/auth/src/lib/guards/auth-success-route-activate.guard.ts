import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';
import { inject } from '@angular/core';
import { Router, UrlTree } from '@angular/router';
import { filter, map } from 'rxjs/operators';

import { LoadStatus, RouterName } from '@cs/core';

import { AuthState } from '../interfaces/state/auth-state.interface';
import { getAuthLoadStatus } from '../stores/auth-selector.store';

const authSuccessRouteActivate = (): Observable<boolean | UrlTree> => {
    const store = inject<Store<AuthState>>(Store);
    const router = inject(Router);

    return store.select(getAuthLoadStatus).pipe(
        filter((status) => [LoadStatus.LOADED, LoadStatus.ERROR].includes(status)),
        map((status) => status === LoadStatus.LOADED),
        map((isAuthorized) => isAuthorized || router.createUrlTree([RouterName.SIGN_IN]))
    );
};

export const authSuccessRouteActivateGuard = () => authSuccessRouteActivate();

export const authSuccessRouteActivateChildGuard = () => authSuccessRouteActivate();
