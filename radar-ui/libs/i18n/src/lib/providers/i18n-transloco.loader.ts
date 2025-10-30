import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
// eslint-disable-next-line import/no-unresolved
import { TranslocoLoaderData } from '@jsverse/transloco/lib/transloco.loader';
import { EMPTY, Observable, catchError, of } from 'rxjs';
import { Translation, TranslocoLoader } from '@jsverse/transloco';

@Injectable({
    providedIn: 'root'
})
export class I18nTranslocoLoader implements TranslocoLoader {
    private readonly http = inject(HttpClient);

    getTranslation(path: string, data?: TranslocoLoaderData): Observable<Translation> {
        if (!data) {
            return of(EMPTY);
        }

        const locale = path.split('/').reverse()[0];
        const translationPath = `/assets/i18n/${locale}/${data?.scope}.json`;

        return this.http.get<Translation>(translationPath).pipe(catchError(() => of({})));
    }
}
