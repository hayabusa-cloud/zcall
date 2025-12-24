// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux && loong64

#include "textflag.h"

// func RawSyscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)
//
// LoongArch Linux syscall convention:
//   R11 (A7) = syscall number
//   R4 (A0) = arg1, R5 (A1) = arg2, R6 (A2) = arg3, R7 (A3) = arg4
//   R4 (A0) = return value (negative values indicate -errno)
//
// Stack layout (from FP):
//   num+0(FP), a1+8(FP), a2+16(FP), a3+24(FP), a4+32(FP)
//   r1+40(FP), errno+48(FP)
//
TEXT ·RawSyscall4(SB), NOSPLIT, $0-56
	MOVV	num+0(FP), R11
	MOVV	a1+8(FP), R4
	MOVV	a2+16(FP), R5
	MOVV	a3+24(FP), R6
	MOVV	a4+32(FP), R7
	SYSCALL
	MOVV	$-4096, R12
	BGEU	R12, R4, ok4
	SUBV	R0, R4, R4
	MOVV	$-1, R12
	MOVV	R12, r1+40(FP)
	MOVV	R4, errno+48(FP)
	RET
ok4:
	MOVV	R4, r1+40(FP)
	MOVV	R0, errno+48(FP)
	RET

// func RawSyscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
//
// LoongArch Linux syscall convention:
//   R11 (A7) = syscall number
//   R4 (A0) = arg1, R5 (A1) = arg2, R6 (A2) = arg3, R7 (A3) = arg4
//   R8 (A4) = arg5, R9 (A5) = arg6
//   R4 (A0) = return value (negative values indicate -errno)
//
// Stack layout (from FP):
//   num+0(FP), a1+8(FP), a2+16(FP), a3+24(FP), a4+32(FP), a5+40(FP), a6+48(FP)
//   r1+56(FP), errno+64(FP)
//
TEXT ·RawSyscall6(SB), NOSPLIT, $0-72
	MOVV	num+0(FP), R11
	MOVV	a1+8(FP), R4
	MOVV	a2+16(FP), R5
	MOVV	a3+24(FP), R6
	MOVV	a4+32(FP), R7
	MOVV	a5+40(FP), R8
	MOVV	a6+48(FP), R9
	SYSCALL
	MOVV	$-4096, R12
	BGEU	R12, R4, ok6
	SUBV	R0, R4, R4
	MOVV	$-1, R12
	MOVV	R12, r1+56(FP)
	MOVV	R4, errno+64(FP)
	RET
ok6:
	MOVV	R4, r1+56(FP)
	MOVV	R0, errno+64(FP)
	RET
