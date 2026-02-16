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

// Socket address families.
const (
	AF_UNIX  = 1
	AF_LOCAL = AF_UNIX
	AF_INET  = 2
	AF_INET6 = 28
)

// Socket types.
const (
	SOCK_STREAM    = 1
	SOCK_DGRAM     = 2
	SOCK_RAW       = 3
	SOCK_SEQPACKET = 5
	SOCK_NONBLOCK  = 0x20000000
	SOCK_CLOEXEC   = 0x10000000
)

// IP protocols.
const (
	IPPROTO_IP   = 0
	IPPROTO_ICMP = 1
	IPPROTO_TCP  = 6
	IPPROTO_UDP  = 17
	IPPROTO_IPV6 = 41
	IPPROTO_SCTP = 132
	IPPROTO_RAW  = 255
)

// Socket option levels.
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
	SO_ACCEPTCONN = 0x0002
	SO_REUSEADDR  = 0x0004
	SO_KEEPALIVE  = 0x0008
	SO_DONTROUTE  = 0x0010
	SO_BROADCAST  = 0x0020
	SO_LINGER     = 0x0080
	SO_OOBINLINE  = 0x0100
	SO_REUSEPORT  = 0x0200
	SO_SNDBUF     = 0x1001
	SO_RCVBUF     = 0x1002
	SO_SNDLOWAT   = 0x1003
	SO_RCVLOWAT   = 0x1004
	SO_SNDTIMEO   = 0x1005
	SO_RCVTIMEO   = 0x1006
	SO_ERROR      = 0x1007
	SO_TYPE       = 0x1008
	SO_ZEROCOPY   = 0x0 // Not supported on FreeBSD
)

// TCP options.
const (
	TCP_NODELAY   = 1
	TCP_MAXSEG    = 2
	TCP_NOPUSH    = 4
	TCP_KEEPINIT  = 128
	TCP_KEEPIDLE  = 256
	TCP_KEEPINTVL = 512
	TCP_KEEPCNT   = 1024
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
	MSG_EOR       = 0x8
	MSG_TRUNC     = 0x10
	MSG_CTRUNC    = 0x20
	MSG_WAITALL   = 0x40
	MSG_DONTWAIT  = 0x80
	MSG_EOF       = 0x100
	MSG_NOSIGNAL  = 0x20000
	MSG_ZEROCOPY  = 0x0 // Not supported on FreeBSD
)

// File descriptor flags.
const (
	O_RDONLY   = 0x0
	O_WRONLY   = 0x1
	O_RDWR     = 0x2
	O_NONBLOCK = 0x4
	O_APPEND   = 0x8
	O_CREAT    = 0x200
	O_TRUNC    = 0x400
	O_EXCL     = 0x800
	O_CLOEXEC  = 0x100000
)

// Iovec is the I/O vector for vectored I/O operations.
type Iovec struct {
	Base *byte
	Len  uint64
}

// Timespec represents a time value with nanosecond precision.
type Timespec struct {
	Sec  int64
	Nsec int64
}
