import { Injectable } from '@angular/core';

import { RuntimeEventProcessorHistoryControl, RuntimeMonitorConfig, RuntimeMonitorConfigExtended } from '../interfaces';

@Injectable({
    providedIn: 'root'
})
export class RuntimeHelperService {
    static convertConfigToDiffValues(
        config: RuntimeMonitorConfig,
        historyControl: RuntimeEventProcessorHistoryControl = RuntimeEventProcessorHistoryControl.NONE
    ): RuntimeMonitorConfigExtended {
        const { aggregation_options, ...configWithoutAggregationOptions } = config;
        const configValues: RuntimeMonitorConfig = {
            ...configWithoutAggregationOptions,
            allow_list: configWithoutAggregationOptions.allow_list.map(({ labels, namespace, pod_regex }) => ({
                labels,
                namespace,
                pod_regex
            })),
            deny_list: configWithoutAggregationOptions.deny_list.map(({ labels, namespace, pod_regex }) => ({
                labels,
                namespace,
                pod_regex
            }))
        };

        return { ...configValues, historyControl };
    }
}
