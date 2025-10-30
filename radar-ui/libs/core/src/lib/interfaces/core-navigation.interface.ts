import { RouterName } from '../constants/core-router.constant';

export interface NavigationMenuChild {
    path: RouterName;
    localizationKey: string;
    testId: string;
    icon: string;
}

export interface NavigationMenu {
    path: RouterName;
    localizationKey?: string;
    children: NavigationMenuChild[];
}
