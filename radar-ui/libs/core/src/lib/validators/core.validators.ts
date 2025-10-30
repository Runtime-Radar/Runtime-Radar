import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

import { FORM_VALIDATION_REG_EXP } from '../constants';

const BIT_LENGTH = 8;

export class CoreValidators {
    static isIpSegmentAllowed(segment: string): ValidatorFn {
        const [ip, pr] = segment.split('/');
        const prefix = Number(pr) / BIT_LENGTH;

        const multiplyBinaryOctets = (a: string[], b: string[]): string[] =>
            a.length === b.length
                ? a.map((octet, i) =>
                      octet
                          .split('')
                          .map((n, j) => Number(n) * Number(b[i][j]))
                          .join('')
                  )
                : [];

        const getBinaryIpOctets = (address: string, shift?: number): string[] =>
            address
                .split('.')
                .filter((_, i) => !!shift && i < shift)
                .map((octet) => {
                    const binary = Number(octet).toString(2);

                    return (
                        Array(BIT_LENGTH - binary.length)
                            .fill('0')
                            .join('') + binary
                    );
                });

        const segmentBinaryIpOctets = getBinaryIpOctets(ip, prefix);
        const segmentBinaryMaskOctets = segmentBinaryIpOctets.map(() => Array(BIT_LENGTH).fill('1').join(''));

        return (control: AbstractControl): ValidationErrors | null => {
            const values = (control.value as string).match(FORM_VALIDATION_REG_EXP.IP);
            if (!values) {
                return null;
            }

            const binaryIpOctets = getBinaryIpOctets(values[0].split(':')[0], prefix);
            const cornerIpOctets = multiplyBinaryOctets(binaryIpOctets, segmentBinaryMaskOctets);
            if (cornerIpOctets.join('') !== segmentBinaryIpOctets.join('')) {
                return null;
            }

            return {
                segment: {
                    allowedValue: ip,
                    actualValue: values[0]
                }
            };
        };
    }
}
