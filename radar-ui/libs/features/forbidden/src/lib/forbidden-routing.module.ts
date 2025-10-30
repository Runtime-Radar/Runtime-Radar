import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ForbiddenFeaturePermissionContainer } from './containers/permission/forbidden-permission.container';

const routes: Routes = [
    {
        path: '',
        component: ForbiddenFeaturePermissionContainer,
        data: {
            localizationTitleKey: 'Common.ForbiddenPage.Header.Title'
        }
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class ForbiddenFeatureRoutingModule {}
