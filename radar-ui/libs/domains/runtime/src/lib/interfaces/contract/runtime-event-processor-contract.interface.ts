export enum RuntimeEventProcessorHistoryControl {
    NONE = 'NONE',
    ALL = 'ALL',
    WITH_THREATS = 'WITH_THREATS'
}

export interface RuntimeEventProcessorConfig {
    version: string;
    history_control: RuntimeEventProcessorHistoryControl;
}

export interface RuntimeEventProcessor {
    id: string;
    config: RuntimeEventProcessorConfig;
}
