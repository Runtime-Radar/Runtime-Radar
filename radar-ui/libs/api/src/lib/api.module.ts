import { NgModule } from '@angular/core';

import { HTTP_PROVIDER } from './providers/api.provider';

@NgModule({
    imports: [],
    providers: [HTTP_PROVIDER]
})
export class ApiModule {}
