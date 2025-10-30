package helm

import (
	"fmt"
	"net/url"
	"strings"
)

// ConvertRegistryToOCI converts registryAddress to desired chart path.
// Prepends "oci://" if needed and forces "oci" scheme.
func ConvertRegistryToOCI(registryAddress string) (*url.URL, error) {
	if !strings.Contains(registryAddress, "://") {
		registryAddress = "oci://" + strings.TrimLeft(registryAddress, "/")
	}

	chartURL, err := url.ParseRequestURI(registryAddress)
	if err != nil {
		return nil, err
	}

	chartURL.Scheme = "oci"

	return chartURL, nil
}

// UpgradeCmd returns command to upgrade CS release
// If values is nil it would use values.yaml file
func UpgradeCmd(registryAddress, chartVersion, namespace string, values *Values) (string, error) {
	if registryAddress == "" {
		registryAddress = DefaultRegistry
	}

	chartURL, err := ConvertRegistryToOCI(registryAddress)
	if err != nil {
		return "", err
	}

	cmd := fmt.Sprintf("helm upgrade --install runtime-radar %s --version %s --namespace %s --create-namespace",
		chartURL.JoinPath("runtime-radar"),
		chartVersion,
		namespace,
	)

	if values == nil {
		return fmt.Sprintf("%s --values values.yaml", cmd), nil
	}

	args, err := values.ToHelmArgs()
	if err != nil {
		return "", err
	}

	lines := append([]string{cmd}, args...)

	cmd = strings.Join(lines, " \\\n")

	return cmd, nil
}

func UninstallCmd(namespace string) string {
	return fmt.Sprintf("helm uninstall runtime-radar --namespace %s", namespace)
}
