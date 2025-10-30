import { ClusterStepName } from './cluster-stepper.interface';
import {
    ClusterAccessForm,
    ClusterDataBaseForm,
    ClusterIngressForm,
    ClusterRabbitForm,
    ClusterRegistryForm
} from './cluster-form.interface';

export interface ClusterFormState {
    id: number;
    step: ClusterStepName;
    registry: ClusterRegistryForm;
    clickhouse: ClusterDataBaseForm;
    postgres: ClusterDataBaseForm;
    redis: ClusterDataBaseForm;
    rabbit: ClusterRabbitForm;
    ingress: ClusterIngressForm;
    access: ClusterAccessForm;
}

export type ClusterFormType = keyof Omit<ClusterFormState, 'id' | 'step'>;
