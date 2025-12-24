// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux && arm64

#include "textflag.h"

// func RawSyscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)
//
// ARM64 Linux syscall convention:
//   X8 = syscall number
//   X0 = arg1, X1 = arg2, X2 = arg3, X3 = arg4
//   X0 = return value (negative values indicate -errno)
//
// Using stack-based ABI for compatibility.
// Stack layout (from FP):
//   num+0(FP), a1+8(FP), a2+16(FP), a3+24(FP), a4+32(FP)
//   r1+40(FP), errno+48(FP)
//
TEXT ·RawSyscall4(SB), NOSPLIT, $0-56
	MOVD num+0(FP), R8
	MOVD a1+8(FP), R0
	MOVD a2+16(FP), R1
	MOVD a3+24(FP), R2
	MOVD a4+32(FP), R3
	SVC $0
	CMN $4095, R0
	BLS ok4
	NEG R0, R0
	MOVD $-1, R1
	MOVD R1, r1+40(FP)
	MOVD R0, errno+48(FP)
	RET
ok4:
	MOVD R0, r1+40(FP)
	MOVD $0, R1
	MOVD R1, errno+48(FP)
	RET

// func RawSyscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
//
// ARM64 Linux syscall convention:
//   X8 = syscall number
//   X0 = arg1, X1 = arg2, X2 = arg3, X3 = arg4, X4 = arg5, X5 = arg6
//   X0 = return value (negative values indicate -errno)
//
// Using stack-based ABI for compatibility.
// Stack layout (from FP):
//   num+0(FP), a1+8(FP), a2+16(FP), a3+24(FP), a4+32(FP), a5+40(FP), a6+48(FP)
//   r1+56(FP), errno+64(FP)
//
TEXT ·RawSyscall6(SB), NOSPLIT, $0-72
	MOVD num+0(FP), R8
	MOVD a1+8(FP), R0
	MOVD a2+16(FP), R1
	MOVD a3+24(FP), R2
	MOVD a4+32(FP), R3
	MOVD a5+40(FP), R4
	MOVD a6+48(FP), R5
	SVC $0
	CMN $4095, R0
	BLS ok6
	NEG R0, R0
	MOVD $-1, R1
	MOVD R1, r1+56(FP)
	MOVD R0, errno+64(FP)
	RET
ok6:
	MOVD R0, r1+56(FP)
	MOVD $0, R1
	MOVD R1, errno+64(FP)
	RET
