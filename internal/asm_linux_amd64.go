// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux && amd64

package internal

// RawSyscall4 executes a syscall with up to 4 arguments.
// It bypasses the Go runtime's syscall machinery entirely.
//
// Register mapping (AMD64 Linux ABI):
//   - RAX: syscall number (num)
//   - RDI: arg1 (a1)
//   - RSI: arg2 (a2)
//   - RDX: arg3 (a3)
//   - R10: arg4 (a4)
//
// Returns:
//   - r1: raw syscall return value
//   - errno: 0 on success, or the error number on failure
//
// On Linux, a negative return value in RAX indicates -errno.
// This function extracts the errno and returns it separately.
//
//go:noescape
func RawSyscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// RawSyscall6 executes a syscall with up to 6 arguments.
// It bypasses the Go runtime's syscall machinery entirely.
//
// Register mapping (AMD64 Linux ABI):
//   - RAX: syscall number (num)
//   - RDI: arg1 (a1)
//   - RSI: arg2 (a2)
//   - RDX: arg3 (a3)
//   - R10: arg4 (a4)
//   - R8:  arg5 (a5)
//   - R9:  arg6 (a6)
//
// Returns:
//   - r1: raw syscall return value
//   - errno: 0 on success, or the error number on failure
//
// On Linux, a negative return value in RAX indicates -errno.
// This function extracts the errno and returns it separately.
//
//go:noescape
func RawSyscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
