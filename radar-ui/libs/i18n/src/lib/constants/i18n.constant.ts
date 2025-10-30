import { I18nLocale } from '../interfaces/i18n.interface';

export const I18N_LOCAL_STORAGE_KEY = 'locale';

export const I18N_AVAILABLE_LOCALES = [I18nLocale.EN];

export const I18N_DEFAULT_LOCALE = I18nLocale.EN;

export const I18N_DATE_LOCALE_FORMATS = {
    [I18nLocale.EN]: 'yyyy-MM-dd'
};
