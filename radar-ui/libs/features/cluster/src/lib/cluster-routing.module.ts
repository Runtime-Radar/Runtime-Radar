import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { clusterActivateGuard } from '@cs/domains/cluster';

import { ClusterFeatureCreateContainer } from './containers/create/cluster-create.container';
import { ClusterFeatureDetailsContainer } from './containers/details/cluster-details.container';
import { ClusterFeatureListContainer } from './containers/list/cluster-list.container';
import { ClusterRouterName } from './interfaces/cluster-navigation.interface';
import { clusterFeatureCreateRouteActivateGuard } from './guards/cluster-create-route-activate.guard';
import { clusterFeatureDetailsResolver } from './resolvers/cluster-details.resolver';

const routes: Routes = [
    {
        path: '',
        component: ClusterFeatureListContainer,
        canActivate: [clusterActivateGuard],
        data: {
            localizationTitleKey: 'Cluster.ListPage.Header.Title'
        }
    },
    {
        path: ClusterRouterName.CREATE,
        component: ClusterFeatureCreateContainer,
        canActivate: [clusterActivateGuard, clusterFeatureCreateRouteActivateGuard],
        data: {
            localizationTitleKey: 'Cluster.CreatePage.Header.Title'
        }
    },
    {
        path: ':clusterId',
        component: ClusterFeatureDetailsContainer,
        resolve: {
            cluster: clusterFeatureDetailsResolver
        },
        data: {
            localizationTitleKey: 'Cluster.DetailsPage.Header.Title'
        }
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class ClusterFeatureRoutingModule {}
