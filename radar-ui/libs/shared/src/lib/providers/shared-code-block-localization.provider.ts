import { KBQ_CODE_BLOCK_LOCALE_CONFIGURATION } from '@koobiq/components/code-block';
import { KBQ_LOCALE_SERVICE } from '@koobiq/components/core';
import { Provider } from '@angular/core';

export const sharedCodeBlockLocalizationProvider = (): Provider[] => [
    {
        provide: KBQ_CODE_BLOCK_LOCALE_CONFIGURATION,
        useFactory: () => ({
            softWrapOnTooltip: 'Enable word wrap',
            softWrapOffTooltip: 'Disable word wrap',
            downloadTooltip: 'Download',
            copiedTooltip: 'Copied',
            copyTooltip: 'Copy',
            viewAllText: 'Open in the external system',
            viewLessText: 'Show all',
            openExternalSystemTooltip: 'Show less'
        })
    },
    {
        provide: KBQ_LOCALE_SERVICE,
        useValue: null
    }
];
