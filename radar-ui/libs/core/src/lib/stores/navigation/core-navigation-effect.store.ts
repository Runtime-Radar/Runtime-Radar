import { Injectable } from '@angular/core';
import { Action, Store } from '@ngrx/store';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { NEVER, Observable } from 'rxjs';
import { ROUTER_CANCEL, ROUTER_ERROR, ROUTER_NAVIGATED, ROUTER_NAVIGATION, ROUTER_REQUEST } from '@ngrx/router-store';
import { concatMap, switchMap, tap } from 'rxjs/operators';

import { I18nService } from '@cs/i18n';

import { CoreMetaService } from '../../services/core-meta.service';
import { CoreNavigationStoreService } from '../../services/core-navigation-store.service';
import { LoadStatus } from '../../constants';
import { CoreNavigationState, getLocalizationTitleKey } from './core-navigation-selector.store';

@Injectable({
    providedIn: 'root'
})
export class CoreNavigationEffectStore {
    readonly updateTitle$: Observable<Action> = createEffect(
        () =>
            this.actions$.pipe(
                ofType(ROUTER_NAVIGATION),
                switchMap(() => this.store.select(getLocalizationTitleKey)),
                tap((localizationTitleKey: string) => {
                    this.coreMetaService.setPageMeta({
                        title: this.i18nService.translate(localizationTitleKey)
                    });
                }),
                concatMap(() => NEVER)
            ),
        { dispatch: false }
    );

    readonly startLoading$: Observable<Action> = createEffect(
        () =>
            this.actions$.pipe(
                ofType(ROUTER_REQUEST),
                tap(() => {
                    this.coreNavigationStoreService.setLoadStatus(LoadStatus.IN_PROGRESS);
                }),
                concatMap(() => NEVER)
            ),
        { dispatch: false }
    );

    readonly endLoading$: Observable<Action> = createEffect(
        () =>
            this.actions$.pipe(
                ofType(ROUTER_NAVIGATED),
                tap(() => {
                    this.coreNavigationStoreService.setLoadStatus(LoadStatus.LOADED);
                }),
                concatMap(() => NEVER)
            ),
        { dispatch: false }
    );

    readonly failLoading$: Observable<Action> = createEffect(
        () =>
            this.actions$.pipe(
                ofType(ROUTER_ERROR, ROUTER_CANCEL),
                tap(() => {
                    this.coreNavigationStoreService.setLoadStatus(LoadStatus.ERROR);
                }),
                concatMap(() => NEVER)
            ),
        { dispatch: false }
    );

    constructor(
        private readonly actions$: Actions,
        private readonly i18nService: I18nService,
        private readonly coreMetaService: CoreMetaService,
        private readonly coreNavigationStoreService: CoreNavigationStoreService,
        private readonly store: Store<CoreNavigationState>
    ) {}
}
