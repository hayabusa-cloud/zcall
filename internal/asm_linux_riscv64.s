// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux && riscv64

#include "textflag.h"

// func RawSyscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)
//
// RISC-V Linux syscall convention:
//   A7 = syscall number
//   A0 = arg1, A1 = arg2, A2 = arg3, A3 = arg4
//   A0 = return value (negative values indicate -errno)
//
// Stack layout (from FP):
//   num+0(FP), a1+8(FP), a2+16(FP), a3+24(FP), a4+32(FP)
//   r1+40(FP), errno+48(FP)
//
TEXT ·RawSyscall4(SB), NOSPLIT, $0-56
	MOV	num+0(FP), A7
	MOV	a1+8(FP), A0
	MOV	a2+16(FP), A1
	MOV	a3+24(FP), A2
	MOV	a4+32(FP), A3
	ECALL
	MOV	$-4096, X5
	BLTU	X5, A0, err4
	MOV	A0, r1+40(FP)
	MOV	X0, errno+48(FP)
	RET
err4:
	NEG	A0, A0
	MOV	$-1, X5
	MOV	X5, r1+40(FP)
	MOV	A0, errno+48(FP)
	RET

// func RawSyscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
//
// RISC-V Linux syscall convention:
//   A7 = syscall number
//   A0 = arg1, A1 = arg2, A2 = arg3, A3 = arg4, A4 = arg5, A5 = arg6
//   A0 = return value (negative values indicate -errno)
//
// Stack layout (from FP):
//   num+0(FP), a1+8(FP), a2+16(FP), a3+24(FP), a4+32(FP), a5+40(FP), a6+48(FP)
//   r1+56(FP), errno+64(FP)
//
TEXT ·RawSyscall6(SB), NOSPLIT, $0-72
	MOV	num+0(FP), A7
	MOV	a1+8(FP), A0
	MOV	a2+16(FP), A1
	MOV	a3+24(FP), A2
	MOV	a4+32(FP), A3
	MOV	a5+40(FP), A4
	MOV	a6+48(FP), A5
	ECALL
	MOV	$-4096, X5
	BLTU	X5, A0, err6
	MOV	A0, r1+56(FP)
	MOV	X0, errno+64(FP)
	RET
err6:
	NEG	A0, A0
	MOV	$-1, X5
	MOV	X5, r1+56(FP)
	MOV	A0, errno+64(FP)
	RET
