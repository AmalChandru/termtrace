//go:build !linux

package record

import "os"

func disablePTYEcho(_ *os.File) error {
	return nil
}
