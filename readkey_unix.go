// +build !windows

package readkey

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type (
	input_event struct {
		data []byte
		err  error
	}
)

var (
	out  *os.File
	infd int
	// termbox inner state
	origTios unix.Termios

	sigio   = make(chan os.Signal, 1)
	quit    = make(chan int)
	opened  bool
	running bool
)

func fcntl(cmd int, arg int) error {
	_, _, e := syscall.Syscall(unix.SYS_FCNTL, uintptr(infd), uintptr(cmd), uintptr(arg))
	if e != 0 {
		return e
	}

	return nil
}

func ioctl(cmd int, termios *unix.Termios) error {
	r, _, e := syscall.Syscall(unix.SYS_IOCTL, out.Fd(), uintptr(cmd), uintptr(unsafe.Pointer(termios)))
	if r != 0 {
		return os.NewSyscallError("SYS_IOCTL", e)
	}
	return nil
}

func Open() (err error) {
	if opened {
		panic("already opened")
	}
	out, err = os.OpenFile("/dev/tty", unix.O_WRONLY, 0)
	if err != nil {
		return
	}
	infd, err = syscall.Open("/dev/tty", unix.O_RDONLY, 0)
	if err != nil {
		return
	}
	signal.Notify(sigio, unix.SIGIO)
	err = fcntl(unix.F_SETFL, unix.O_ASYNC|unix.O_NONBLOCK)
	if err != nil {
		return
	}
	err = fcntl(unix.F_SETOWN, unix.Getpid())
	if runtime.GOOS != "darwin" && err != nil {
		return
	}

	err = ioctl(ioctl_GETATTR, &origTios)
	if err != nil {
		return
	}

	tios := origTios
	tios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK |
		unix.ISTRIP | unix.INLCR | unix.IGNCR |
		unix.ICRNL | unix.IXON
	tios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON |
		unix.ISIG | unix.IEXTEN
	tios.Cflag &^= unix.CSIZE | unix.PARENB
	tios.Cflag |= unix.CS8
	tios.Cc[unix.VMIN] = 1
	tios.Cc[unix.VTIME] = 0

	err = ioctl(ioctl_SETATTR, &tios)
	opened = true
	return
}

func Loop(onkey func([]byte)) {
	if running {
		panic("already running")
	}
	running = true
	readed := 0
	var err error
	buf := make([]byte, 128)
	for {
		select {
		case <-sigio:
			for {
				readed, err = syscall.Read(infd, buf)
				if err == unix.EAGAIN || err == unix.EWOULDBLOCK {
					break
				}
				onkey(buf[:readed])
			}
		case <-quit:
			return
		}
	}
}

func Close() {
	quit <- 1
	ioctl(ioctl_SETATTR, &origTios)
	out.Close()
	unix.Close(infd)
	running = false
	opened = false
}

type Buffer struct {
	buf chan []byte
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		buf: make(chan []byte, size),
	}
}

func (b *Buffer) OnKey(key []byte) {
	b.buf <- key
}

func (b *Buffer) ReadKey() (key []byte) {
	key = <-b.buf
	return
}

var SharedBuffer = NewBuffer(10240)

func ReadKey() (key []byte, err error) {
	if !running {
		if !opened {
			err = Open()
			if err != nil {
				return
			}
		}
		go Loop(SharedBuffer.OnKey)
	}
	key = SharedBuffer.ReadKey()
	return
}
