import { RuntimeEventProcessorHistoryControl } from './contract/runtime-event-processor-contract.interface';
import { RuntimeEventType } from './contract/runtime-event-contract.interface';

export enum RuntimeContext {
    CURRENT = 1,
    PARENT = 2,
    SIBLING = 3,
    CHILDREN = 4
}

export interface RuntimeEventProcessorOption {
    id: RuntimeEventProcessorHistoryControl;
    localizationKey: string;
}

export interface RuntimeEventTypeOption {
    id: RuntimeEventType;
    localizationKey: string;
}

export interface RuntimeContextOption {
    id: RuntimeContext;
    localizationKey: string;
}
