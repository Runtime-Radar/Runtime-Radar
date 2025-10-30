import { EntityState } from '@ngrx/entity';

import { LoadStatus } from '@cs/core';

import { Detector, DetectorType } from '../contract/detector-contract.interface';

export interface DetectorExtended extends Detector {
    key: string;
    type: DetectorType;
}

export interface DetectorConfig {
    id: DetectorType;
    loadStatus: LoadStatus;
    lastUpdate: number;
}

export type DetectorConfigEntityState = EntityState<DetectorConfig>;
export type DetectorEntityState = EntityState<DetectorExtended>;

export interface DetectorState {
    config: DetectorConfigEntityState;
    list: DetectorEntityState;
}
