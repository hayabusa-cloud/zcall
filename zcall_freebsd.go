// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build freebsd

package zcall

import (
	"unsafe"

	"code.hybscloud.com/zcall/internal"
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
	if len(buf) == 0 {
		return Syscall6(SYS_SENDTO, fd, 0, 0, flags, uintptr(noescape(addr)), addrlen)
	}
	return Syscall6(SYS_SENDTO, fd, uintptr(noescape(unsafe.Pointer(&buf[0]))), uintptr(len(buf)), flags, uintptr(noescape(addr)), addrlen)
}

// Recvfrom receives a message from a socket.
func Recvfrom(fd uintptr, buf []byte, flags uintptr, addr unsafe.Pointer, addrlen unsafe.Pointer) (n uintptr, errno uintptr) {
	if len(buf) == 0 {
		return Syscall6(SYS_RECVFROM, fd, 0, 0, flags, uintptr(noescape(addr)), uintptr(noescape(addrlen)))
	}
	return Syscall6(SYS_RECVFROM, fd, uintptr(noescape(unsafe.Pointer(&buf[0]))), uintptr(len(buf)), flags, uintptr(noescape(addr)), uintptr(noescape(addrlen)))
}

// Sendmsg sends a message on a socket using a msghdr structure.
func Sendmsg(fd uintptr, msg unsafe.Pointer, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_SENDMSG, fd, uintptr(noescape(msg)), flags, 0)
}

// Recvmsg receives a message from a socket using a msghdr structure.
func Recvmsg(fd uintptr, msg unsafe.Pointer, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_RECVMSG, fd, uintptr(noescape(msg)), flags, 0)
}

// Sendmmsg sends multiple messages on a socket.
func Sendmmsg(fd uintptr, msgvec unsafe.Pointer, vlen uintptr, flags uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_SENDMMSG, fd, uintptr(noescape(msgvec)), vlen, flags)
}

// Recvmmsg receives multiple messages from a socket.
func Recvmmsg(fd uintptr, msgvec unsafe.Pointer, vlen uintptr, flags uintptr, timeout unsafe.Pointer) (n uintptr, errno uintptr) {
	return Syscall6(SYS_RECVMMSG, fd, uintptr(noescape(msgvec)), vlen, flags, uintptr(noescape(timeout)), 0)
}

// Readv reads from a file descriptor into multiple buffers.
func Readv(fd uintptr, iov unsafe.Pointer, iovcnt uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_READV, fd, uintptr(noescape(iov)), iovcnt, 0)
}

// Writev writes to a file descriptor from multiple buffers.
func Writev(fd uintptr, iov unsafe.Pointer, iovcnt uintptr) (n uintptr, errno uintptr) {
	return Syscall4(SYS_WRITEV, fd, uintptr(noescape(iov)), iovcnt, 0)
}

// Preadv reads from a file descriptor at an offset into multiple buffers.
func Preadv(fd uintptr, iov unsafe.Pointer, iovcnt uintptr, offset int64) (n uintptr, errno uintptr) {
	return Syscall6(SYS_PREADV, fd, uintptr(noescape(iov)), iovcnt, uintptr(offset), 0, 0)
}

// Pwritev writes to a file descriptor at an offset from multiple buffers.
func Pwritev(fd uintptr, iov unsafe.Pointer, iovcnt uintptr, offset int64) (n uintptr, errno uintptr) {
	return Syscall6(SYS_PWRITEV, fd, uintptr(noescape(iov)), iovcnt, uintptr(offset), 0, 0)
}

// Pipe2 creates a pipe with flags.
func Pipe2(fds *[2]int32, flags uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_PIPE2, uintptr(noescape(unsafe.Pointer(fds))), flags, 0, 0)
	return
}

// Mmap maps files or devices into memory.
func Mmap(addr unsafe.Pointer, length uintptr, prot uintptr, flags uintptr, fd uintptr, offset uintptr) (r1 uintptr, errno uintptr) {
	return Syscall6(SYS_MMAP, uintptr(addr), length, prot, flags, fd, offset)
}

// Munmap unmaps files or devices from memory.
func Munmap(addr uintptr, length uintptr) (errno uintptr) {
	_, errno = Syscall4(SYS_MUNMAP, addr, length, 0, 0)
	return
}
