// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build freebsd && amd64

package zcall

// Syscall numbers for FreeBSD on amd64.
// Reference: /usr/include/sys/syscall.h (FreeBSD 14.x)
const (
	// Basic I/O
	SYS_READ      = 3
	SYS_WRITE     = 4
	SYS_CLOSE     = 6
	SYS_FSTAT     = 551 // freebsd12_fstat
	SYS_MMAP      = 477 // freebsd6_mmap
	SYS_MUNMAP    = 73
	SYS_FTRUNCATE = 480 // freebsd6_ftruncate

	// Vectored I/O
	SYS_READV   = 120
	SYS_WRITEV  = 121
	SYS_PREADV  = 289
	SYS_PWRITEV = 290

	// Pipe
	SYS_PIPE2 = 542

	// Networking - basic
	SYS_SOCKET      = 97
	SYS_CONNECT     = 98
	SYS_ACCEPT      = 30
	SYS_SENDTO      = 133
	SYS_RECVFROM    = 29
	SYS_SENDMSG     = 28
	SYS_RECVMSG     = 27
	SYS_SHUTDOWN    = 134
	SYS_BIND        = 104
	SYS_LISTEN      = 106
	SYS_GETSOCKNAME = 32
	SYS_GETPEERNAME = 31
	SYS_SOCKETPAIR  = 135
	SYS_SETSOCKOPT  = 105
	SYS_GETSOCKOPT  = 118
	SYS_ACCEPT4     = 541

	// Multi-message (FreeBSD 11+)
	SYS_SENDMMSG = 475
	SYS_RECVMMSG = 474
)

// Socket constants for FreeBSD.
const (
	AF_UNIX  = 1
	AF_INET  = 2
	AF_INET6 = 28

	SOCK_STREAM    = 1
	SOCK_DGRAM     = 2
	SOCK_RAW       = 3
	SOCK_SEQPACKET = 5
	SOCK_NONBLOCK  = 0x20000000
	SOCK_CLOEXEC   = 0x10000000

	IPPROTO_TCP  = 6
	IPPROTO_UDP  = 17
	IPPROTO_SCTP = 132
)

// File descriptor flags for FreeBSD.
const (
	O_NONBLOCK = 0x4
	O_CLOEXEC  = 0x100000
)
