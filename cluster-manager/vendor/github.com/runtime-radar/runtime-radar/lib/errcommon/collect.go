package errcommon

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

// CollectErrors joins multiple errors into one report and returns it as single error.
// It can output debug level logs.
func CollectErrors(caller string, errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	report := ""
	count := 0

	for i, e := range errs {
		log.Debug().Msgf("%s, error [%d]: %v", caller, i, e)
		report += "[" + strconv.Itoa(i) + "] " + e.Error() + "; "
		count++
	}

	return fmt.Errorf("collected %d errors from %q: %s", count, caller, report)
}
