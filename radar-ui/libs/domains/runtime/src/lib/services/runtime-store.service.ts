import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';

import {
    CHECK_RUNTIME_CHANGES_TODO_ACTION,
    CREATE_RUNTIME_CONFIG_TODO_ACTION,
    HIDE_RUNTIME_OVERLAY_TODO_ACTION,
    SWITCH_RUNTIME_EXPERT_MODE_TODO_ACTION
} from '../stores/runtime-action.store';
import {
    RuntimeEventProcessorHistoryControl,
    RuntimeMonitorConfig,
    RuntimeMonitorConfigExtended,
    RuntimeState
} from '../interfaces';
import {
    getRuntimeEventProcessorHistoryControl,
    getRuntimeHasChanges,
    getRuntimeHasPoliciesChanges,
    getRuntimeIsExpertMode,
    getRuntimeIsOverlayed,
    getRuntimeMonitorConfig
} from '../stores/runtime-selector.store';

@Injectable({
    providedIn: 'root'
})
export class RuntimeStoreService {
    readonly runtimeMonitorConfig$: Observable<RuntimeMonitorConfig> = this.store.select(getRuntimeMonitorConfig);

    readonly eventProcessorHistoryControl$: Observable<RuntimeEventProcessorHistoryControl | undefined> =
        this.store.select(getRuntimeEventProcessorHistoryControl);

    readonly runtimeHasChanges$: Observable<boolean> = this.store.select(getRuntimeHasChanges);

    readonly runtimeHasPoliciesChanges$: Observable<boolean> = this.store.select(getRuntimeHasPoliciesChanges);

    readonly runtimeIsExpertMode$: Observable<boolean> = this.store.select(getRuntimeIsExpertMode);

    readonly runtimeIsOverlayed$: Observable<boolean> = this.store.select(getRuntimeIsOverlayed);

    constructor(private readonly store: Store<RuntimeState>) {}

    createConfig(config: RuntimeMonitorConfig, historyControl: RuntimeEventProcessorHistoryControl) {
        this.store.dispatch(
            CREATE_RUNTIME_CONFIG_TODO_ACTION({
                config,
                historyControl
            })
        );
    }

    checkChanges(config: RuntimeMonitorConfigExtended) {
        this.store.dispatch(CHECK_RUNTIME_CHANGES_TODO_ACTION({ config }));
    }

    switchExpertMode() {
        this.store.dispatch(SWITCH_RUNTIME_EXPERT_MODE_TODO_ACTION());
    }

    hideOverlay() {
        this.store.dispatch(HIDE_RUNTIME_OVERLAY_TODO_ACTION());
    }
}
