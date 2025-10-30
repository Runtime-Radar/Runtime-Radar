import { KBQ_SIDEPANEL_DATA } from '@koobiq/components/sidepanel';
import { KbqCodeBlockFile } from '@koobiq/components/code-block';
import { ChangeDetectionStrategy, Component, Inject, OnInit } from '@angular/core';

import { RuntimeSidepanelPolicyProps } from '../../interfaces/runtime-sidepanel.interface';

@Component({
    templateUrl: './runtime-sidepanel-policy.component.html',
    styleUrl: './runtime-sidepanel-policy.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuntimeFeatureSidepanelPolicyComponent implements OnInit {
    files: KbqCodeBlockFile[] = [];

    constructor(@Inject(KBQ_SIDEPANEL_DATA) public readonly props: Partial<RuntimeSidepanelPolicyProps>) {}

    ngOnInit() {
        if (this.props.yaml) {
            this.files.push({
                filename: 'source.yaml',
                content: this.props.yaml,
                language: 'yaml'
            });
        }
    }
}
