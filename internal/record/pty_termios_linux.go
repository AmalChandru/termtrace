//go:build linux

package record

import (
	"os"
	"syscall"
	"unsafe"
)

func disablePTYEcho(f *os.File) error {
	fd := int(f.Fd())
	termios, err := ioctlGetTermios(fd, syscall.TCGETS)
	if err != nil {
		return err
	}
	termios.Lflag &^= syscall.ECHO
	return ioctlSetTermios(fd, syscall.TCSETS, termios)
}

func ioctlGetTermios(fd int, req uint) (*syscall.Termios, error) {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(req),
		uintptr(unsafe.Pointer(&termios)),
		0,
		0,
		0,
	)
	if errno != 0 {
		return nil, errno
	}
	return &termios, nil
}

func ioctlSetTermios(fd int, req uint, termios *syscall.Termios) error {
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(req),
		uintptr(unsafe.Pointer(termios)),
		0,
		0,
		0,
	)
	if errno != 0 {
		return errno
	}
	return nil
}
