import { Store } from '@ngrx/store';
import { inject } from '@angular/core';
import { Observable, filter, map, tap } from 'rxjs';

import { UPDATE_RUNTIME_STATE_DOC_ACTION } from '../stores/runtime-action.store';
import { getRuntimeConfigStatus } from '../stores/runtime-selector.store';
import { RuntimeConfigStatus, RuntimeState } from '../interfaces';

const runtimeConfigModifyDeactivate = (): Observable<boolean> => {
    const store = inject<Store<RuntimeState>>(Store);

    return store.select(getRuntimeConfigStatus).pipe(
        tap((status) => {
            if (status === RuntimeConfigStatus.MODIFY) {
                store.dispatch(UPDATE_RUNTIME_STATE_DOC_ACTION({ isOverlayed: true }));
            }
        }),
        map((status) => status === RuntimeConfigStatus.STAY),
        filter((isNavigateAllowed) => isNavigateAllowed),
        tap(() => {
            store.dispatch(
                UPDATE_RUNTIME_STATE_DOC_ACTION({ configStatus: RuntimeConfigStatus.INIT, isOverlayed: false })
            );
        })
    );
};

export const runtimeConfigModifyDeactivateGuard = () => runtimeConfigModifyDeactivate();
