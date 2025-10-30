import { provideTranslocoMessageformat } from '@jsverse/transloco-messageformat';
import {
    DateAdapter,
    DateFormatter,
    KBQ_DATE_FORMATS,
    KBQ_DATE_LOCALE,
    KBQ_LOCALE_SERVICE,
    KbqFormattersModule,
    KbqLocaleService
} from '@koobiq/components/core';
import {
    KBQ_LUXON_DATE_ADAPTER_OPTIONS,
    KBQ_LUXON_DATE_FORMATS,
    LuxonDateAdapter
} from '@koobiq/angular-luxon-adapter/adapter';
import { LOCALE_ID, ModuleWithProviders, NgModule } from '@angular/core';
import { TranslocoConfig, TranslocoModule, provideTransloco } from '@jsverse/transloco';
import {
    TranslocoMarkupModule,
    defaultTranslocoMarkupTranspilers,
    provideTranslationMarkupTranspiler
} from 'ngx-transloco-markup';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

import { I18nTemplateTranspiler } from './transpilers/i18n-template.transpiler';
import { I18nTranslocoLoader } from './providers/i18n-transloco.loader';
import { I18N_AVAILABLE_LOCALES, I18N_DEFAULT_LOCALE } from './constants/i18n.constant';

const DEFAULT_TRANSLOCO_CONFIG: Partial<TranslocoConfig> = {
    availableLangs: I18N_AVAILABLE_LOCALES,
    defaultLang: I18N_DEFAULT_LOCALE,
    fallbackLang: I18N_DEFAULT_LOCALE,
    interpolation: ['<<<', '>>>'],
    reRenderOnLangChange: true
};

const translocoMarkupTranspilers = [
    defaultTranslocoMarkupTranspilers(),
    provideTranslationMarkupTranspiler(I18nTemplateTranspiler)
];

@NgModule({
    imports: [KbqFormattersModule, TranslocoMarkupModule, TranslocoModule],
    providers: [
        {
            provide: LOCALE_ID,
            useValue: I18N_DEFAULT_LOCALE
        },
        {
            provide: DateFormatter,
            deps: [DateAdapter, KBQ_DATE_LOCALE]
        },
        ...translocoMarkupTranspilers
    ],
    exports: [TranslocoMarkupModule, TranslocoModule]
})
export class I18nModule {
    static forRoot(config: Partial<TranslocoConfig> = {}): ModuleWithProviders<I18nModule> {
        return {
            ngModule: I18nModule,
            providers: [
                provideTransloco({
                    config: {
                        ...DEFAULT_TRANSLOCO_CONFIG,
                        ...config
                    },
                    loader: I18nTranslocoLoader
                }),
                provideTranslocoMessageformat({
                    locales: I18N_AVAILABLE_LOCALES
                }),
                {
                    provide: KBQ_DATE_FORMATS,
                    useValue: KBQ_LUXON_DATE_FORMATS
                },
                {
                    provide: DateAdapter,
                    useClass: LuxonDateAdapter,
                    deps: [KBQ_DATE_LOCALE, KBQ_LUXON_DATE_ADAPTER_OPTIONS, KBQ_LOCALE_SERVICE]
                },
                {
                    provide: KBQ_LOCALE_SERVICE,
                    useClass: KbqLocaleService
                },
                provideHttpClient(withInterceptorsFromDi())
            ]
        };
    }
}
