import {
    RuntimeContext,
    RuntimeContextOption,
    RuntimeEventProcessorHistoryControl,
    RuntimeEventProcessorOption,
    RuntimeEventType,
    RuntimeEventTypeOption
} from '../interfaces';

export const RUNTIME_EVENT_PROCESSOR: RuntimeEventProcessorOption[] = [
    {
        id: RuntimeEventProcessorHistoryControl.ALL,
        localizationKey: 'Runtime.Pseudo.EventProcessorHistory.All'
    },
    {
        id: RuntimeEventProcessorHistoryControl.WITH_THREATS,
        localizationKey: 'Runtime.Pseudo.EventProcessorHistory.WithThreats'
    },
    {
        id: RuntimeEventProcessorHistoryControl.NONE,
        localizationKey: 'Runtime.Pseudo.EventProcessorHistory.None'
    }
];

export const RUNTIME_EVENT_TYPE: RuntimeEventTypeOption[] = [
    {
        id: RuntimeEventType.EXEC,
        localizationKey: 'Runtime.Pseudo.EventType.Exec'
    },
    {
        id: RuntimeEventType.EXIT,
        localizationKey: 'Runtime.Pseudo.EventType.Exit'
    },
    {
        id: RuntimeEventType.KPROBE,
        localizationKey: 'Runtime.Pseudo.EventType.Kprobe'
    },
    {
        id: RuntimeEventType.LOADER,
        localizationKey: 'Runtime.Pseudo.EventType.Loader'
    },
    {
        id: RuntimeEventType.TRACEPOINT,
        localizationKey: 'Runtime.Pseudo.EventType.Tracepoint'
    },
    {
        id: RuntimeEventType.UPROBE,
        localizationKey: 'Runtime.Pseudo.EventType.Uprobe'
    }
];

export const RUNTIME_CONTEXT: RuntimeContextOption[] = [
    {
        id: RuntimeContext.PARENT,
        localizationKey: 'Runtime.Pseudo.Context.Parent'
    },
    {
        id: RuntimeContext.SIBLING,
        localizationKey: 'Runtime.Pseudo.Context.Sibling'
    },
    {
        id: RuntimeContext.CHILDREN,
        localizationKey: 'Runtime.Pseudo.Context.Children'
    },
    {
        id: RuntimeContext.CURRENT,
        localizationKey: 'Runtime.Pseudo.Context.Current'
    }
];
