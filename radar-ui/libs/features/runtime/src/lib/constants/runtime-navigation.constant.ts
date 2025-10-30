import { RuntimeNavigationTab, RuntimeRouterName } from '../interfaces/runtime-navigation.interface';

export const RUNTIME_NAVIGATION_TABS: RuntimeNavigationTab[] = [
    {
        path: RuntimeRouterName.SETTINGS,
        localizationKey: 'Runtime.Pseudo.Navigation.Settings'
    },
    {
        path: RuntimeRouterName.EVENTS,
        localizationKey: 'Runtime.Pseudo.Navigation.Events'
    },
    {
        path: RuntimeRouterName.DETECTORS,
        localizationKey: 'Runtime.Pseudo.Navigation.Detectors'
    }
];
