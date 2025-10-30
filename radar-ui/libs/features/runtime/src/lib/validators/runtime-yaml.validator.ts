import { parseDocument } from 'yaml';
import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

export class RuntimeFeatureYamlValidator {
    static isYamlCodeValid(): ValidatorFn {
        return (control: AbstractControl): ValidationErrors | null => {
            if (!control.value) {
                return null;
            }

            const result = parseDocument(control.value);
            if (!result.errors.length) {
                return null;
            }

            const errors = result.errors.map((item) => ({
                type: item.code,
                line: item.linePos ? item.linePos.at(0)?.line : undefined
            }));

            return {
                code: {
                    type: errors.at(0)?.type,
                    line: errors.at(0)?.line
                }
            };
        };
    }
}
