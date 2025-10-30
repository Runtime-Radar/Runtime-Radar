import { ClusterStepName, ClusterStepperTab } from '../interfaces/cluster-stepper.interface';

export const CLUSTER_STEPPER_TABS: ClusterStepperTab[] = [
    {
        id: ClusterStepName.REGISTRY,
        title: 'Cluster.CreatePage.Tab.Label.Registry',
        description: 'Cluster.CreatePage.Tab.Hint.Registry'
    },
    {
        id: ClusterStepName.DATABASE,
        title: 'Cluster.CreatePage.Tab.Label.Database',
        description: 'Cluster.CreatePage.Tab.Hint.Database'
    },
    {
        id: ClusterStepName.INGRESS,
        title: 'Cluster.CreatePage.Tab.Label.Ingress',
        description: 'Cluster.CreatePage.Tab.Hint.Ingress'
    },
    {
        id: ClusterStepName.ACCESS,
        title: 'Cluster.CreatePage.Tab.Label.Access',
        description: 'Cluster.CreatePage.Tab.Hint.Access'
    }
];
