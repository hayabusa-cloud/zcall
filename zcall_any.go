// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package zcall

import "unsafe"

// noescape hides a pointer from escape analysis.
//
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

// Errno is a raw system error number.
// It implements the error interface and provides helper methods
// for common error patterns.
type Errno uintptr

func (e Errno) Error() string {
	if e == 0 {
		return "success"
	}
	if int(e) >= 0 && int(e) < len(errnoStrings) {
		s := errnoStrings[e]
		if s != "" {
			return s
		}
	}
	return "errno " + uitoa(uint(e))
}

// Is reports whether the error is equal to the target.
// This enables errors.Is() compatibility.
func (e Errno) Is(target error) bool {
	t, ok := target.(Errno)
	return ok && e == t
}

// Temporary reports whether the error is temporary.
// Temporary errors include EAGAIN, EWOULDBLOCK, EINTR, and EINPROGRESS.
func (e Errno) Temporary() bool {
	return e == EAGAIN || e == EWOULDBLOCK || e == EINTR || e == EINPROGRESS
}

// Timeout reports whether the error represents a timeout.
func (e Errno) Timeout() bool {
	return e == EAGAIN || e == EWOULDBLOCK || e == ETIMEDOUT
}

// uitoa converts an unsigned integer to a string without allocations.
func uitoa(val uint) string {
	if val == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf) - 1
	for val > 0 {
		buf[i] = byte('0' + val%10)
		val /= 10
		i--
	}
	return string(buf[i+1:])
}
