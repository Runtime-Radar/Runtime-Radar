import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { detectorRuntimeLazyActivateGuard } from '@cs/domains/detector';
import { notificationLazyActivateGuard } from '@cs/domains/notification';
import { ruleActivateGuard } from '@cs/domains/rule';

import { RuleFeatureListContainer } from './containers/list/rule-list.container';

const routes: Routes = [
    {
        path: '',
        component: RuleFeatureListContainer,
        canActivate: [ruleActivateGuard, notificationLazyActivateGuard, detectorRuntimeLazyActivateGuard],
        data: {
            localizationTitleKey: 'Rule.ListPage.Header.Title'
        }
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class RuleFeatureRoutingModule {}
