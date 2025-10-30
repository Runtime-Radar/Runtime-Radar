import { HTTP_INTERCEPTORS } from '@angular/common/http';

import { AuthHeadersInterceptor } from '../interceptors/auth-headers.interceptor';

const authHeadersInterceptorProvider = {
    provide: HTTP_INTERCEPTORS,
    useClass: AuthHeadersInterceptor,
    multi: true
};

export const AUTH_HEADERS_PROVIDER = [authHeadersInterceptorProvider];
