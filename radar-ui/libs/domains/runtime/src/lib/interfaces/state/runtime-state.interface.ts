import { LoadStatus } from '@cs/core';

import { RuntimeEventProcessorHistoryControl } from '../contract/runtime-event-processor-contract.interface';
import { RuntimeMonitorConfig } from '../contract/runtime-monitor-contract.interface';

export enum RuntimeConfigStatus {
    INIT = 'INIT',
    STAY = 'STAY',
    MODIFY = 'MODIFY'
}

export interface RuntimeState {
    loadStatus: LoadStatus;
    hasChanges: boolean;
    hasPoliciesChanges: boolean;
    configStatus: RuntimeConfigStatus;
    isExpertMode: boolean;
    isOverlayed: boolean;
    historyControl?: RuntimeEventProcessorHistoryControl;
    config: RuntimeMonitorConfig;
}
