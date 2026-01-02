// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build darwin

package zcall

import (
	"unsafe"

	"code.hybscloud.com/zcall/internal"
)

const bsdClass = 0x2000000

// Syscall4 executes a syscall with up to 4 arguments.
func Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr) {
	return internal.RawSyscall4(num|bsdClass, a1, a2, a3, a4)
}

// Syscall6 executes a syscall with up to 6 arguments.
func Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr) {
	return internal.RawSyscall6(num|bsdClass, a1, a2, a3, a4, a5, a6)
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

// Mmap maps files or devices into memory.
// Returns unsafe.Pointer to enable vet-clean pointer arithmetic with unsafe.Add.
func Mmap(addr unsafe.Pointer, length, prot, flags, fd, offset uintptr) (ptr unsafe.Pointer, errno uintptr) {
	r1, errno := Syscall6(SYS_MMAP, uintptr(noescape(addr)), length, prot, flags, fd, offset)
	return unsafe.Pointer(r1), errno
}

// Munmap unmaps files or devices from memory.
func Munmap(addr unsafe.Pointer, length uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_MUNMAP, uintptr(addr), length, 0, 0)
	return
}

// Socket address families.
const (
	AF_UNIX  = 1
	AF_LOCAL = AF_UNIX
	AF_INET  = 2
	AF_INET6 = 30
)

// Socket types.
const (
	SOCK_STREAM    = 1
	SOCK_DGRAM     = 2
	SOCK_RAW       = 3
	SOCK_SEQPACKET = 5
	SOCK_NONBLOCK  = 0x20000000 // Darwin-specific
	SOCK_CLOEXEC   = 0x10000000 // Darwin-specific
)

// IP protocols.
const (
	IPPROTO_IP   = 0
	IPPROTO_ICMP = 1
	IPPROTO_TCP  = 6
	IPPROTO_UDP  = 17
	IPPROTO_IPV6 = 41
	IPPROTO_RAW  = 255
)

// Socket options levels.
const (
	SOL_SOCKET = 0xffff
	SOL_IP     = 0
	SOL_TCP    = 6
	SOL_UDP    = 17
	SOL_IPV6   = 41
)

// Socket options (SOL_SOCKET level).
const (
	SO_DEBUG      = 0x0001
	SO_REUSEADDR  = 0x0004
	SO_TYPE       = 0x1008
	SO_ERROR      = 0x1007
	SO_DONTROUTE  = 0x0010
	SO_BROADCAST  = 0x0020
	SO_SNDBUF     = 0x1001
	SO_RCVBUF     = 0x1002
	SO_KEEPALIVE  = 0x0008
	SO_OOBINLINE  = 0x0100
	SO_LINGER     = 0x0080
	SO_REUSEPORT  = 0x0200
	SO_RCVLOWAT   = 0x1004
	SO_SNDLOWAT   = 0x1003
	SO_RCVTIMEO   = 0x1006
	SO_SNDTIMEO   = 0x1005
	SO_ACCEPTCONN = 0x0002
	SO_NOSIGPIPE  = 0x1022
)

// TCP options.
const (
	TCP_NODELAY   = 0x01
	TCP_MAXSEG    = 0x02
	TCP_KEEPALIVE = 0x10
	TCP_KEEPINTVL = 0x101
	TCP_KEEPCNT   = 0x102
)

// File descriptor flags.
const (
	O_RDONLY   = 0x0
	O_WRONLY   = 0x1
	O_RDWR     = 0x2
	O_CREAT    = 0x200
	O_EXCL     = 0x800
	O_NOCTTY   = 0x20000
	O_TRUNC    = 0x400
	O_APPEND   = 0x8
	O_NONBLOCK = 0x4
	O_SYNC     = 0x80
	O_CLOEXEC  = 0x1000000
)

// mmap protection flags.
const (
	PROT_NONE  = 0x0
	PROT_READ  = 0x1
	PROT_WRITE = 0x2
	PROT_EXEC  = 0x4
)

// mmap flags.
const (
	MAP_SHARED    = 0x1
	MAP_PRIVATE   = 0x2
	MAP_FIXED     = 0x10
	MAP_ANONYMOUS = 0x1000
)

// Shutdown how.
const (
	SHUT_RD   = 0
	SHUT_WR   = 1
	SHUT_RDWR = 2
)

// MSG flags for send/recv.
const (
	MSG_OOB       = 0x1
	MSG_PEEK      = 0x2
	MSG_DONTROUTE = 0x4
	MSG_TRUNC     = 0x10
	MSG_CTRUNC    = 0x20
	MSG_WAITALL   = 0x40
	MSG_DONTWAIT  = 0x80
	MSG_EOF       = 0x100
	MSG_NOSIGNAL  = 0x80000
)

// Iovec represents a scatter/gather I/O vector.
type Iovec struct {
	Base *byte
	Len  uint64
}

// Timespec represents a time value with nanosecond precision.
type Timespec struct {
	Sec  int64
	Nsec int64
}
