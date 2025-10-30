import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { SwitchFeaturePageContainer } from './containers/page/switch-page.container';
import { clusterFeatureSwitchRouteActivateGuard } from './guards/cluster-switch-route-activate.guard';

const routes: Routes = [
    {
        path: '',
        component: SwitchFeaturePageContainer,
        canActivate: [clusterFeatureSwitchRouteActivateGuard],
        data: {
            localizationTitleKey: 'Cluster.SwitchPage.Header.Title'
        }
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class SwitchFeatureRoutingModule {}
