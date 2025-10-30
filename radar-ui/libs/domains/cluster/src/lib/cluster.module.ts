import { EffectsModule } from '@ngrx/effects';
import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';

import { ApiModule } from '@cs/api';

import { ClusterEffectStore } from './stores/cluster-effect.store';
import { ClusterStatusLocalizationPipe } from './pipes/cluster-status.pipe';
import { CLUSTER_DOMAIN_KEY, clusterDomainReducer } from './stores/cluster-selector.store';

@NgModule({
    imports: [
        ApiModule,
        StoreModule.forFeature(CLUSTER_DOMAIN_KEY, clusterDomainReducer),
        EffectsModule.forFeature([ClusterEffectStore])
    ],
    declarations: [ClusterStatusLocalizationPipe],
    exports: [ClusterStatusLocalizationPipe]
})
export class ClusterDomainModule {}
