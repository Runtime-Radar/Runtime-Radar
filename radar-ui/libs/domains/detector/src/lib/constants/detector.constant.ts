import { DetectorType } from '../interfaces/contract/detector-contract.interface';
import { DetectorTypeOption } from '../interfaces/detector-type.interface';

export const DETECTOR_TYPE: DetectorTypeOption[] = [
    {
        id: DetectorType.RUNTIME,
        localizationKey: 'Common.Pseudo.DetectorType.Runtime'
    }
];
