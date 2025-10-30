import { HTTP_INTERCEPTORS } from '@angular/common/http';

import { ApiErrorInterceptor } from '../interceptors/api.interceptor';

const apiInterceptorProvider = {
    provide: HTTP_INTERCEPTORS,
    useClass: ApiErrorInterceptor,
    multi: true
};

export const HTTP_PROVIDER = [apiInterceptorProvider];
