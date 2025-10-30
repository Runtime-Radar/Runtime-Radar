import { ChangeDetectionStrategy, Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Observable, debounceTime, distinctUntilChanged, map, startWith } from 'rxjs';

import { AuthStoreService } from '@cs/domains/auth';
import { FormScheme, LoadStatus, CoreUtilsService as utils } from '@cs/core';

import { SignInForm } from '../../interfaces/sign-in-form.interface';

@Component({
    templateUrl: './sign-in-page.container.html',
    styleUrl: './sign-in-page.container.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class SignInFeaturePageContainer {
    readonly isSignInProgress$ = this.authStoreService.loadStatus$.pipe(
        map((loadStatus) => loadStatus === LoadStatus.IN_PROGRESS)
    );

    readonly form: FormGroup<FormScheme<SignInForm>> = this.formBuilder.group({
        username: ['', Validators.required],
        password: ['', Validators.required],
        rememberMe: [true]
    });

    readonly isFormValid$: Observable<boolean> = this.form.valueChanges.pipe(
        startWith(this.form.value),
        /* eslint @typescript-eslint/no-magic-numbers: "off" */
        debounceTime(250),
        distinctUntilChanged(),
        map(() => utils.isFormValid(this.form.controls))
    );

    constructor(
        private readonly formBuilder: FormBuilder,
        private readonly authStoreService: AuthStoreService
    ) {}

    signIn() {
        const formValues = utils.getFormValues<SignInForm>(this.form.controls);
        const { username, password } = utils.getTrimmedFormValues<SignInForm>(formValues);

        if (username && password) {
            this.authStoreService.signIn(username, password);
        }
    }
}
