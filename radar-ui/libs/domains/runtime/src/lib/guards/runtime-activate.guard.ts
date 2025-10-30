import { Store } from '@ngrx/store';
import { inject } from '@angular/core';
import { Observable, filter, map, tap } from 'rxjs';
import { Router, UrlTree } from '@angular/router';

import { LoadStatus, RouterName } from '@cs/core';

import { LOAD_RUNTIME_CONFIG_TODO_ACTION } from '../stores/runtime-action.store';
import { RuntimeState } from '../interfaces';
import { getRuntimeLoadStatus } from '../stores/runtime-selector.store';

const runtimeActivate = (): Observable<boolean | UrlTree> => {
    const router = inject(Router);
    const store = inject<Store<RuntimeState>>(Store);

    return store.select(getRuntimeLoadStatus).pipe(
        tap((status) => {
            if (status === LoadStatus.INIT) {
                store.dispatch(LOAD_RUNTIME_CONFIG_TODO_ACTION());
            }
        }),
        filter((status) => [LoadStatus.LOADED, LoadStatus.ERROR].includes(status)),
        map((status) => status === LoadStatus.LOADED),
        map((isLoaded) => isLoaded || router.createUrlTree([RouterName.ERROR]))
    );
};

export const runtimeActivateGuard = () => runtimeActivate();

export const runtimeActivateChildGuard = () => runtimeActivate();
