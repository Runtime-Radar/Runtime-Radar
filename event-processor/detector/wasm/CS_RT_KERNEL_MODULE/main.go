//go:build tinygo.wasm

package main

import (
	"context"

	"github.com/runtime-radar/runtime-radar/event-processor/detector/api"
	"github.com/runtime-radar/runtime-radar/event-processor/detector/api/tetragon"
)

const (
	ID          = "CS_RT_KERNEL_MODULE"
	Name        = "Kernel module loading"
	Description = "The detector detects loading of kernel modules that may contain malicious code used to attack the target system."
	Version     = 1
	Author      = "Runtime Radar Team"

	License = "Apache License 2.0"
)

// main is required for TinyGo to compile to Wasm.
func main() {
	api.RegisterDetector(Detector{})
}

type Detector struct{}

func (d Detector) Info(ctx context.Context, req *api.InfoReq) (*api.InfoResp, error) {
	return &api.InfoResp{
		Id:          ID,
		Name:        Name,
		Description: Description,
		Version:     Version,
		Author:      Author,
		Contact:     Contact,
		License:     License,
	}, nil
}

func (d Detector) Detect(ctx context.Context, req *api.DetectReq) (*api.DetectResp, error) {
	// Detector info added to DetectResp because detector info is always correlated to response, thus
	// to avoid +1 Wasm call on detect.
	resp := &api.DetectResp{
		Id:          ID,
		Name:        Name,
		Description: Description,
		Version:     Version,
		Author:      Author,
		Contact:     Contact,

		// Default response indicates that nothing detected (this is redundant and put here just for reference,
		// as Severity == api.DetectResp_NONE == 0 when omitted (default zero value)).
		Severity: api.DetectResp_NONE,
	}

	event := req.GetEvent().GetEvent()

	switch ev := event.(type) {
	case *tetragon.GetEventsResponse_ProcessExec:
		// Nothing here
	case *tetragon.GetEventsResponse_ProcessExit:
		// Nothing here
	case *tetragon.GetEventsResponse_ProcessKprobe:
		kprobe := ev.ProcessKprobe
		functionName := kprobe.GetFunctionName()

		if functionName == "do_init_module" {
			resp.Severity = api.DetectResp_HIGH // <-- threat detected

			return resp, nil
		}

	case *tetragon.GetEventsResponse_ProcessTracepoint:
		// Nothing here
	}

	return resp, nil
}

/* Example event (JSON):

{
    "process_kprobe": {
        "process": {
            "exec_id": "cHRleHBlcnRzLWs4cy1wdGNzOjM0NDM0NDE5NjUyMTg5NTk6MzI2Nzg0NQ==",
            "pid": 3267845,
            "uid": 0,
            "cwd": "/root/kernel_expl",
            "binary": "/usr/sbin/insmod",
            "arguments": "reverse-shell.ko",
            "flags": "execve clone",
            "start_time": "2024-03-24T07:44:28.758670169Z",
            "auid": 4294967295,
            "pod": {
                "namespace": "default",
                "name": "test-pod-modload-ubuntu",
                "labels": [],
                "container": {
                    "id": "cri-o://7ce6c223e9d5800d8610c1e1d2e0759a822de07bce7ae216b2c6250b4755bf9a",
                    "name": "test-pod-modload-ubuntu",
                    "image": {
                        "id": "docker.io/library/ubuntu@sha256:a4fab1802f08df089c4b2e0a1c8f1a06f573bd1775687d07fef4076d3a2e4900",
                        "name": "docker.io/library/ubuntu:focal"
                    },
                    "start_time": "2024-03-01T07:58:52Z",
                    "pid": 5622,
                    "maybe_exec_probe": false
                },
                "pod_labels": {},
                "workload": "test-pod-modload-ubuntu",
                "workload_kind": "Pod"
            },
            "docker": "7ce6c223e9d5800d8610c1e1d2e0759",
            "parent_exec_id": "cHRleHBlcnRzLWs4cy1wdGNzOjM0NDMyNzMwMzQ0NTk1OTI6MzI2NjEwMQ==",
            "refcnt": 1,
            "cap": {
                "permitted": [
                    "CAP_CHOWN",
                    "DAC_OVERRIDE",
                    "CAP_FOWNER",
                    "CAP_FSETID",
                    "CAP_KILL",
                    "CAP_SETGID",
                    "CAP_SETUID",
                    "CAP_SETPCAP",
                    "CAP_NET_BIND_SERVICE",
                    "CAP_SYS_MODULE"
                ],
                "effective": [
                    "CAP_CHOWN",
                    "DAC_OVERRIDE",
                    "CAP_FOWNER",
                    "CAP_FSETID",
                    "CAP_KILL",
                    "CAP_SETGID",
                    "CAP_SETUID",
                    "CAP_SETPCAP",
                    "CAP_NET_BIND_SERVICE",
                    "CAP_SYS_MODULE"
                ],
                "inheritable": []
            },
            "ns": {
                "uts": {
                    "inum": 4026533759,
                    "is_host": false
                },
                "ipc": {
                    "inum": 4026533760,
                    "is_host": false
                },
                "mnt": {
                    "inum": 4026534054,
                    "is_host": false
                },
                "pid": {
                    "inum": 4026534055,
                    "is_host": false
                },
                "pid_for_children": {
                    "inum": 4026534055,
                    "is_host": false
                },
                "net": {
                    "inum": 4026533761,
                    "is_host": false
                },
                "time": {
                    "inum": 4026531834,
                    "is_host": true
                },
                "time_for_children": {
                    "inum": 4026531834,
                    "is_host": true
                },
                "cgroup": {
                    "inum": 4026534056,
                    "is_host": false
                },
                "user": {
                    "inum": 4026531837,
                    "is_host": true
                }
            },
            "tid": 3267845,
            "process_credentials": {
                "uid": 0,
                "gid": 0,
                "euid": 0,
                "egid": 0,
                "suid": 0,
                "sgid": 0,
                "fsuid": 0,
                "fsgid": 0,
                "securebits": [],
                "caps": null,
                "user_ns": null
            },
            "binary_properties": null
        },
        "parent": {
            "exec_id": "cHRleHBlcnRzLWs4cy1wdGNzOjM0NDMyNzMwMzQ0NTk1OTI6MzI2NjEwMQ==",
            "pid": 3266101,
            "uid": 0,
            "cwd": "/",
            "binary": "/usr/bin/bash",
            "arguments": "",
            "flags": "execve rootcwd",
            "start_time": "2024-03-24T07:41:39.827910823Z",
            "auid": 4294967295,
            "pod": {
                "namespace": "default",
                "name": "test-pod-modload-ubuntu",
                "labels": [],
                "container": {
                    "id": "cri-o://7ce6c223e9d5800d8610c1e1d2e0759a822de07bce7ae216b2c6250b4755bf9a",
                    "name": "test-pod-modload-ubuntu",
                    "image": {
                        "id": "docker.io/library/ubuntu@sha256:a4fab1802f08df089c4b2e0a1c8f1a06f573bd1775687d07fef4076d3a2e4900",
                        "name": "docker.io/library/ubuntu:focal"
                    },
                    "start_time": "2024-03-01T07:58:52Z",
                    "pid": 5603,
                    "maybe_exec_probe": false
                },
                "pod_labels": {},
                "workload": "test-pod-modload-ubuntu",
                "workload_kind": "Pod"
            },
            "docker": "7ce6c223e9d5800d8610c1e1d2e0759",
            "parent_exec_id": "cHRleHBlcnRzLWs4cy1wdGNzOjM0NDMyNzMwMjU0MzI0NTI6MzI2NjEwMQ==",
            "refcnt": 0,
            "cap": {
                "permitted": [
                    "CAP_CHOWN",
                    "DAC_OVERRIDE",
                    "CAP_FOWNER",
                    "CAP_FSETID",
                    "CAP_KILL",
                    "CAP_SETGID",
                    "CAP_SETUID",
                    "CAP_SETPCAP",
                    "CAP_NET_BIND_SERVICE",
                    "CAP_SYS_MODULE"
                ],
                "effective": [
                    "CAP_CHOWN",
                    "DAC_OVERRIDE",
                    "CAP_FOWNER",
                    "CAP_FSETID",
                    "CAP_KILL",
                    "CAP_SETGID",
                    "CAP_SETUID",
                    "CAP_SETPCAP",
                    "CAP_NET_BIND_SERVICE",
                    "CAP_SYS_MODULE"
                ],
                "inheritable": []
            },
            "ns": {
                "uts": {
                    "inum": 4026533759,
                    "is_host": false
                },
                "ipc": {
                    "inum": 4026533760,
                    "is_host": false
                },
                "mnt": {
                    "inum": 4026534054,
                    "is_host": false
                },
                "pid": {
                    "inum": 4026534055,
                    "is_host": false
                },
                "pid_for_children": {
                    "inum": 4026534055,
                    "is_host": false
                },
                "net": {
                    "inum": 4026533761,
                    "is_host": false
                },
                "time": {
                    "inum": 4026531834,
                    "is_host": true
                },
                "time_for_children": {
                    "inum": 4026531834,
                    "is_host": true
                },
                "cgroup": {
                    "inum": 4026534056,
                    "is_host": false
                },
                "user": {
                    "inum": 4026531837,
                    "is_host": true
                }
            },
            "tid": 3266101,
            "process_credentials": {
                "uid": 0,
                "gid": 0,
                "euid": 0,
                "egid": 0,
                "suid": 0,
                "sgid": 0,
                "fsuid": 0,
                "fsgid": 0,
                "securebits": [],
                "caps": null,
                "user_ns": null
            },
            "binary_properties": null
        },
        "function_name": "do_init_module",
        "args": [
            {
                "module_arg": {
                    "name": "reverse_shell",
                    "signature_ok": null,
                    "tainted": [
                        "TAINT_OUT_OF_TREE_MODULE",
                        "TAINT_UNSIGNED_MODULE"
                    ]
                },
                "label": ""
            }
        ],
        "return": {
            "int_arg": 0,
            "label": ""
        },
        "action": "KPROBE_ACTION_POST",
        "stack_trace": [],
        "policy_name": "kmod-load-monitoring"
    },
    "node_name": "experts-k8s-cs",
    "time": "2024-03-24T07:44:28.761391447Z",
    "aggregation_info": null
}

*/
