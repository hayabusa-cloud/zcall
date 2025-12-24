// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build darwin && arm64

package zcall

// Syscall numbers for Darwin on arm64.
// Reference: /usr/include/sys/syscall.h (XNU kernel)
// Darwin syscall numbers are BSD-style without the 0x2000000 class prefix.
// The class prefix is added by the caller when needed.
const (
	// Basic I/O
	SYS_READ   = 3
	SYS_WRITE  = 4
	SYS_CLOSE  = 6
	SYS_MMAP   = 197
	SYS_MUNMAP = 73

	// Vectored I/O
	SYS_READV   = 120
	SYS_WRITEV  = 121
	SYS_PREADV  = 267
	SYS_PWRITEV = 268

	// Pipe
	SYS_PIPE = 42

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

	// kqueue (Darwin's event notification)
	SYS_KQUEUE   = 362
	SYS_KEVENT   = 363
	SYS_KEVENT64 = 369
)
