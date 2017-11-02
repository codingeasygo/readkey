package readkey

/*
void rk_init();
void rk_release();
int rk_read(char *buf, int len);
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

var opened = false

func Open() {
	if opened {
		return
	}
	C.rk_init()
	opened = true
}

var OnKey = func(key []byte) {

}

// //export onKey
// func onKey(buf *C.char, l C.int) {
// 	OnKey(C.GoBytes(unsafe.Pointer(buf), l))
// }

var buf [100]C.char

func Read() (key []byte, err error) {
	readed := C.rk_read(&buf[0], 100)
	if readed < 1 {
		err = fmt.Errorf("return code(%v)", readed)
	} else {
		key = C.GoBytes(unsafe.Pointer(&buf[0]), readed)
	}
	return
}

func Close() {
	C.rk_release()
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func GetSize() (w, h int) {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	w, h = int(ws.Col), int(ws.Row)
	return
}

func SetSize(fd uintptr, w, h int) (err error) {
	ws := &winsize{
		Col: uint16(w),
		Row: uint16(h),
	}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd,
		uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		err = fmt.Errorf("%v -%v", retCode, errno)
	}
	return
}
