package prompt

import (
	"bytes"
	"syscall"
	"unsafe"

	"github.com/pkg/term/termios"
)

type VT100Parser struct {
	fd          int
	origTermios syscall.Termios
}

// winsize is winsize struct got from the ioctl(2) system call.
type ioctlWinsize struct {
	Row uint16
	Col uint16
	X   uint16 // pixel value
	Y   uint16 // pixel value
}

func (t *VT100Parser) Setup() error {
	return t.setRawMode()
}

func (t *VT100Parser) setRawMode() error {
	var n syscall.Termios
	if err := termios.Tcgetattr(uintptr(t.fd), &t.origTermios); err != nil {
		return err
	}
	n = t.origTermios
	// "&^=" used like: https://play.golang.org/p/8eJw3JxS4O
	n.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	n.Cc[syscall.VMIN] = 1
	n.Cc[syscall.VTIME] = 0
	termios.Tcsetattr(uintptr(t.fd), termios.TCSANOW, (*syscall.Termios)(&n))
	return nil
}

func (t *VT100Parser) TearDown() error {
	return termios.Tcsetattr(uintptr(t.fd), termios.TCSANOW, &t.origTermios)
}

func (t *VT100Parser) GetASCIICode(b []byte) *ASCIICode {
	for _, k := range asciiSequences {
		if bytes.Compare(k.ASCIICode, b) == 0 {
			return k
		}
	}
	return nil
}

// GetWinSize returns winsize struct which is the response of ioctl(2).
func (t *VT100Parser) GetWinSize() (row, col uint16) {
	ws := &ioctlWinsize{}
	retCode, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(t.fd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return ws.Row, ws.Col
}

var asciiSequences []*ASCIICode = []*ASCIICode{
	&ASCIICode{Key: Escape, ASCIICode: []byte{0x1b}},

	&ASCIICode{Key: ControlSpace, ASCIICode: []byte{0x00}},
	&ASCIICode{Key: ControlA, ASCIICode: []byte{0x1}},
	&ASCIICode{Key: ControlB, ASCIICode: []byte{0x2}},
	&ASCIICode{Key: ControlC, ASCIICode: []byte{0x3}},
	&ASCIICode{Key: ControlD, ASCIICode: []byte{0x4}},
	&ASCIICode{Key: ControlE, ASCIICode: []byte{0x5}},
	&ASCIICode{Key: ControlF, ASCIICode: []byte{0x6}},
	&ASCIICode{Key: ControlG, ASCIICode: []byte{0x7}},
	&ASCIICode{Key: ControlH, ASCIICode: []byte{0x8}},
	&ASCIICode{Key: ControlI, ASCIICode: []byte{0x9}},
	&ASCIICode{Key: ControlJ, ASCIICode: []byte{0xa}},
	&ASCIICode{Key: ControlK, ASCIICode: []byte{0xb}},
	&ASCIICode{Key: ControlL, ASCIICode: []byte{0xc}},
	&ASCIICode{Key: ControlM, ASCIICode: []byte{0xd}},
	&ASCIICode{Key: ControlN, ASCIICode: []byte{0xe}},
	&ASCIICode{Key: ControlO, ASCIICode: []byte{0xf}},
	&ASCIICode{Key: ControlP, ASCIICode: []byte{0x10}},
	&ASCIICode{Key: ControlQ, ASCIICode: []byte{0x11}},
	&ASCIICode{Key: ControlR, ASCIICode: []byte{0x12}},
	&ASCIICode{Key: ControlS, ASCIICode: []byte{0x13}},
	&ASCIICode{Key: ControlT, ASCIICode: []byte{0x14}},
	&ASCIICode{Key: ControlU, ASCIICode: []byte{0x15}},
	&ASCIICode{Key: ControlV, ASCIICode: []byte{0x16}},
	&ASCIICode{Key: ControlW, ASCIICode: []byte{0x17}},
	&ASCIICode{Key: ControlX, ASCIICode: []byte{0x18}},
	&ASCIICode{Key: ControlY, ASCIICode: []byte{0x19}},
	&ASCIICode{Key: ControlZ, ASCIICode: []byte{0x1a}},

	&ASCIICode{Key: ControlBackslash, ASCIICode: []byte{0x1c}},
	&ASCIICode{Key: ControlSquareClose, ASCIICode: []byte{0x1d}},
	&ASCIICode{Key: ControlCircumflex, ASCIICode: []byte{0x1e}},
	&ASCIICode{Key: ControlUnderscore, ASCIICode: []byte{0x1f}},
	&ASCIICode{Key: Backspace, ASCIICode: []byte{0x7f}},

	&ASCIICode{Key: Up, ASCIICode: []byte{0x1b, 0x5b, 0x41}},
	&ASCIICode{Key: Down, ASCIICode: []byte{0x1b, 0x5b, 0x42}},
	&ASCIICode{Key: Right, ASCIICode: []byte{0x1b, 0x5b, 0x43}},
	&ASCIICode{Key: Left, ASCIICode: []byte{0x1b, 0x5b, 0x44}},
	&ASCIICode{Key: Home, ASCIICode: []byte{0x1b, 0x5b, 0x48}},
	&ASCIICode{Key: Home, ASCIICode: []byte{0x1b, 0x4f, 0x48}},
	&ASCIICode{Key: End, ASCIICode: []byte{0x1b, 0x5b, 0x70}},
	&ASCIICode{Key: End, ASCIICode: []byte{0x1b, 0x4f, 0x70}},

	&ASCIICode{Key: Delete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	&ASCIICode{Key: ShiftDelete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x3b, 0x02, 0x7e}},
	&ASCIICode{Key: ControlDelete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x3b, 0x05, 0x7e}},
	&ASCIICode{Key: Home, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x7e}},
	&ASCIICode{Key: End, ASCIICode: []byte{0x1b, 0x5b, 0x04, 0x7e}},
	&ASCIICode{Key: PageUp, ASCIICode: []byte{0x1b, 0x5b, 0x05, 0x7e}},
	&ASCIICode{Key: PageDown, ASCIICode: []byte{0x1b, 0x5b, 0x06, 0x7e}},
	&ASCIICode{Key: Home, ASCIICode: []byte{0x1b, 0x5b, 0x07, 0x7e}},
	&ASCIICode{Key: End, ASCIICode: []byte{0x1b, 0x5b, 0x09, 0x7e}},
	&ASCIICode{Key: BackTab, ASCIICode: []byte{0x1b, 0x5b, 0x5a}},
	&ASCIICode{Key: Insert, ASCIICode: []byte{0x1b, 0x5b, 0x02, 0x7e}},

	&ASCIICode{Key: F1, ASCIICode: []byte{0x1b, 0x4f, 0x50}},
	&ASCIICode{Key: F2, ASCIICode: []byte{0x1b, 0x4f, 0x51}},
	&ASCIICode{Key: F3, ASCIICode: []byte{0x1b, 0x4f, 0x52}},
	&ASCIICode{Key: F4, ASCIICode: []byte{0x1b, 0x4f, 0x53}},

	&ASCIICode{Key: F1, ASCIICode: []byte{0x1b, 0x4f, 0x50, 0x41}}, // Linux console
	&ASCIICode{Key: F2, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x42}}, // Linux console
	&ASCIICode{Key: F3, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x43}}, // Linux console
	&ASCIICode{Key: F4, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x44}}, // Linux console
	&ASCIICode{Key: F5, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x45}}, // Linux console

	&ASCIICode{Key: F1, ASCIICode: []byte{0x1b, 0x5b, 0x11, 0x7e}}, // rxvt-unicode
	&ASCIICode{Key: F2, ASCIICode: []byte{0x1b, 0x5b, 0x12, 0x7e}}, // rxvt-unicode
	&ASCIICode{Key: F3, ASCIICode: []byte{0x1b, 0x5b, 0x13, 0x7e}}, // rxvt-unicode
	&ASCIICode{Key: F4, ASCIICode: []byte{0x1b, 0x5b, 0x14, 0x7e}}, // rxvt-unicode

	&ASCIICode{Key: F5, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x35, 0x7e}},
	&ASCIICode{Key: F6, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x37, 0x7e}},
	&ASCIICode{Key: F7, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x38, 0x7e}},
	&ASCIICode{Key: F8, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x39, 0x7e}},
	&ASCIICode{Key: F9, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x30, 0x7e}},
	&ASCIICode{Key: F10, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x31, 0x7e}},
	&ASCIICode{Key: F11, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x32, 0x7e}},
	&ASCIICode{Key: F12, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x34, 0x7e, 0x8}},
	&ASCIICode{Key: F13, ASCIICode: []byte{0x1b, 0x5b, 0x25, 0x7e}},
	&ASCIICode{Key: F14, ASCIICode: []byte{0x1b, 0x5b, 0x26, 0x7e}},
	&ASCIICode{Key: F15, ASCIICode: []byte{0x1b, 0x5b, 0x28, 0x7e}},
	&ASCIICode{Key: F16, ASCIICode: []byte{0x1b, 0x5b, 0x29, 0x7e}},
	&ASCIICode{Key: F17, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x7e}},
	&ASCIICode{Key: F18, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x7e}},
	&ASCIICode{Key: F19, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	&ASCIICode{Key: F20, ASCIICode: []byte{0x1b, 0x5b, 0x34, 0x7e}},

	// Xterm
	&ASCIICode{Key: F13, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x3b, 0x02, 0x50}},
	&ASCIICode{Key: F14, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x3b, 0x02, 0x51}},
	// &ASCIICode{Key: F15, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x3b, 0x02, 0x52}},  // Conflicts with CPR response
	&ASCIICode{Key: F16, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x3b, 0x02, 0x52}},
	&ASCIICode{Key: F17, ASCIICode: []byte{0x1b, 0x5b, 0x15, 0x3b, 0x02, 0x7e}},
	&ASCIICode{Key: F18, ASCIICode: []byte{0x1b, 0x5b, 0x17, 0x3b, 0x02, 0x7e}},
	&ASCIICode{Key: F19, ASCIICode: []byte{0x1b, 0x5b, 0x18, 0x3b, 0x02, 0x7e}},
	&ASCIICode{Key: F20, ASCIICode: []byte{0x1b, 0x5b, 0x19, 0x3b, 0x02, 0x7e}},
	&ASCIICode{Key: F21, ASCIICode: []byte{0x1b, 0x5b, 0x20, 0x3b, 0x02, 0x7e}},
	&ASCIICode{Key: F22, ASCIICode: []byte{0x1b, 0x5b, 0x21, 0x3b, 0x02, 0x7e}},
	&ASCIICode{Key: F23, ASCIICode: []byte{0x1b, 0x5b, 0x23, 0x3b, 0x02, 0x7e}},
	&ASCIICode{Key: F24, ASCIICode: []byte{0x1b, 0x5b, 0x24, 0x3b, 0x02, 0x7e}},

	&ASCIICode{Key: ControlUp, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x3b, 0x5a}},
	&ASCIICode{Key: ControlDown, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x3b, 0x5b}},
	&ASCIICode{Key: ControlRight, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x3b, 0x5c}},
	&ASCIICode{Key: ControlLeft, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x3b, 0x5d}},

	&ASCIICode{Key: ShiftUp, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x2a}},
	&ASCIICode{Key: ShiftDown, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x2b}},
	&ASCIICode{Key: ShiftRight, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x2c}},
	&ASCIICode{Key: ShiftLeft, ASCIICode: []byte{0x1b, 0x5b, 0x01, 0x2d}},

	// Tmux sends following keystrokes when control+arrow is pressed, but for
	// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	// it a normal arrow press, because that's more important.
	&ASCIICode{Key: Up, ASCIICode: []byte{0x1b, 0x4f, 0x41}},
	&ASCIICode{Key: Down, ASCIICode: []byte{0x1b, 0x4f, 0x42}},
	&ASCIICode{Key: Right, ASCIICode: []byte{0x1b, 0x4f, 0x43}},
	&ASCIICode{Key: Left, ASCIICode: []byte{0x1b, 0x4f, 0x44}},

	&ASCIICode{Key: ControlUp, ASCIICode: []byte{0x1b, 0x5b, 0x05, 0x41}},
	&ASCIICode{Key: ControlDown, ASCIICode: []byte{0x1b, 0x5b, 0x05, 0x42}},
	&ASCIICode{Key: ControlRight, ASCIICode: []byte{0x1b, 0x5b, 0x05, 0x43}},
	&ASCIICode{Key: ControlLeft, ASCIICode: []byte{0x1b, 0x5b, 0x05, 0x44}},

	&ASCIICode{Key: ControlRight, ASCIICode: []byte{0x1b, 0x5b, 0x4f, 0x63}}, // rxvt
	&ASCIICode{Key: ControlLeft, ASCIICode: []byte{0x1b, 0x5b, 0x4f, 0x64}},  // rxvt

	&ASCIICode{Key: Ignore, ASCIICode: []byte{0x1b, 0x5b, 0x45}}, // Xterm
	&ASCIICode{Key: Ignore, ASCIICode: []byte{0x1b, 0x5b, 0x46}}, // Linux console
}

func NewVT100Parser() *VT100Parser {
	return &VT100Parser{
		fd: syscall.Stdin,
	}
}