// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux

package zcall

import (
	"unsafe"

	"code.hybscloud.com/zcall/internal"
)

type ifreq struct {
	name [IFNAMSIZ]byte
	data [24]byte
}

type sockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

type nlMsghdr struct {
	Len   uint32
	Type  uint16
	Flags uint16
	Seq   uint32
	Pid   uint32
}

type rtGenmsg struct {
	Family uint8
}

type ifInfomsg struct {
	Family uint8
	_      uint8
	Type   uint16
	Index  int32
	Flags  uint32
	Change uint32
}

type rtAttr struct {
	Len  uint16
	Type uint16
}

// Link is a raw netlink record for a network link (ifinfomsg + rtattrs).
type Link struct {
	Index        int
	MTU          int
	Name         string
	HardwareAddr []byte
	Flags        uint32
}

const (
	nlmsgAlignTo = 4
	rtaAlignTo   = 4
)

// Syscall4 executes a syscall with up to 4 arguments.
func Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr) {
	return internal.RawSyscall4(num, a1, a2, a3, a4)
}

// Syscall6 executes a syscall with up to 6 arguments.
func Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr) {
	return internal.RawSyscall6(num, a1, a2, a3, a4, a5, a6)
}

// Close closes a file descriptor.
func Close(fd uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_CLOSE, fd, 0, 0, 0)
	return
}

// Read reads from a file descriptor into buf.
func Read(fd uintptr, buf []byte) (n uintptr, errno uintptr) {
	if len(buf) == 0 {
		return 0, 0
	}
	return Syscall4(SYS_READ, fd, uintptr(noescape(unsafe.Pointer(&buf[0]))), uintptr(len(buf)), 0)
}

// Write writes buf to a file descriptor.
func Write(fd uintptr, buf []byte) (n uintptr, errno uintptr) {
	if len(buf) == 0 {
		return 0, 0
	}
	return Syscall4(SYS_WRITE, fd, uintptr(noescape(unsafe.Pointer(&buf[0]))), uintptr(len(buf)), 0)
}

// Socket creates a socket.
func Socket(domain, typ, protocol uintptr) (fd uintptr, errno uintptr) {
	return Syscall4(SYS_SOCKET, domain, typ, protocol, 0)
}

// Ioctl executes an ioctl syscall.
func Ioctl(fd, req uintptr, arg unsafe.Pointer) (errno uintptr) {
	_, errno = Syscall4(SYS_IOCTL, fd, req, uintptr(noescape(arg)), 0)
	return
}

// Links returns the current Linux network link table.
func Links() ([]Link, error) {
	fd, pid, err := openRouteNetlink()
	return linksFromRouteNetlink(fd, pid, err)
}

func linksFromRouteNetlink(fd uintptr, pid uint32, openErr error) ([]Link, error) {
	if openErr != nil {
		return nil, openErr
	}
	defer Close(fd)

	const seq = 1
	if err := sendNetlinkGetLinkDump(fd, seq); err != nil {
		return nil, err
	}
	return recvLinks(fd, pid, seq)
}

func findLinkByName(name string, links []Link, err error) (*Link, error) {
	if err != nil {
		return nil, err
	}
	for i := range links {
		if links[i].Name == name {
			return &links[i], nil
		}
	}
	return nil, Errno(ENODEV)
}

func findLinkByIndex(index int, links []Link, err error) (*Link, error) {
	if err != nil {
		return nil, err
	}
	for i := range links {
		if links[i].Index == index {
			return &links[i], nil
		}
	}
	return nil, Errno(ENODEV)
}

// LinkByName resolves a network link by name from the Linux link table.
func LinkByName(name string) (*Link, error) {
	if name == "" || len(name) >= IFNAMSIZ {
		return nil, Errno(EINVAL)
	}
	links, err := Links()
	return findLinkByName(name, links, err)
}

// LinkByIndex resolves a network link by index from the Linux link table.
func LinkByIndex(index int) (*Link, error) {
	if index <= 0 {
		return nil, Errno(EINVAL)
	}
	links, err := Links()
	return findLinkByIndex(index, links, err)
}

// IfNameToIndex resolves a network link name to its kernel index.
func IfNameToIndex(name string) (uint32, error) {
	link, err := LinkByName(name)
	if err != nil {
		return 0, err
	}
	return uint32(link.Index), nil
}

func openRouteNetlink() (uintptr, uint32, error) {
	fd, errno := Socket(AF_NETLINK, SOCK_RAW|SOCK_CLOEXEC, NETLINK_ROUTE)
	return openRouteNetlinkFD(fd, errno)
}

func openRouteNetlinkFD(fd uintptr, socketErr uintptr) (uintptr, uint32, error) {
	errno := socketErr
	if errno != 0 {
		return 0, 0, Errno(errno)
	}

	sa := sockaddrNetlink{Family: AF_NETLINK}
	if errno = Bind(fd, unsafe.Pointer(&sa), uintptr(unsafe.Sizeof(sa))); errno != 0 {
		Close(fd)
		return 0, 0, Errno(errno)
	}

	var local sockaddrNetlink
	localLen := uint32(unsafe.Sizeof(local))
	if errno = Getsockname(fd, unsafe.Pointer(&local), unsafe.Pointer(&localLen)); errno != 0 {
		Close(fd)
		return 0, 0, Errno(errno)
	}
	return fd, local.Pid, nil
}

func sendNetlinkGetLinkDump(fd uintptr, seq uint32) error {
	const (
		nlmsgHdrLen   = uint32(unsafe.Sizeof(nlMsghdr{}))
		rtGenmsgSize  = uint32(unsafe.Sizeof(rtGenmsg{}))
		requestLength = nlmsgHdrLen + rtGenmsgSize
	)

	var req [requestLength]byte
	hdr := (*nlMsghdr)(unsafe.Pointer(&req[0]))
	hdr.Len = requestLength
	hdr.Type = RTM_GETLINK
	hdr.Flags = NLM_F_REQUEST | NLM_F_DUMP
	hdr.Seq = seq
	(*rtGenmsg)(unsafe.Pointer(&req[nlmsgHdrLen])).Family = AF_UNSPEC

	sa := sockaddrNetlink{Family: AF_NETLINK}
	if _, errno := Sendto(fd, req[:], 0, unsafe.Pointer(&sa), uintptr(unsafe.Sizeof(sa))); errno != 0 {
		return Errno(errno)
	}
	return nil
}

func recvLinks(fd uintptr, pid, seq uint32) ([]Link, error) {
	buf := make([]byte, 64*1024)
	links := make([]Link, 0, 16)
	for {
		n, errno := Recvfrom(fd, buf, 0, nil, nil)
		if errno != 0 {
			return nil, Errno(errno)
		}
		if int(n) < int(unsafe.Sizeof(nlMsghdr{})) {
			return nil, Errno(EINVAL)
		}

		msgs := buf[:n]
		for len(msgs) >= int(unsafe.Sizeof(nlMsghdr{})) {
			hdr := (*nlMsghdr)(unsafe.Pointer(&msgs[0]))
			msgLen := int(hdr.Len)
			alignedLen := nlmsgAlignOf(msgLen)
			if msgLen < int(unsafe.Sizeof(nlMsghdr{})) || alignedLen > len(msgs) {
				return nil, Errno(EINVAL)
			}
			if hdr.Seq != seq || hdr.Pid != pid {
				return nil, Errno(EINVAL)
			}

			payload := msgs[unsafe.Sizeof(nlMsghdr{}):msgLen]
			switch hdr.Type {
			case NLMSG_DONE:
				return links, nil
			case NLMSG_ERROR:
				if len(payload) < 4 {
					return nil, Errno(EINVAL)
				}
				errno := *(*int32)(unsafe.Pointer(&payload[0]))
				if errno >= 0 {
					return nil, Errno(EINVAL)
				}
				return nil, Errno(-errno)
			case RTM_NEWLINK:
				link, ok, err := parseLink(payload)
				if err != nil {
					return nil, err
				}
				if ok {
					links = append(links, link)
				}
			}

			msgs = msgs[alignedLen:]
		}
	}
}

func parseLink(payload []byte) (Link, bool, error) {
	if len(payload) < int(unsafe.Sizeof(ifInfomsg{})) {
		return Link{}, false, nil
	}
	ifim := (*ifInfomsg)(unsafe.Pointer(&payload[0]))
	if ifim.Index <= 0 {
		return Link{}, false, nil
	}
	link := Link{Index: int(ifim.Index), Flags: ifim.Flags}
	attrs := payload[unsafe.Sizeof(ifInfomsg{}):]
	for len(attrs) >= int(unsafe.Sizeof(rtAttr{})) {
		attr := (*rtAttr)(unsafe.Pointer(&attrs[0]))
		attrLen := int(attr.Len)
		alignedLen := rtaAlignOf(attrLen)
		if attrLen < int(unsafe.Sizeof(rtAttr{})) || alignedLen > len(attrs) {
			return Link{}, false, Errno(EINVAL)
		}
		value := attrs[unsafe.Sizeof(rtAttr{}):attrLen]
		switch attr.Type {
		case IFLA_IFNAME:
			if len(value) != 0 && value[len(value)-1] == 0 {
				value = value[:len(value)-1]
			}
			link.Name = string(value)
		case IFLA_MTU:
			if len(value) >= 4 {
				link.MTU = int(*(*uint32)(unsafe.Pointer(&value[0])))
			}
		case IFLA_ADDRESS:
			if len(value) != 0 {
				link.HardwareAddr = append([]byte(nil), value...)
			}
		}
		attrs = attrs[alignedLen:]
	}
	if link.Name == "" {
		return Link{}, false, nil
	}
	return link, true, nil
}

func nlmsgAlignOf(length int) int {
	return (length + nlmsgAlignTo - 1) & ^(nlmsgAlignTo - 1)
}

func rtaAlignOf(length int) int {
	return (length + rtaAlignTo - 1) & ^(rtaAlignTo - 1)
}

// Bind binds a socket to an address.
func Bind(fd uintptr, addr unsafe.Pointer, addrlen uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_BIND, fd, uintptr(noescape(addr)), addrlen, 0)
	return
}

// Listen marks a socket as listening.
func Listen(fd, backlog uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_LISTEN, fd, backlog, 0, 0)
	return
}

// Accept accepts a connection on a socket.
func Accept(fd uintptr, addr unsafe.Pointer, addrlen unsafe.Pointer) (nfd uintptr, errno uintptr) {
	return Syscall4(SYS_ACCEPT, fd, uintptr(noescape(addr)), uintptr(noescape(addrlen)), 0)
}

// Accept4 accepts a connection on a socket with flags.
func Accept4(fd uintptr, addr unsafe.Pointer, addrlen unsafe.Pointer, flags uintptr) (nfd uintptr, errno uintptr) {
	return Syscall4(SYS_ACCEPT4, fd, uintptr(noescape(addr)), uintptr(noescape(addrlen)), flags)
}

// Connect connects a socket to an address.
func Connect(fd uintptr, addr unsafe.Pointer, addrlen uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_CONNECT, fd, uintptr(noescape(addr)), addrlen, 0)
	return
}

// Shutdown shuts down part of a full-duplex connection.
func Shutdown(fd, how uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_SHUTDOWN, fd, how, 0, 0)
	return
}

// Setsockopt sets a socket option.
func Setsockopt(fd, level, optname uintptr, optval unsafe.Pointer, optlen uintptr) (errno uintptr) {
	_, errno = Syscall6(SYS_SETSOCKOPT, fd, level, optname, uintptr(noescape(optval)), optlen, 0)
	return
}

// Getsockopt gets a socket option.
func Getsockopt(fd, level, optname uintptr, optval unsafe.Pointer, optlen unsafe.Pointer) (errno uintptr) {
	_, errno = Syscall6(SYS_GETSOCKOPT, fd, level, optname, uintptr(noescape(optval)), uintptr(noescape(optlen)), 0)
	return
}

// Socketpair creates a pair of connected sockets.
func Socketpair(domain, typ, protocol uintptr, fds *[2]int32) (errno uintptr) {
	_, errno = Syscall4(SYS_SOCKETPAIR, domain, typ, protocol, uintptr(noescape(unsafe.Pointer(fds))))
	return
}

// Getsockname gets the local address of a socket.
func Getsockname(fd uintptr, addr unsafe.Pointer, addrlen unsafe.Pointer) (errno uintptr) {
	_, errno = Syscall4(SYS_GETSOCKNAME, fd, uintptr(noescape(addr)), uintptr(noescape(addrlen)), 0)
	return
}

// Getpeername gets the remote address of a socket.
func Getpeername(fd uintptr, addr unsafe.Pointer, addrlen unsafe.Pointer) (errno uintptr) {
	_, errno = Syscall4(SYS_GETPEERNAME, fd, uintptr(noescape(addr)), uintptr(noescape(addrlen)), 0)
	return
}

// Sendto sends a message on a socket.
func Sendto(fd uintptr, buf []byte, flags uintptr, addr unsafe.Pointer, addrlen uintptr) (n uintptr, errno uintptr) {
	var p uintptr
	if len(buf) > 0 {
		p = uintptr(noescape(unsafe.Pointer(&buf[0])))
	}
	return Syscall6(SYS_SENDTO, fd, p, uintptr(len(buf)), flags, uintptr(noescape(addr)), addrlen)
}

// Recvfrom receives a message from a socket.
func Recvfrom(fd uintptr, buf []byte, flags uintptr, addr unsafe.Pointer, addrlen unsafe.Pointer) (n uintptr, errno uintptr) {
	var p uintptr
	if len(buf) > 0 {
		p = uintptr(noescape(unsafe.Pointer(&buf[0])))
	}
	return Syscall6(SYS_RECVFROM, fd, p, uintptr(len(buf)), flags, uintptr(noescape(addr)), uintptr(noescape(addrlen)))
}

// Sendmsg sends a message on a socket using a msghdr.
func Sendmsg(fd uintptr, msg unsafe.Pointer, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_SENDMSG, fd, uintptr(noescape(msg)), flags, 0)
}

// Recvmsg receives a message from a socket using a msghdr.
func Recvmsg(fd uintptr, msg unsafe.Pointer, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_RECVMSG, fd, uintptr(noescape(msg)), flags, 0)
}

// Sendmmsg sends multiple messages on a socket.
func Sendmmsg(fd uintptr, msgvec unsafe.Pointer, vlen, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_SENDMMSG, fd, uintptr(noescape(msgvec)), vlen, flags)
}

// Recvmmsg receives multiple messages from a socket.
func Recvmmsg(fd uintptr, msgvec unsafe.Pointer, vlen, flags uintptr, timeout unsafe.Pointer) (n uintptr, errno uintptr) {
	return Syscall6(SYS_RECVMMSG, fd, uintptr(noescape(msgvec)), vlen, flags, uintptr(noescape(timeout)), 0)
}

// Readv reads into multiple buffers.
func Readv(fd uintptr, iov unsafe.Pointer, iovcnt uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_READV, fd, uintptr(noescape(iov)), iovcnt, 0)
}

// Writev writes from multiple buffers.
func Writev(fd uintptr, iov unsafe.Pointer, iovcnt uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_WRITEV, fd, uintptr(noescape(iov)), iovcnt, 0)
}

// Preadv reads into multiple buffers at a given offset.
func Preadv(fd uintptr, iov unsafe.Pointer, iovcnt uintptr, offset int64) (n uintptr, errno uintptr) {
	return Syscall4(SYS_PREADV, fd, uintptr(noescape(iov)), iovcnt, uintptr(offset))
}

// Pwritev writes from multiple buffers at a given offset.
func Pwritev(fd uintptr, iov unsafe.Pointer, iovcnt uintptr, offset int64) (n uintptr, errno uintptr) {
	return Syscall4(SYS_PWRITEV, fd, uintptr(noescape(iov)), iovcnt, uintptr(offset))
}

// Preadv2 reads into multiple buffers at a given offset with flags.
func Preadv2(fd uintptr, iov unsafe.Pointer, iovcnt uintptr, offset int64, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall6(SYS_PREADV2, fd, uintptr(noescape(iov)), iovcnt, uintptr(offset), 0, flags)
}

// Pwritev2 writes from multiple buffers at a given offset with flags.
func Pwritev2(fd uintptr, iov unsafe.Pointer, iovcnt uintptr, offset int64, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall6(SYS_PWRITEV2, fd, uintptr(noescape(iov)), iovcnt, uintptr(offset), 0, flags)
}

// Pipe2 creates a pipe with flags.
func Pipe2(fds *[2]int32, flags uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_PIPE2, uintptr(noescape(unsafe.Pointer(fds))), flags, 0, 0)
	return
}

// Splice moves data between two file descriptors.
func Splice(fdIn uintptr, offIn *int64, fdOut uintptr, offOut *int64, length, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall6(SYS_SPLICE, fdIn, uintptr(noescape(unsafe.Pointer(offIn))), fdOut, uintptr(noescape(unsafe.Pointer(offOut))), length, flags)
}

// Tee duplicates data between two pipe file descriptors.
func Tee(fdIn, fdOut, length, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_TEE, fdIn, fdOut, length, flags)
}

// Vmsplice maps user pages into a pipe.
func Vmsplice(fd uintptr, iov unsafe.Pointer, nrSegs, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_VMSPLICE, fd, uintptr(noescape(iov)), nrSegs, flags)
}

// Eventfd2 creates an eventfd.
func Eventfd2(initval, flags uintptr) (fd uintptr, errno uintptr) {
	return Syscall4(SYS_EVENTFD2, initval, flags, 0, 0)
}

// Signalfd4 creates or modifies a file descriptor for signal handling.
// If fd is -1 (^uintptr(0)), a new signalfd is created; otherwise, the existing
// fd is modified. The mask points to a sigset_t specifying which signals to accept.
// Flags may include SFD_NONBLOCK and SFD_CLOEXEC.
// Returns the file descriptor and errno.
func Signalfd4(fd uintptr, mask unsafe.Pointer, maskSize, flags uintptr) (newfd uintptr, errno uintptr) {
	return Syscall4(SYS_SIGNALFD4, fd, uintptr(noescape(mask)), maskSize, flags)
}

// PidfdOpen refers to a process.
func PidfdOpen(pid, flags uintptr) (fd uintptr, errno uintptr) {
	return Syscall4(SYS_PIDFD_OPEN, pid, flags, 0, 0)
}

// PidfdGetfd duplicates a file descriptor from another process.
func PidfdGetfd(pidfd, targetfd, flags uintptr) (fd uintptr, errno uintptr) {
	return Syscall4(SYS_PIDFD_GETFD, pidfd, targetfd, flags, 0)
}

// PidfdSendSignal sends a signal to a process.
func PidfdSendSignal(pidfd, sig uintptr, info unsafe.Pointer, flags uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_PIDFD_SEND_SIGNAL, pidfd, sig, uintptr(noescape(info)), flags)
	return
}

// MemfdCreate creates an anonymous, memory-backed file descriptor.
// The name is a null-terminated string used for debugging (visible in /proc/self/fd/).
// Flags may include MFD_CLOEXEC, MFD_ALLOW_SEALING, and MFD_HUGETLB.
// Returns the file descriptor and errno.
func MemfdCreate(name unsafe.Pointer, flags uintptr) (fd uintptr, errno uintptr) {
	return Syscall4(SYS_MEMFD_CREATE, uintptr(noescape(name)), flags, 0, 0)
}

// TimerfdCreate creates a timerfd.
func TimerfdCreate(clockid, flags uintptr) (fd uintptr, errno uintptr) {
	return Syscall4(SYS_TIMERFD_CREATE, clockid, flags, 0, 0)
}

// TimerfdSettime arms or disarms a timerfd.
func TimerfdSettime(fd, flags uintptr, newValue, oldValue unsafe.Pointer) (errno uintptr) {
	_, errno = Syscall4(SYS_TIMERFD_SETTIME, fd, flags, uintptr(noescape(newValue)), uintptr(noescape(oldValue)))
	return
}

// TimerfdGettime gets the current setting of a timerfd.
func TimerfdGettime(fd uintptr, currValue unsafe.Pointer) (errno uintptr) {
	_, errno = Syscall4(SYS_TIMERFD_GETTIME, fd, uintptr(noescape(currValue)), 0, 0)
	return
}

// IoUringSetup sets up an io_uring instance.
func IoUringSetup(entries uintptr, params unsafe.Pointer) (fd uintptr, errno uintptr) {
	return Syscall4(SYS_IO_URING_SETUP, entries, uintptr(noescape(params)), 0, 0)
}

// IoUringEnter submits I/O requests and/or waits for completions.
func IoUringEnter(fd, toSubmit, minComplete, flags uintptr, sig unsafe.Pointer, sigsetSize uintptr) (r1 uintptr, errno uintptr) {
	return Syscall6(SYS_IO_URING_ENTER, fd, toSubmit, minComplete, flags, uintptr(noescape(sig)), sigsetSize)
}

// IoUringRegister registers resources with an io_uring instance.
func IoUringRegister(fd, opcode uintptr, arg unsafe.Pointer, nrArgs uintptr) (r1 uintptr, errno uintptr) {
	return Syscall4(SYS_IO_URING_REGISTER, fd, opcode, uintptr(noescape(arg)), nrArgs)
}

// Mmap maps files or devices into memory.
// Returns unsafe.Pointer to enable vet-clean pointer arithmetic with unsafe.Add.
func Mmap(addr unsafe.Pointer, length, prot, flags, fd, offset uintptr) (ptr unsafe.Pointer, errno uintptr) {
	r1, errno := Syscall6(SYS_MMAP, uintptr(noescape(addr)), length, prot, flags, fd, offset)
	return asPointer(r1), errno
}

// Munmap unmaps files or devices from memory.
func Munmap(addr unsafe.Pointer, length uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_MUNMAP, uintptr(addr), length, 0, 0)
	return
}
