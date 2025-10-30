import { Store } from '@ngrx/store';
import { inject } from '@angular/core';
import { Observable, filter, map, tap } from 'rxjs';

import { LoadStatus } from '@cs/core';

import { DEACTIVATE_RUNTIME_CONFIG_TODO_ACTION } from '../stores/runtime-action.store';
import { RuntimeState } from '../interfaces';
import { getRuntimeLoadStatus } from '../stores/runtime-selector.store';

const runtimeDeactivate = (): Observable<boolean> => {
    const store = inject<Store<RuntimeState>>(Store);

    return store.select(getRuntimeLoadStatus).pipe(
        tap((status) => {
            if ([LoadStatus.LOADED, LoadStatus.ERROR].includes(status)) {
                store.dispatch(DEACTIVATE_RUNTIME_CONFIG_TODO_ACTION());
            }
        }),
        filter((status) => status === LoadStatus.INIT),
        map((status) => status === LoadStatus.INIT)
    );
};

export const runtimeDeactivateGuard = () => runtimeDeactivate();
