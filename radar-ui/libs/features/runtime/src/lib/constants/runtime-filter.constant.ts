import {
    RuntimeEventContext,
    RuntimeEventDateTimePeriod,
    RuntimeEventDateTimePeriodOption,
    RuntimeEventFilters
} from '../interfaces/runtime-filter.interface';

export const RUNTIME_FILTER_LABEL_LIMIT = 3;

export const RUNTIME_FILTER_DATETIME_PERIOD_SEPARATOR = '|';

export const RUNTIME_FILTER_INITIAL_STATE: RuntimeEventFilters = {
    type: null,
    argument: '',
    binary: '',
    container: '',
    function: '',
    image: '',
    namespace: '',
    pod: '',
    period: '',
    hasThreats: false,
    hasIncident: false,
    detectors: [],
    rules: []
};

export const RUNTIME_FILTER_INITIAL_CONTEXT_STATE: RuntimeEventContext = {
    activeContextId: undefined,
    context: undefined,
    execId: '',
    parentExecId: ''
};

export const RUNTIME_FILTER_DATETIME_PERIOD: RuntimeEventDateTimePeriodOption[] = [
    {
        id: RuntimeEventDateTimePeriod.ONE_MINUTE,
        localizationKey: 'Runtime.Pseudo.DateTimePeriod.OneMinute'
    },
    {
        id: RuntimeEventDateTimePeriod.TEN_MINUTES,
        localizationKey: 'Runtime.Pseudo.DateTimePeriod.TenMinutes'
    },
    {
        id: RuntimeEventDateTimePeriod.ONE_HOUR,
        localizationKey: 'Runtime.Pseudo.DateTimePeriod.OneHour'
    },
    {
        id: RuntimeEventDateTimePeriod.ONE_DAY,
        localizationKey: 'Runtime.Pseudo.DateTimePeriod.OneDay'
    },
    {
        id: RuntimeEventDateTimePeriod.CUSTOM,
        localizationKey: 'Runtime.Pseudo.DateTimePeriod.Custom'
    }
];
