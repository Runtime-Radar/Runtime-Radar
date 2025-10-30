import { FormArray, FormControl, FormRecord } from '@angular/forms';

type Includes<T, K> = K extends T ? true : false;

/**
 * @arg - FormScheme<FormInterface, FormRecord fields, FormArray fields>
 */
export type FormScheme<T, R extends keyof T = never, A extends keyof T = never> = {
    [K in keyof T]: Includes<R, K> extends true
        ? T[K] extends infer C
            ? FormRecord<FormControl<C | null>>
            : never
        : Includes<A, K> extends true
          ? T[K] extends (infer D)[]
              ? FormArray<FormControl<D | null>>
              : never
          : FormControl<T[K] | null>;
};
