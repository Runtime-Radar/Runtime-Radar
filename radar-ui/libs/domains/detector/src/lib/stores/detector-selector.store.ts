import { ActionReducerMap, createFeatureSelector, createSelector } from '@ngrx/store';

import { DetectorType } from '../interfaces';
import {
    DetectorConfigEntityState,
    DetectorEntityState,
    DetectorState
} from '../interfaces/state/detector-state.interface';
import { detectorListEntitySelector, detectorReducer } from './detector-reducer.store';

export const DETECTOR_DOMAIN_KEY = 'detector';

export interface DetectorDomainState {
    readonly domain: DetectorState;
}

const selectDetectorDomainState = createFeatureSelector<DetectorDomainState>(DETECTOR_DOMAIN_KEY);
const selectDetectorState = createSelector(selectDetectorDomainState, (state: DetectorDomainState) => state.domain);
const selectDetectorListEntityState = createSelector(selectDetectorState, (state: DetectorState) => state.list);
const selectDetectorConfigEntityState = createSelector(selectDetectorState, (state: DetectorState) => state.config);

export const getDetectorConfigs = createSelector(selectDetectorConfigEntityState, (state: DetectorConfigEntityState) =>
    Object.entries(state.entities).map(([_, value]) => value)
);

export const getDetectorConfig = (type: DetectorType) =>
    createSelector(selectDetectorConfigEntityState, (state: DetectorConfigEntityState) => state.entities[type]);

export const getDetectors = (types?: DetectorType[]) =>
    createSelector(selectDetectorListEntityState, (state: DetectorEntityState) => {
        const detectors = detectorListEntitySelector.selectAll(state);

        return types && types.length ? detectors.filter((item) => types.includes(item.type)) : detectors;
    });

export const detectorDomainReducer: ActionReducerMap<DetectorDomainState> = {
    domain: detectorReducer
};
