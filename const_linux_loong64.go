// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux && loong64

package zcall

// Syscall numbers for Linux on loong64.
// Reference: include/uapi/asm-generic/unistd.h (loong64 uses the generic table)
const (
	// Basic I/O
	SYS_FTRUNCATE = 46
	SYS_CLOSE     = 57
	SYS_READ      = 63
	SYS_WRITE     = 64
	SYS_FSTAT     = 80
	SYS_MUNMAP    = 215
	SYS_MMAP      = 222

	// Vectored I/O
	SYS_READV    = 65
	SYS_WRITEV   = 66
	SYS_PREADV   = 69
	SYS_PWRITEV  = 70
	SYS_PREADV2  = 286
	SYS_PWRITEV2 = 287

	// Zero-copy and pipe
	SYS_PIPE2    = 59
	SYS_VMSPLICE = 75
	SYS_SPLICE   = 76
	SYS_TEE      = 77

	// Timers and events
	SYS_EVENTFD2        = 19
	SYS_TIMERFD_CREATE  = 85
	SYS_TIMERFD_SETTIME = 86
	SYS_TIMERFD_GETTIME = 87

	// Networking - basic
	SYS_SOCKET      = 198
	SYS_SOCKETPAIR  = 199
	SYS_BIND        = 200
	SYS_LISTEN      = 201
	SYS_ACCEPT      = 202
	SYS_CONNECT     = 203
	SYS_GETSOCKNAME = 204
	SYS_GETPEERNAME = 205
	SYS_SENDTO      = 206
	SYS_RECVFROM    = 207
	SYS_SETSOCKOPT  = 208
	SYS_GETSOCKOPT  = 209
	SYS_SHUTDOWN    = 210
	SYS_SENDMSG     = 211
	SYS_RECVMSG     = 212
	SYS_ACCEPT4     = 242

	// Multi-message
	SYS_RECVMMSG = 243
	SYS_SENDMMSG = 269

	// io_uring
	SYS_IO_URING_SETUP    = 425
	SYS_IO_URING_ENTER    = 426
	SYS_IO_URING_REGISTER = 427

	// signalfd, pidfd, memfd
	SYS_SIGNALFD4         = 74
	SYS_MEMFD_CREATE      = 279
	SYS_PIDFD_SEND_SIGNAL = 424
	SYS_PIDFD_OPEN        = 434
	SYS_PIDFD_GETFD       = 438
)
