// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Package zcall provides zero-overhead syscall primitives for Linux and Darwin.
//
// # Overview
//
// zcall bypasses Go's runtime syscall machinery (entersyscall/exitsyscall)
// by invoking the kernel directly via raw assembly. This eliminates the
// latency tax imposed by Go's scheduler hooks, making it suitable for
// low-latency I/O paths such as io_uring submission and completion.
//
// # Design Principles
//
// Zero Overhead: All syscalls are implemented in pure assembly without
// any Go runtime interaction. The caller is responsible for cooperative
// scheduling.
//
// Zero Dependencies: This package does not import "syscall" or
// "golang.org/x/sys/unix". All syscall numbers and constants are
// defined internally.
//
// Raw Semantics: Functions return (result, errno) directly from the kernel.
// errno == 0 indicates success; errno != 0 is the raw kernel error number.
// The caller must handle error translation using zcall.Errno(errno).
//
// # Supported Architectures
//
//   - linux/amd64: Uses SYSCALL instruction
//   - linux/arm64: Uses SVC #0 instruction
//   - linux/riscv64: Uses ECALL instruction
//   - linux/loong64: Uses SYSCALL instruction
//   - darwin/arm64: Uses SVC #0x80 instruction
//   - freebsd/amd64: Uses SYSCALL instruction
//
// # Usage
//
// zcall is designed as a building block for high-performance I/O libraries.
// Direct usage requires understanding of Linux syscall semantics and
// careful attention to memory safety.
//
//	r1, errno := zcall.Syscall6(zcall.SYS_IO_URING_ENTER, fd, toSubmit, minComplete, flags, sigset, sigsetSize)
//	if errno != 0 {
//	    // handle error
//	}
//
// # Safety
//
// This package uses unsafe operations for pointer-to-uintptr conversion.
// Callers must ensure that:
//   - Pointers passed to syscalls remain valid for the duration of the call
//   - Memory referenced by pointers is not garbage collected during the call
//   - Proper synchronization is used for concurrent access
//
// # Manual Cooperation
//
// Since zcall bypasses the Go scheduler's syscall hooks, long-running
// syscalls may starve other goroutines. Callers should:
//   - Use non-blocking syscalls where possible
//   - Call spin.Yield() periodically in tight loops
//   - Consider using GOMAXPROCS > 1 for concurrent workloads
package zcall
