import { Injectable } from '@angular/core';
import { Store } from '@ngrx/store';
import { BehaviorSubject, Observable } from 'rxjs';

import { LoadStatus, RouterName } from '@cs/core';

import {
    CoreNavigationState,
    getCurrentRouteSlug,
    getCurrentRouterName
} from '../stores/navigation/core-navigation-selector.store';

@Injectable({
    providedIn: 'root'
})
export class CoreNavigationStoreService {
    readonly loadStatus$ = new BehaviorSubject<LoadStatus>(LoadStatus.INIT);

    readonly routeSlug$: Observable<string> = this.store.select(getCurrentRouteSlug);

    readonly routerName$: Observable<RouterName> = this.store.select(getCurrentRouterName);

    constructor(private readonly store: Store<CoreNavigationState>) {}

    setLoadStatus(status: LoadStatus) {
        this.loadStatus$.next(status);
    }
}
