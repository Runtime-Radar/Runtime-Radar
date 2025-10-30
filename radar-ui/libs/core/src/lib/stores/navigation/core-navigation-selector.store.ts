import * as fromRouter from '@ngrx/router-store';
import { Data } from '@angular/router';
import { ActionReducerMap, createFeatureSelector, createSelector } from '@ngrx/store';

import { RouterName } from '../../constants/core-router.constant';

export const CORE_NAVIGATION_KEY = 'router';

export interface CoreNavigationState {
    readonly router: fromRouter.RouterReducerState;
}

const selectCoreNavigationState = createFeatureSelector<fromRouter.RouterReducerState>(CORE_NAVIGATION_KEY);
const { selectRouteData, selectUrl, selectCurrentRoute } = fromRouter.getRouterSelectors(selectCoreNavigationState);
const selectCoreNavigationRouteData = createSelector(selectRouteData, (data: Data) => data);

export const getCurrentRouteSlug = createSelector(
    selectCurrentRoute,
    (route) => (route?.routeConfig?.path as string) || ''
);

export const getCurrentRouterName = createSelector(selectUrl, (url) => {
    const paths = url?.split('/') || [];

    return ((paths.at(1) === RouterName.SETTINGS ? paths.at(2) : paths.at(1)) as RouterName) || RouterName.DEFAULT;
});

export const getLocalizationTitleKey = createSelector(
    selectCoreNavigationRouteData,
    (data) => data['localizationTitleKey'] ?? ''
);

export const coreNavigationReducer: ActionReducerMap<CoreNavigationState> = {
    router: fromRouter.routerReducer
};
