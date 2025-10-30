import { Injectable } from '@angular/core';
import { Observable, filter, map, switchMap, take } from 'rxjs';

import { ApiEmptyRequest, ApiService } from '@cs/api';

import {
    CreateRuntimeEventProcessorRequest,
    CreateRuntimeMonitorRequest,
    EmptyRuntimeResponse,
    GetRuntimeEventCountResponse,
    GetRuntimeEventsByFilterRequest,
    GetRuntimeEventsRequest,
    GetRuntimeEventsResponse,
    RuntimeEvent,
    RuntimeEventCursorDirection,
    RuntimeEventProcessor,
    RuntimeEventProcessorConfig,
    RuntimeEventProcessorHistoryControl,
    RuntimeMonitor,
    RuntimeMonitorConfig
} from '../interfaces';

@Injectable({
    providedIn: 'root'
})
export class RuntimeRequestService {
    constructor(private readonly apiService: ApiService) {}

    /** @external */
    getEventCount(from: string, to: string): Observable<number> {
        return this.apiService
            .get<
                ApiEmptyRequest,
                GetRuntimeEventCountResponse
            >(`stats/runtime-event/count?period.from=${from}&period.to=${to}`)
            .pipe(map((response) => response.count));
    }

    getRuntimeMonitor(): Observable<RuntimeMonitor> {
        return this.apiService.get<ApiEmptyRequest, RuntimeMonitor>('config/runtime-monitor');
    }

    getEventProcessor(): Observable<RuntimeEventProcessor> {
        return this.apiService.get<ApiEmptyRequest, RuntimeEventProcessor>('config/event-processor');
    }

    createRuntimeMonitor(config: RuntimeMonitorConfig): Observable<RuntimeMonitor> {
        return this.apiService
            .post<CreateRuntimeMonitorRequest, EmptyRuntimeResponse>('config/runtime-monitor', { config })
            .pipe(
                map((response) => response && !Object.keys(response).length),
                filter((isCreated) => isCreated),
                switchMap(() => this.getRuntimeMonitor().pipe(take(1)))
            );
    }

    createEventProcessor(historyControl: RuntimeEventProcessorHistoryControl): Observable<RuntimeEventProcessor> {
        const config: RuntimeEventProcessorConfig = {
            version: '1', // @todo: create environment constant
            history_control: historyControl
        };

        return this.apiService
            .post<CreateRuntimeEventProcessorRequest, EmptyRuntimeResponse>('config/event-processor', { config })
            .pipe(
                map((response) => response && !Object.keys(response).length),
                filter((isCreated) => isCreated),
                switchMap(() => this.getEventProcessor().pipe(take(1)))
            );
    }

    /** @external */
    getEvents(
        direction: RuntimeEventCursorDirection,
        request: GetRuntimeEventsRequest
    ): Observable<GetRuntimeEventsResponse> {
        return this.apiService.get<GetRuntimeEventsRequest, GetRuntimeEventsResponse>(
            `runtime-event/slice/${direction}`,
            request
        );
    }

    /** @external */
    getEventsByFilter(
        direction: RuntimeEventCursorDirection,
        request: GetRuntimeEventsByFilterRequest
    ): Observable<GetRuntimeEventsResponse> {
        return this.apiService.post<GetRuntimeEventsByFilterRequest, GetRuntimeEventsResponse>(
            `runtime-event/by-filter/slice/${direction}`,
            request
        );
    }

    /** @external */
    getEvent(id: string): Observable<RuntimeEvent> {
        return this.apiService.get<ApiEmptyRequest, RuntimeEvent>(`runtime-event/${id}`);
    }
}
