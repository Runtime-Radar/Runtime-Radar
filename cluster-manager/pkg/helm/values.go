package helm

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"sigs.k8s.io/yaml"
)

const (
	DefaultUser              = "runtime-radar"
	DefaultNamespace         = "runtime-radar"
	DefaultRetentionInterval = time.Hour * 8460
	DefaultRegistry          = "ghcr.io/runtime-radar"
)

var (
	errInputNotStruct = errors.New("input must be a struct")
	errMarshal        = errors.New("failed to marshal field")
)

// Values represents helm values returned as a complete install command or yaml.
// `sigs.k8s.io/yaml` library ignores `yaml` tag and uses only `json`. So only
// `json` tag is used here.
// `struct` fields are followed by `json:"...,omitzero"` because `omitempty`
// does not check for struct's zero value.
type Values struct {
	Global struct {
		CSVersion      string `json:"csVersion,omitempty"`
		IsChildCluster bool   `json:"isChildCluster"`
		OwnCSURL       string `json:"ownCsUrl,omitempty"`
		CentralCSURL   string `json:"centralCsUrl,omitempty"`

		ImageRegistry   string `json:"imageRegistry,omitempty"`
		ImageShortNames bool   `json:"imageShortNames,omitempty"`

		Keys struct {
			Encryption            string `json:"encryption,omitempty"`
			Token                 string `json:"token,omitempty"`
			PublicAccessTokenSalt string `json:"publicAccessTokenSalt,omitempty"`
		} `json:"keys,omitzero"`

		Postgresql TLSGlobal `json:"postgresql"`
		Redis      TLSGlobal `json:"redis"`
		Clickhouse TLSGlobal `json:"clickhouse"`
		Grafana    TLSGlobal `json:"grafana,omitzero"`
		Loki       TLSGlobal `json:"loki,omitzero"`
	} `json:"global,omitzero"`

	TLS struct {
		Verify  bool   `json:"verify"`
		CertCA  string `json:"certCA,omitempty"`
		Cert    string `json:"cert,omitempty"`
		CertKey string `json:"certKey,omitempty"`
	} `json:"tls,omitzero"`

	AuthAPI struct {
		Administrator struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		} `json:"administrator,omitzero"`
	} `json:"auth-center,omitzero"`

	ImagePullSecret struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"imagePullSecret,omitzero"`

	Postgresql struct {
		Deploy       bool   `json:"deploy"`
		ExternalHost string `json:"externalHost,omitempty"`

		Auth struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
			Database string `json:"database,omitempty"`
		} `json:"auth,omitzero"`

		Persistence struct {
			Enabled      bool   `json:"enabled"`
			StorageClass string `json:"storageClass,omitempty"`
		} `json:"persistence"`

		TLS struct {
			CertCA  string `json:"certCA,omitempty"`
			Cert    string `json:"cert,omitempty"`
			CertKey string `json:"certKey,omitempty"`
		} `json:"tls,omitzero"`
	} `json:"postgresql"`

	Redis struct {
		Deploy       bool   `json:"deploy"`
		ExternalHost string `json:"externalHost,omitempty"`

		Auth struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		} `json:"auth,omitzero"`

		Persistence struct {
			Enabled      bool   `json:"enabled"`
			StorageClass string `json:"storageClass,omitempty"`
		} `json:"persistence"`

		TLS struct {
			CertCA  string `json:"certCA,omitempty"`
			Cert    string `json:"cert,omitempty"`
			CertKey string `json:"certKey,omitempty"`
		} `json:"tls,omitzero"`
	} `json:"redis"`

	Rabbitmq struct {
		Deploy       bool   `json:"deploy"`
		ExternalHost string `json:"externalHost,omitempty"`

		Auth struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		} `json:"auth,omitzero"`

		Persistence struct {
			Enabled      bool   `json:"enabled"`
			StorageClass string `json:"storageClass,omitempty"`
		} `json:"persistence"`
	} `json:"rabbitmq"`

	Clickhouse struct {
		Deploy       bool   `json:"deploy"`
		ExternalHost string `json:"externalHost,omitempty"`

		Auth struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
			Database string `json:"database,omitempty"`
		} `json:"auth,omitzero"`

		Persistence struct {
			Enabled      bool   `json:"enabled"`
			StorageClass string `json:"storageClass,omitempty"`
		} `json:"persistence"`

		TLS struct {
			CertCA  string `json:"certCA,omitempty"`
			Cert    string `json:"cert,omitempty"`
			CertKey string `json:"certKey,omitempty"`
		} `json:"tls,omitzero"`
	} `json:"clickhouse"`

	Grafana struct {
		Deploy       bool   `json:"deploy"`
		ExternalHost string `json:"externalHost,omitempty"`

		Admin struct {
			User     string `json:"user,omitempty"`
			Password string `json:"password,omitempty"`
		} `json:"admin,omitzero"`

		Persistence struct {
			Enabled      bool   `json:"enabled"`
			StorageClass string `json:"storageClass,omitempty"`
		} `json:"persistence"`

		TLS struct {
			CertCA  string `json:"certCA,omitempty"`
			Cert    string `json:"cert,omitempty"`
			CertKey string `json:"certKey,omitempty"`
		} `json:"tls,omitzero"`
	} `json:"grafana,omitzero"`

	Loki struct {
		Deploy       bool   `json:"deploy"`
		ExternalHost string `json:"externalHost,omitempty"`

		TenantID string `json:"tenant_id,omitempty"`

		SingleBinary struct {
			Persistence struct {
				Enabled      bool   `json:"enabled"`
				StorageClass string `json:"storageClass,omitempty"`
			} `json:"persistence"`
		}

		TLS struct {
			CertCA  string `json:"certCA,omitempty"`
			Cert    string `json:"cert,omitempty"`
			CertKey string `json:"certKey,omitempty"`
		} `json:"tls,omitzero"`
	} `json:"loki,omitzero"`

	ReverseProxy struct {
		Ingress struct {
			Enabled    bool   `json:"enabled,omitempty"`
			Class      string `json:"class,omitempty"`
			Hostname   string `json:"hostname,omitempty"`
			SecretName string `json:"secretName,omitempty"`
			TLS        struct {
				CertCA  string `json:"certCA,omitempty"`
				Cert    string `json:"cert,omitempty"`
				CertKey string `json:"certKey,omitempty"`
			} `json:"tls,omitzero"`
		} `json:"ingress,omitzero"`

		Service struct {
			Type      string `json:"type,omitempty"`
			NodePorts struct {
				HTTP string `json:"http,omitempty"`
			} `json:"nodePorts,omitzero"`
		} `json:"service,omitzero"`
	} `json:"reverse-proxy,omitzero"`

	Notifier struct {
		OverwriteEnv []Env `json:"overwriteEnv,omitempty"`
	} `json:"notifier,omitzero"`

	CSManager struct {
		RegistrationToken string `json:"registrationToken,omitempty"`
	} `json:"cs-manager,omitzero"`
}

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TLSGlobal struct {
	TLS struct {
		Enabled bool `json:"enabled"`
		Verify  bool `json:"verify"`
	} `json:"tls"`
}

// buildHelmArgs recursively converts a struct's fields to Helm command-line arguments.
// It traverses the struct and generates the appropriate --set or --set-string arguments
// based on field types.
//
// Each field is processed according to its kind:
//   - Strings use --set-string
//   - Bool, numeric types use --set
//   - Arrays/Slices are JSON marshaled and use --set
//   - Structs are processed recursively
//
// Fields with `json:"-"` are skipped.
// Fields with `json:",omitempty"` or `json:",omitzero"` are skipped if they contain zero values.
//
// Returns an error if JSON marshaling fails for array/slice fields.
func buildHelmArgs(v any, prefix string) ([]string, error) {
	var res []string

	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, errInputNotStruct
	}

	t := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("json")
		if tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		fieldName := parts[0]
		if fieldName == "" {
			fieldName = fieldType.Name
		}
		hasOmit := (slices.Index(parts, "omitempty") != -1) || (slices.Index(parts, "omitzero") != -1)

		// skip field on `omitempty` and `omitzero`
		// theoretically `field.IsZero()` can possibly panic
		// but currently it should never happen
		if hasOmit && field.IsZero() {
			continue
		}

		currentPrefix := fieldName
		if prefix != "" {
			currentPrefix = prefix + "." + fieldName
		}

		switch field.Kind() {
		case reflect.String:
			res = append(res, fmt.Sprintf("--set-string '%s=%s'", currentPrefix, field.String()))
		case reflect.Bool:
			res = append(res, fmt.Sprintf("--set '%s=%t'", currentPrefix, field.Bool()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			res = append(res, fmt.Sprintf("--set '%s=%d'", currentPrefix, field.Int()))
		case reflect.Float32, reflect.Float64:
			res = append(res, fmt.Sprintf("--set '%s=%f'", currentPrefix, field.Float()))
		case reflect.Array, reflect.Slice:
			if hasOmit && field.Len() == 0 {
				continue
			}
			b, err := json.Marshal(field.Interface())
			if err != nil {
				return nil, fmt.Errorf("%w: %w", errMarshal, err)
			}
			res = append(res, fmt.Sprintf("--set '%s=%s'", currentPrefix, string(b)))
		case reflect.Struct:
			nestedRes, err := buildHelmArgs(field.Interface(), currentPrefix)
			if err != nil {
				return nil, err
			}
			res = append(res, nestedRes...)
		}
	}
	return res, nil
}

func (v *Values) ToYAML() (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

func (v *Values) ToHelmArgs() ([]string, error) {
	return buildHelmArgs(v, "")
}
