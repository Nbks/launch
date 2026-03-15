package launcher

import (
	"fmt"
	"strings"
)

func joinErrors(errs []error) error {
	msgs := make([]string, len(errs))
	for i, err := range errs {
		msgs[i] = err.Error()
	}
	return fmt.Errorf("launcher errors:\n%s", strings.Join(msgs, "\n"))
}

func replaceVar(s, key, val string) string {
	return strings.ReplaceAll(s, key, val)
}
