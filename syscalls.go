// +build !windows,!linux

package readkey

import (
	"golang.org/x/sys/unix"
)

const (
	ioctl_GETATTR = unix.TIOCGETA
	ioctl_SETATTR = unix.TIOCSETA
)
