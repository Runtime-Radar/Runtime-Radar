import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { detectorRuntimeLazyActivateGuard } from '@cs/domains/detector';
import { integrationActivateGuard } from '@cs/domains/integration';
import { notificationActivateGuard } from '@cs/domains/notification';
import { ruleLazyActivateGuard } from '@cs/domains/rule';

import { IntegrationFeatureListContainer } from './containers/list/integration-list.container';

const routes: Routes = [
    {
        path: '',
        component: IntegrationFeatureListContainer,
        canActivate: [
            integrationActivateGuard,
            notificationActivateGuard,
            ruleLazyActivateGuard,
            detectorRuntimeLazyActivateGuard
        ],
        data: {
            localizationTitleKey: 'Integration.ListPage.Header.Title'
        }
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class IntegrationFeatureRoutingModule {}
