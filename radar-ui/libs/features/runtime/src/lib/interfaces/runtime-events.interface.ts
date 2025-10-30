import { RuleSeverity } from '@cs/domains/rule';
import {
    GetRuntimeEventsResponse,
    RuntimeContext,
    RuntimeEvent,
    RuntimeEventCursorDirection
} from '@cs/domains/runtime';

export interface RuntimeEventsPagination {
    direction: RuntimeEventCursorDirection;
    cursor: string; // RFC3339
}

export interface RuntimeEventsGridContextId {
    id: string;
    context: RuntimeContext;
    execId: string;
    parentExecId: string;
}

export interface RuntimeEventExtended extends RuntimeEvent {
    threatSeverity?: RuleSeverity;
}

export interface GetRuntimeEventsResponseExtended extends GetRuntimeEventsResponse {
    runtime_events: RuntimeEventExtended[];
    isPrevResponse: boolean;
}
