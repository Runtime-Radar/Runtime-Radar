import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { notificationLazyActivateGuard } from '@cs/domains/notification';
import { ruleLazyActivateGuard } from '@cs/domains/rule';
import { PermissionName, rolePermissionsResolver } from '@cs/domains/role';
import { detectorRuntimeActivateGuard, detectorRuntimeLazyActivateGuard } from '@cs/domains/detector';
import { runtimeActivateGuard, runtimeConfigModifyDeactivateGuard, runtimeDeactivateGuard } from '@cs/domains/runtime';

import { RuntimeFeatureDetailsContainer } from './containers/details/runtime-details.container';
import { RuntimeFeatureDetectorsContainer } from './containers/detectors/runtime-detectors.container';
import { RuntimeFeatureEventsContainer } from './containers/events/runtime-events.container';
import { RuntimeFeatureSettingsContainer } from './containers/settings/runtime-settings.container';
import { RuntimeRouterName } from './interfaces/runtime-navigation.interface';

const routes: Routes = [
    {
        path: '',
        canActivate: [runtimeActivateGuard],
        canDeactivate: [runtimeDeactivateGuard],
        children: [
            {
                path: RuntimeRouterName.SETTINGS,
                component: RuntimeFeatureSettingsContainer,
                canActivate: [ruleLazyActivateGuard, notificationLazyActivateGuard, detectorRuntimeLazyActivateGuard],
                canDeactivate: [runtimeConfigModifyDeactivateGuard],
                resolve: {
                    permissions: rolePermissionsResolver
                },
                data: {
                    localizationTitleKey: 'Runtime.SettingsPage.Header.Title',
                    permissions: [PermissionName.SYSTEM, PermissionName.RULES]
                }
            },
            {
                path: RuntimeRouterName.EVENTS,
                component: RuntimeFeatureEventsContainer,
                canActivate: [ruleLazyActivateGuard, notificationLazyActivateGuard, detectorRuntimeLazyActivateGuard],
                resolve: {
                    permissions: rolePermissionsResolver
                },
                data: {
                    localizationTitleKey: 'Runtime.EventsPage.Header.Title',
                    permissions: [PermissionName.RULES]
                }
            },
            {
                path: `${RuntimeRouterName.EVENTS}/:eventId`,
                component: RuntimeFeatureDetailsContainer,
                canActivate: [ruleLazyActivateGuard, notificationLazyActivateGuard, detectorRuntimeLazyActivateGuard],
                resolve: {
                    permissions: rolePermissionsResolver
                },
                data: {
                    localizationTitleKey: 'Runtime.DetailsPage.Header.Title',
                    permissions: [PermissionName.RULES]
                }
            },
            {
                path: RuntimeRouterName.DETECTORS,
                component: RuntimeFeatureDetectorsContainer,
                canActivate: [detectorRuntimeActivateGuard],
                resolve: {
                    permissions: rolePermissionsResolver
                },
                data: {
                    localizationTitleKey: 'Runtime.DetectorsPage.Header.Title',
                    permissions: [PermissionName.SYSTEM]
                }
            },
            {
                path: '**',
                redirectTo: RuntimeRouterName.EVENTS
            }
        ]
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class RuntimeFeatureRoutingModule {}
