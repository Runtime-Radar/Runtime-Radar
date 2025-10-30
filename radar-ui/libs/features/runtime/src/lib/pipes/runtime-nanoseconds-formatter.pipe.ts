import { DateAdapter } from '@koobiq/components/core';
import { DateTime } from 'luxon';
import { Pipe, PipeTransform } from '@angular/core';

import { I18nService } from '@cs/i18n';

const DEFAULT_FRACTION = 9;

@Pipe({
    name: 'runtimeNanosecondsFormatter',
    pure: false
})
export class RuntimeFeatureNanosecondsFormatterPipe implements PipeTransform {
    constructor(
        private readonly dateAdapter: DateAdapter<DateTime>,
        private readonly i18nService: I18nService
    ) {}

    transform(date: string, fraction = DEFAULT_FRACTION): string {
        const deserializedDate = this.dateAdapter.deserialize(date);
        if (!deserializedDate) {
            console.warn('date must be valid');

            return '';
        }

        /* eslint @typescript-eslint/no-magic-numbers: "off" */
        return (
            deserializedDate.setLocale(this.i18nService.getLocale()).toFormat('D, TT.') +
            date.substring(20, fraction >= 9 ? 29 : 20 + fraction)
        );
    }
}
