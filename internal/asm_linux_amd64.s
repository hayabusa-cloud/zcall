// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux && amd64

#include "textflag.h"

// func RawSyscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)
//
// AMD64 Linux syscall convention:
//   RAX = syscall number
//   RDI = arg1, RSI = arg2, RDX = arg3, R10 = arg4
//   RAX = return value (negative values indicate -errno)
//
// Using stack-based ABI for compatibility.
// Stack layout (from FP):
//   num+0(FP), a1+8(FP), a2+16(FP), a3+24(FP), a4+32(FP)
//   r1+40(FP), errno+48(FP)
//
TEXT ·RawSyscall4(SB), NOSPLIT, $0-56
	MOVQ num+0(FP), AX
	MOVQ a1+8(FP), DI
	MOVQ a2+16(FP), SI
	MOVQ a3+24(FP), DX
	MOVQ a4+32(FP), R10
	SYSCALL
	CMPQ AX, $-4095
	JLS ok4
	NEGQ AX
	MOVQ $-1, r1+40(FP)
	MOVQ AX, errno+48(FP)
	RET
ok4:
	MOVQ AX, r1+40(FP)
	MOVQ $0, errno+48(FP)
	RET

// func RawSyscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
//
// AMD64 Linux syscall convention:
//   RAX = syscall number
//   RDI = arg1, RSI = arg2, RDX = arg3, R10 = arg4, R8 = arg5, R9 = arg6
//   RAX = return value (negative values indicate -errno)
//
// Using stack-based ABI for compatibility.
// Stack layout (from FP):
//   num+0(FP), a1+8(FP), a2+16(FP), a3+24(FP), a4+32(FP), a5+40(FP), a6+48(FP)
//   r1+56(FP), errno+64(FP)
//
TEXT ·RawSyscall6(SB), NOSPLIT, $0-72
	MOVQ num+0(FP), AX
	MOVQ a1+8(FP), DI
	MOVQ a2+16(FP), SI
	MOVQ a3+24(FP), DX
	MOVQ a4+32(FP), R10
	MOVQ a5+40(FP), R8
	MOVQ a6+48(FP), R9
	SYSCALL
	CMPQ AX, $-4095
	JLS ok6
	NEGQ AX
	MOVQ $-1, r1+56(FP)
	MOVQ AX, errno+64(FP)
	RET
ok6:
	MOVQ AX, r1+56(FP)
	MOVQ $0, errno+64(FP)
	RET
