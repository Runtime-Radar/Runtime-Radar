import { HttpErrorResponse } from '@angular/common/http';
import { Router } from '@angular/router';
import { ErrorHandler, Injectable } from '@angular/core';
import { KbqToastService, KbqToastStyle } from '@koobiq/components/toast';

import { I18nService } from '@cs/i18n';
import { RouterName } from '@cs/core';

import { AuthLocalStorageService } from '../services/auth-local-storage.services';

@Injectable({
    providedIn: 'root'
})
export class AuthChunksLoadHandler implements ErrorHandler {
    constructor(
        private readonly authLocalStorageService: AuthLocalStorageService,
        private readonly i18nService: I18nService,
        private readonly toastService: KbqToastService,
        private readonly router: Router
    ) {}

    handleError(error: HttpErrorResponse) {
        const chunkRegExp = /Loading chunk [^\s]+ failed/;
        if (chunkRegExp.test(error.message)) {
            this.toastService.show({
                style: KbqToastStyle.Warning,
                title: this.i18nService.translate('Common.Pseudo.Notification.ChunkLoadFailed'),
                caption: this.i18nService.translate('Common.Pseudo.Notification.UpdatePage')
            });

            // @todo: reset auth state
            this.authLocalStorageService.removeTokens();
            this.router.navigate([RouterName.SIGN_IN]);
        }

        console.error(error);
    }
}
