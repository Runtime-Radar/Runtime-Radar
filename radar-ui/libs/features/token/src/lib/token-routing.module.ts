import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { tokenActivateGuard } from '@cs/domains/token';

import { TokenFeatureListContainer } from './containers/list/token-list.container';

const routes: Routes = [
    {
        path: '',
        component: TokenFeatureListContainer,
        canActivate: [tokenActivateGuard],
        data: {
            localizationTitleKey: 'Token.ListPage.Header.Title'
        }
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class TokenFeatureRoutingModule {}
