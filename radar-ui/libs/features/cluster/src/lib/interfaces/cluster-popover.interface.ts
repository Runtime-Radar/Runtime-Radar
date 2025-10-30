import { ClusterConfig } from '@cs/domains/cluster';

export interface ClusterEditPopoverOutputs {
    id: string;
    name: string;
    config: ClusterConfig;
}
