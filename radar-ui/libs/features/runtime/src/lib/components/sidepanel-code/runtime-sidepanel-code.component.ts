import { KBQ_SIDEPANEL_DATA } from '@koobiq/components/sidepanel';
import { KbqCodeBlockFile } from '@koobiq/components/code-block';
import { ChangeDetectionStrategy, Component, Inject, OnInit } from '@angular/core';

import { RuntimeSidepanelCodeProps } from '../../interfaces/runtime-sidepanel.interface';

const RUNTIME_CODE_SPACE_INDENT = 4;

@Component({
    templateUrl: './runtime-sidepanel-code.component.html',
    styleUrl: './runtime-sidepanel-code.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuntimeFeatureSidepanelCodeComponent implements OnInit {
    files: KbqCodeBlockFile[] = [];

    constructor(@Inject(KBQ_SIDEPANEL_DATA) public readonly props: RuntimeSidepanelCodeProps) {}

    ngOnInit() {
        this.files.push({
            filename: 'event.json',
            content: this.format(this.props.content),
            language: 'json'
        });
    }

    private format(content: string) {
        return JSON.stringify(JSON.parse(content), null, RUNTIME_CODE_SPACE_INDENT);
    }
}
