import { Injectable } from '@angular/core';
import { Action, Store } from '@ngrx/store';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { KbqToastService, KbqToastStyle } from '@koobiq/components/toast';
import { Observable, combineLatest, forkJoin, of } from 'rxjs';
import { catchError, filter, map, switchMap, take, tap } from 'rxjs/operators';

import { I18nService } from '@cs/i18n';
import { SIGN_OUT_EVENT_ACTION } from '@cs/domains/auth';
import { SWITCH_CLUSTER_EVENT_ACTION } from '@cs/domains/cluster';
import { CoreWindowService, LoadStatus, CoreUtilsService as utils } from '@cs/core';

import { RuntimeRequestService } from '../services/runtime-request.service';
import { RuntimeHelperService as runtimeHelper } from '../services/runtime-helper.service';
import {
    CHECK_RUNTIME_CHANGES_TODO_ACTION,
    CREATE_RUNTIME_CONFIG_TODO_ACTION,
    DEACTIVATE_RUNTIME_CONFIG_TODO_ACTION,
    HIDE_RUNTIME_OVERLAY_TODO_ACTION,
    LOAD_RUNTIME_CONFIG_TODO_ACTION,
    SWITCH_RUNTIME_EXPERT_MODE_TODO_ACTION,
    UPDATE_RUNTIME_STATE_DOC_ACTION
} from './runtime-action.store';
import { RuntimeConfigStatus, RuntimeState } from '../interfaces';
import {
    getRuntimeEventProcessorHistoryControl,
    getRuntimeIsExpertMode,
    getRuntimeLoadStatus,
    getRuntimeMonitorConfig
} from './runtime-selector.store';

const RUNTIME_EXPERT_MODE_LOCAL_STORAGE_KEY = 'xprtmd';

@Injectable({
    providedIn: 'root'
})
export class RuntimeEffectStore {
    readonly checkExpertMode$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(LOAD_RUNTIME_CONFIG_TODO_ACTION),
            map(() => this.coreWindowService.localStorage.getItem(RUNTIME_EXPERT_MODE_LOCAL_STORAGE_KEY)),
            map((value) =>
                UPDATE_RUNTIME_STATE_DOC_ACTION({
                    isExpertMode: value ? value === 'true' : false
                })
            )
        )
    );

    readonly loadConfig$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(LOAD_RUNTIME_CONFIG_TODO_ACTION),
            switchMap(() =>
                forkJoin({
                    config: this.runtimeRequestService.getRuntimeMonitor().pipe(
                        map((response) => response.config),
                        catchError(() => of(undefined))
                    ),
                    historyControl: this.runtimeRequestService.getEventProcessor().pipe(
                        map((response) => response.config.history_control),
                        catchError(() => of(undefined))
                    )
                })
            ),
            map(({ config, historyControl }) => {
                if (config === undefined) {
                    return UPDATE_RUNTIME_STATE_DOC_ACTION({
                        loadStatus: LoadStatus.ERROR
                    });
                }

                return UPDATE_RUNTIME_STATE_DOC_ACTION({
                    loadStatus: LoadStatus.LOADED,
                    historyControl,
                    config
                });
            })
        )
    );

    readonly reloadConfig$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(SWITCH_CLUSTER_EVENT_ACTION),
            switchMap(() => this.store.select(getRuntimeLoadStatus).pipe(take(1))),
            filter((status) => status !== LoadStatus.INIT),
            map(() => LOAD_RUNTIME_CONFIG_TODO_ACTION())
        )
    );

    readonly deactivateState$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(DEACTIVATE_RUNTIME_CONFIG_TODO_ACTION),
            map(() =>
                UPDATE_RUNTIME_STATE_DOC_ACTION({
                    loadStatus: LoadStatus.INIT // @todo: set separate status e.g. PARTIAL
                })
            )
        )
    );

    readonly createRuntimeMonitor$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(CREATE_RUNTIME_CONFIG_TODO_ACTION),
            switchMap((action) =>
                this.runtimeRequestService.createRuntimeMonitor(action.config).pipe(
                    take(1),
                    map((response) => response.config)
                )
            ),
            filter((config) => config && !!Object.keys(config).length),
            map((config) =>
                UPDATE_RUNTIME_STATE_DOC_ACTION({
                    config,
                    configStatus: RuntimeConfigStatus.STAY
                })
            ),
            tap(() => {
                this.toastService.show({
                    style: KbqToastStyle.Success,
                    title: this.i18nService.translate('Runtime.Pseudo.Notification.Created')
                });
            })
        )
    );

    readonly createEventProcessor$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(CREATE_RUNTIME_CONFIG_TODO_ACTION),
            switchMap((action) =>
                this.runtimeRequestService.createEventProcessor(action.historyControl).pipe(
                    take(1),
                    map((response) => response.config.history_control)
                )
            ),
            filter((historyControl) => !!historyControl),
            map((historyControl) => UPDATE_RUNTIME_STATE_DOC_ACTION({ historyControl }))
        )
    );

    readonly checkChanges$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(CHECK_RUNTIME_CHANGES_TODO_ACTION),
            switchMap((action) =>
                combineLatest([
                    this.store.select(getRuntimeEventProcessorHistoryControl),
                    this.store.select(getRuntimeMonitorConfig)
                ]).pipe(
                    take(1),
                    map(([historyControl, config]) => ({
                        previous: runtimeHelper.convertConfigToDiffValues(config, historyControl),
                        current: action.config
                    }))
                )
            ),
            map(({ previous, current }) => {
                const hasChanges = !utils.isEqual(previous, current);

                return UPDATE_RUNTIME_STATE_DOC_ACTION({
                    hasChanges,
                    hasPoliciesChanges: !utils.isEqual(
                        Object.entries(previous.tracing_policies).map(([key, { enabled, ...rest }]) => [key, rest]),
                        Object.entries(current.tracing_policies).map(([key, { enabled, ...rest }]) => [key, rest])
                    ),
                    configStatus: hasChanges ? RuntimeConfigStatus.MODIFY : RuntimeConfigStatus.STAY
                });
            })
        )
    );

    readonly switchExpertMode$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(SWITCH_RUNTIME_EXPERT_MODE_TODO_ACTION),
            switchMap(() => this.store.select(getRuntimeIsExpertMode).pipe(take(1))),
            tap((isExpertMode) => {
                this.coreWindowService.localStorage.setItem(
                    RUNTIME_EXPERT_MODE_LOCAL_STORAGE_KEY,
                    (!isExpertMode).toString()
                );
            }),
            map((isExpertMode) => UPDATE_RUNTIME_STATE_DOC_ACTION({ isExpertMode: !isExpertMode }))
        )
    );

    readonly deactivateExpertMode$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(SIGN_OUT_EVENT_ACTION),
            map(() => UPDATE_RUNTIME_STATE_DOC_ACTION({ isExpertMode: false })),
            tap(() => {
                this.coreWindowService.localStorage.removeItem(RUNTIME_EXPERT_MODE_LOCAL_STORAGE_KEY);
            })
        )
    );

    readonly hideOverlay$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(HIDE_RUNTIME_OVERLAY_TODO_ACTION),
            map(() =>
                UPDATE_RUNTIME_STATE_DOC_ACTION({
                    isOverlayed: false
                })
            )
        )
    );

    constructor(
        private readonly actions$: Actions,
        private readonly i18nService: I18nService,
        private readonly coreWindowService: CoreWindowService,
        private readonly runtimeRequestService: RuntimeRequestService,
        private readonly toastService: KbqToastService,
        private readonly store: Store<RuntimeState>
    ) {}
}
