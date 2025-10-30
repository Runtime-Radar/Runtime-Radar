import { DateAdapter } from '@koobiq/components/core';
import { DateTime } from 'luxon';
import { ChangeDetectionStrategy, Component, EventEmitter, Output } from '@angular/core';

import { RUNTIME_FILTER_DATETIME_PERIOD_SEPARATOR } from '../../constants/runtime-filter.constant';
import { RuntimeEventFilters } from '../../interfaces/runtime-filter.interface';

@Component({
    selector: 'cs-runtime-feature-preset-dropdown-component',
    templateUrl: './runtime-preset-dropdown.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuntimeFeaturePresetDropdownComponent {
    @Output() presetChange = new EventEmitter<RuntimeEventFilters>();

    presetDropdownFiltersCollection: Map<RuntimeEventFilters, string> = new Map();

    constructor(private readonly dateAdapter: DateAdapter<DateTime>) {}

    select(filters: RuntimeEventFilters) {
        this.presetChange.emit(filters);
    }

    applyItems() {
        const today = this.dateAdapter.today().startOf('second');

        this.presetDropdownFiltersCollection = new Map([
            [
                this.getFilterSettings(today.minus({ minute: 1 }), today, true),
                'Runtime.EventsPage.Filter.PresetDrowdown.Option.OneMinuteWithThreats'
            ],
            [
                this.getFilterSettings(today.minus({ hour: 1 }), today, true),
                'Runtime.EventsPage.Filter.PresetDrowdown.Option.OneHourWithThreats'
            ],
            [
                this.getFilterSettings(today.minus({ day: 1 }), today, true),
                'Runtime.EventsPage.Filter.PresetDrowdown.Option.OneDayWithThreats'
            ],
            [
                this.getFilterSettings(today.minus({ week: 1 }), today, true),
                'Runtime.EventsPage.Filter.PresetDrowdown.Option.OneWeekWithThreats'
            ],
            [
                this.getFilterSettings(today.minus({ minute: 1 }), today),
                'Runtime.EventsPage.Filter.PresetDrowdown.Option.OneMinute'
            ],
            [
                this.getFilterSettings(today.minus({ hour: 1 }), today),
                'Runtime.EventsPage.Filter.PresetDrowdown.Option.OneHour'
            ],
            [
                this.getFilterSettings(today.minus({ day: 1 }), today),
                'Runtime.EventsPage.Filter.PresetDrowdown.Option.OneDay'
            ],
            [
                this.getFilterSettings(today.minus({ week: 1 }), today),
                'Runtime.EventsPage.Filter.PresetDrowdown.Option.OneWeek'
            ]
        ]);
    }

    private getFilterSettings(from: DateTime, to: DateTime, hasThreats = false): RuntimeEventFilters {
        return {
            type: null,
            argument: '',
            binary: '',
            container: '',
            function: '',
            image: '',
            namespace: '',
            pod: '',
            period: from
                .toJSDate()
                .toISOString()
                .concat(RUNTIME_FILTER_DATETIME_PERIOD_SEPARATOR, to.toJSDate().toISOString()), // RFC3339
            hasThreats,
            hasIncident: false,
            detectors: [],
            rules: []
        };
    }
}
