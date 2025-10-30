import { RuntimeEventContext, RuntimeEventFilters } from './runtime-filter.interface';

export interface RuntimeEventFilterEntity extends RuntimeEventFilters, RuntimeEventContext {
    id: string; // RFC3339
}
