package clickhouse

import (
	"context"
	"math/rand"
	"time"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/event-processor/api"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model/convert"
	enf_model "github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

func populate(db *gorm.DB, count int) error {
	if count == 0 {
		return nil
	}

	eventDB := &RuntimeEventDatabase{db}

	runtimeEvents := []model.RuntimeEvent{} // some events can be added manually
	runtimeEvents = addRuntimeEvents(runtimeEvents, count)

	return eventDB.Add(context.Background(), &runtimeEvents)
}

func addRuntimeEvents(es []model.RuntimeEvent, count int) []model.RuntimeEvent {
	for i := count; i >= 1; i-- {
		createdAt := time.Now().Add(time.Duration(-i) * time.Minute)
		es = append(es, newRuntimeEvent(createdAt, i%2 == 0)) // every second event won't have any threats
	}

	return es
}

func newRuntimeEvent(createdAt time.Time, withIncident bool) model.RuntimeEvent {
	process := tetragon.Process{
		ExecId:    "a2luZC1jb250cm9sLXBsYW5lOjE1ODkwMjg4MTI5NTg0MDo1NTY2NTc=",
		Pid:       wrapperspb.UInt32(556657),
		Uid:       wrapperspb.UInt32(0),
		Cwd:       "/",
		Binary:    "/usr/bin/nc",
		Arguments: "-vz postgres 5432",
		Flags:     "execve rootcwd clone",
		StartTime: timestamppb.Now(),
		Auid:      wrapperspb.UInt32(4294967295),
		Pod: &tetragon.Pod{
			Namespace: "default",
			Name:      "deathstar-8464cdd4d9-7slzs",
			Container: &tetragon.Container{
				Id:   "containerd://beb265602c5db9b2b8a3cf603397642fe7c3c3d7c2092b4353b019e1779f2810",
				Name: "deathstar",
				Image: &tetragon.Image{
					Id:   "docker.io/cilium/starwars@sha256:f92c8cd25372bac56f55111469fe9862bf682385a4227645f5af155eee7f58d9",
					Name: "docker.io/cilium/starwars:latest",
				},
				StartTime: timestamppb.Now(),
				Pid:       wrapperspb.UInt32(51),
			},
			PodLabels: map[string]string{
				"app.kubernetes.io/name": "deathstar",
				"class":                  "deathstar",
				"org":                    "empire",
				"pod-template-hash":      "8464cdd4d9",
			},
			Workload: "workload",
		},
		Docker:       "beb265602c5db9b2b8a3cf603397642",
		ParentExecId: "a2luZC1jb250cm9sLXBsYW5lOjE1ODc1ODY3MDM4MDY2NDo1NTU5NjY=",
		Refcnt:       342423,
		Cap: &tetragon.Capabilities{
			Permitted: []tetragon.CapabilitiesType{
				tetragon.CapabilitiesType_CAP_CHOWN,
				tetragon.CapabilitiesType_CAP_FOWNER,
				tetragon.CapabilitiesType_CAP_FSETID,
				tetragon.CapabilitiesType_CAP_KILL,
				tetragon.CapabilitiesType_CAP_SETGID,
				tetragon.CapabilitiesType_CAP_SETUID,
				tetragon.CapabilitiesType_CAP_SETPCAP,
				tetragon.CapabilitiesType_CAP_NET_BIND_SERVICE,
				tetragon.CapabilitiesType_CAP_NET_RAW,
				tetragon.CapabilitiesType_CAP_SYS_CHROOT,
				tetragon.CapabilitiesType_CAP_MKNOD,
				tetragon.CapabilitiesType_CAP_AUDIT_WRITE,
				tetragon.CapabilitiesType_CAP_SETFCAP,
			},
			Effective: []tetragon.CapabilitiesType{
				tetragon.CapabilitiesType_CAP_CHOWN,
				tetragon.CapabilitiesType_CAP_FOWNER,
				tetragon.CapabilitiesType_CAP_FSETID,
				tetragon.CapabilitiesType_CAP_KILL,
				tetragon.CapabilitiesType_CAP_SETGID,
				tetragon.CapabilitiesType_CAP_SETUID,
				tetragon.CapabilitiesType_CAP_SETPCAP,
				tetragon.CapabilitiesType_CAP_NET_BIND_SERVICE,
				tetragon.CapabilitiesType_CAP_NET_RAW,
				tetragon.CapabilitiesType_CAP_SYS_CHROOT,
				tetragon.CapabilitiesType_CAP_MKNOD,
				tetragon.CapabilitiesType_CAP_AUDIT_WRITE,
				tetragon.CapabilitiesType_CAP_SETFCAP,
			},
			Inheritable: []tetragon.CapabilitiesType{
				tetragon.CapabilitiesType_CAP_CHOWN,
				tetragon.CapabilitiesType_CAP_FOWNER,
				tetragon.CapabilitiesType_CAP_FSETID,
				tetragon.CapabilitiesType_CAP_KILL,
				tetragon.CapabilitiesType_CAP_SETGID,
				tetragon.CapabilitiesType_CAP_SETUID,
				tetragon.CapabilitiesType_CAP_SETPCAP,
				tetragon.CapabilitiesType_CAP_NET_BIND_SERVICE,
				tetragon.CapabilitiesType_CAP_NET_RAW,
				tetragon.CapabilitiesType_CAP_SYS_CHROOT,
				tetragon.CapabilitiesType_CAP_MKNOD,
				tetragon.CapabilitiesType_CAP_AUDIT_WRITE,
				tetragon.CapabilitiesType_CAP_SETFCAP,
			},
		},
		Ns: &tetragon.Namespaces{
			Uts: &tetragon.Namespace{
				Inum: 4026532823,
			},
			Ipc: &tetragon.Namespace{
				Inum: 4026532824,
			},
			Mnt: &tetragon.Namespace{
				Inum: 4026532830,
			},
			Pid: &tetragon.Namespace{
				Inum: 4026532831,
			},
			PidForChildren: &tetragon.Namespace{
				Inum: 4026532831,
			},
			Net: &tetragon.Namespace{
				Inum: 4026532699,
			},
			Cgroup: &tetragon.Namespace{
				Inum:   4026531835,
				IsHost: false,
			},
			User: &tetragon.Namespace{
				Inum:   4026531837,
				IsHost: false,
			},
		},
		Tid: wrapperspb.UInt32(556657),
	}

	event := &api.RuntimeEvent{
		Id:              uuid.NewString(),
		TetragonVersion: "1.0.0",
		Event: &tetragon.GetEventsResponse{
			Event: &tetragon.GetEventsResponse_ProcessExec{
				ProcessExec: &tetragon.ProcessExec{
					Process: &process,
					Parent:  &process,
				},
			},
			NodeName: "kind-control-plane",
			Time:     timestamppb.New(createdAt),
		},
	}

	if withIncident {
		s := enf_model.Severity(rand.Intn(5)).String() // nolint:gosec

		event.Threats = []*api.Threat{
			{
				Detector: &api.Detector{
					Id:          "1",
					Name:        "test_detector",
					Version:     1,
					Description: "dummy",
					Author:      "CS Team",
				},
				Severity: s,
			},
		}

		event.IsIncident = true
		event.IncidentSeverity = s
		event.BlockBy = []string{"74875022-32b9-4089-a665-63c230ed4d63", "74875022-32b9-4089-a665-63c230ed4d68"}
		event.NotifyBy = []string{"74875022-32b9-4089-a665-63c230ed4d63", "74875022-32b9-4089-a665-63c230ed4d68"}
	}

	m, err := convert.RuntimeEventFromProto(event)
	if err != nil {
		log.Fatal().Err(err).Msg("can't convert event to model")
	}

	return m
}
