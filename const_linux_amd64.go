// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux && amd64

package zcall

// Syscall numbers for Linux on amd64.
// Reference: arch/x86/entry/syscalls/syscall_64.tbl
const (
	// Basic I/O
	SYS_READ      = 0
	SYS_WRITE     = 1
	SYS_CLOSE     = 3
	SYS_FSTAT     = 5
	SYS_MMAP      = 9
	SYS_MUNMAP    = 11
	SYS_FTRUNCATE = 77

	// Vectored I/O
	SYS_READV    = 19
	SYS_WRITEV   = 20
	SYS_PREADV   = 295
	SYS_PWRITEV  = 296
	SYS_PREADV2  = 327
	SYS_PWRITEV2 = 328

	// Networking - basic
	SYS_SOCKET      = 41
	SYS_CONNECT     = 42
	SYS_ACCEPT      = 43
	SYS_SENDTO      = 44
	SYS_RECVFROM    = 45
	SYS_SENDMSG     = 46
	SYS_RECVMSG     = 47
	SYS_SHUTDOWN    = 48
	SYS_BIND        = 49
	SYS_LISTEN      = 50
	SYS_GETSOCKNAME = 51
	SYS_GETPEERNAME = 52
	SYS_SOCKETPAIR  = 53
	SYS_SETSOCKOPT  = 54
	SYS_GETSOCKOPT  = 55

	// Zero-copy and pipe
	SYS_SPLICE   = 275
	SYS_TEE      = 276
	SYS_VMSPLICE = 278
	SYS_PIPE2    = 293

	// Timers and events
	SYS_TIMERFD_CREATE  = 283
	SYS_TIMERFD_SETTIME = 286
	SYS_TIMERFD_GETTIME = 287
	SYS_ACCEPT4         = 288
	SYS_EVENTFD2        = 290

	// Multi-message
	SYS_RECVMMSG = 299
	SYS_SENDMMSG = 307

	// io_uring
	SYS_IO_URING_SETUP    = 425
	SYS_IO_URING_ENTER    = 426
	SYS_IO_URING_REGISTER = 427

	// signalfd, pidfd, memfd
	SYS_SIGNALFD4          = 289
	SYS_MEMFD_CREATE       = 319
	SYS_PIDFD_SEND_SIGNAL  = 424
	SYS_PIDFD_OPEN         = 434
	SYS_PIDFD_GETFD        = 438
)
