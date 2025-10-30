export enum ClusterStepName {
    REGISTRY = 'registry',
    DATABASE = 'database',
    INGRESS = 'ingress',
    ACCESS = 'access'
}

export interface ClusterStepperTab {
    id: ClusterStepName;
    title: string;
    description: string;
}
