import { DateAdapter } from '@koobiq/components/core';
import { DateTime } from 'luxon';
import { Pipe, PipeTransform } from '@angular/core';

import { I18nService } from '@cs/i18n';

@Pipe({
    name: 'dateFormatter',
    pure: false
})
export class SharedDateFormatterPipe implements PipeTransform {
    constructor(
        private readonly dateAdapter: DateAdapter<DateTime>,
        private readonly i18nService: I18nService
    ) {}

    transform(originDate: string, format?: Intl.DateTimeFormatOptions | string): string {
        const deserializedDate = this.dateAdapter.deserialize(originDate);
        if (!deserializedDate) {
            console.warn('date must be valid');

            return '';
        }

        const date = deserializedDate.setLocale(this.i18nService.getLocale());
        if (typeof format === 'string') {
            return date.toFormat(format);
        }

        const str = format
            ? date.toLocaleString(format)
            : `${date.toLocaleString(DateTime.DATE_FULL)}, ${date.toLocaleString(DateTime.TIME_SIMPLE)}`;

        return str.replace(/GMT\+\d+/, '');
    }
}
