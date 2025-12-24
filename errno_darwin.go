// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build darwin

package zcall

// Common Darwin/BSD errno values.
// Reference: /usr/include/sys/errno.h
const (
	EPERM           Errno = 1
	ENOENT          Errno = 2
	ESRCH           Errno = 3
	EINTR           Errno = 4
	EIO             Errno = 5
	ENXIO           Errno = 6
	E2BIG           Errno = 7
	ENOEXEC         Errno = 8
	EBADF           Errno = 9
	ECHILD          Errno = 10
	EDEADLK         Errno = 11
	ENOMEM          Errno = 12
	EACCES          Errno = 13
	EFAULT          Errno = 14
	ENOTBLK         Errno = 15
	EBUSY           Errno = 16
	EEXIST          Errno = 17
	EXDEV           Errno = 18
	ENODEV          Errno = 19
	ENOTDIR         Errno = 20
	EISDIR          Errno = 21
	EINVAL          Errno = 22
	ENFILE          Errno = 23
	EMFILE          Errno = 24
	ENOTTY          Errno = 25
	ETXTBSY         Errno = 26
	EFBIG           Errno = 27
	ENOSPC          Errno = 28
	ESPIPE          Errno = 29
	EROFS           Errno = 30
	EMLINK          Errno = 31
	EPIPE           Errno = 32
	EDOM            Errno = 33
	ERANGE          Errno = 34
	EAGAIN          Errno = 35
	EWOULDBLOCK     Errno = EAGAIN
	EINPROGRESS     Errno = 36
	EALREADY        Errno = 37
	ENOTSOCK        Errno = 38
	EDESTADDRREQ    Errno = 39
	EMSGSIZE        Errno = 40
	EPROTOTYPE      Errno = 41
	ENOPROTOOPT     Errno = 42
	EPROTONOSUPPORT Errno = 43
	ESOCKTNOSUPPORT Errno = 44
	ENOTSUP         Errno = 45
	EOPNOTSUPP      Errno = ENOTSUP
	EPFNOSUPPORT    Errno = 46
	EAFNOSUPPORT    Errno = 47
	EADDRINUSE      Errno = 48
	EADDRNOTAVAIL   Errno = 49
	ENETDOWN        Errno = 50
	ENETUNREACH     Errno = 51
	ENETRESET       Errno = 52
	ECONNABORTED    Errno = 53
	ECONNRESET      Errno = 54
	ENOBUFS         Errno = 55
	EISCONN         Errno = 56
	ENOTCONN        Errno = 57
	ESHUTDOWN       Errno = 58
	ETOOMANYREFS    Errno = 59
	ETIMEDOUT       Errno = 60
	ECONNREFUSED    Errno = 61
	ELOOP           Errno = 62
	ENAMETOOLONG    Errno = 63
	EHOSTDOWN       Errno = 64
	EHOSTUNREACH    Errno = 65
	ENOTEMPTY       Errno = 66
	EPROCLIM        Errno = 67
	EUSERS          Errno = 68
	EDQUOT          Errno = 69
	ESTALE          Errno = 70
	EREMOTE         Errno = 71
	EBADRPC         Errno = 72
	ERPCMISMATCH    Errno = 73
	EPROGUNAVAIL    Errno = 74
	EPROGMISMATCH   Errno = 75
	EPROCUNAVAIL    Errno = 76
	ENOLCK          Errno = 77
	ENOSYS          Errno = 78
	EFTYPE          Errno = 79
	EAUTH           Errno = 80
	ENEEDAUTH       Errno = 81
	EPWROFF         Errno = 82
	EDEVERR         Errno = 83
	EOVERFLOW       Errno = 84
	EBADEXEC        Errno = 85
	EBADARCH        Errno = 86
	ESHLIBVERS      Errno = 87
	EBADMACHO       Errno = 88
	ECANCELED       Errno = 89
	EIDRM           Errno = 90
	ENOMSG          Errno = 91
	EILSEQ          Errno = 92
	ENOATTR         Errno = 93
	EBADMSG         Errno = 94
	EMULTIHOP       Errno = 95
	ENODATA         Errno = 96
	ENOLINK         Errno = 97
	ENOSR           Errno = 98
	ENOSTR          Errno = 99
	EPROTO          Errno = 100
	ETIME           Errno = 101
	EOPNOTSUPP_ALT  Errno = 102
	ENOPOLICY       Errno = 103
	ENOTRECOVERABLE Errno = 104
	EOWNERDEAD      Errno = 105
	EQFULL          Errno = 106
)

var errnoStrings = [...]string{
	1:   "operation not permitted",
	2:   "no such file or directory",
	3:   "no such process",
	4:   "interrupted system call",
	5:   "input/output error",
	6:   "device not configured",
	7:   "argument list too long",
	8:   "exec format error",
	9:   "bad file descriptor",
	10:  "no child processes",
	11:  "resource deadlock avoided",
	12:  "cannot allocate memory",
	13:  "permission denied",
	14:  "bad address",
	15:  "block device required",
	16:  "resource busy",
	17:  "file exists",
	18:  "cross-device link",
	19:  "operation not supported by device",
	20:  "not a directory",
	21:  "is a directory",
	22:  "invalid argument",
	23:  "too many open files in system",
	24:  "too many open files",
	25:  "inappropriate ioctl for device",
	26:  "text file busy",
	27:  "file too large",
	28:  "no space left on device",
	29:  "illegal seek",
	30:  "read-only file system",
	31:  "too many links",
	32:  "broken pipe",
	33:  "numerical argument out of domain",
	34:  "result too large",
	35:  "resource temporarily unavailable",
	36:  "operation now in progress",
	37:  "operation already in progress",
	38:  "socket operation on non-socket",
	39:  "destination address required",
	40:  "message too long",
	41:  "protocol wrong type for socket",
	42:  "protocol not available",
	43:  "protocol not supported",
	44:  "socket type not supported",
	45:  "operation not supported",
	46:  "protocol family not supported",
	47:  "address family not supported by protocol family",
	48:  "address already in use",
	49:  "can't assign requested address",
	50:  "network is down",
	51:  "network is unreachable",
	52:  "network dropped connection on reset",
	53:  "software caused connection abort",
	54:  "connection reset by peer",
	55:  "no buffer space available",
	56:  "socket is already connected",
	57:  "socket is not connected",
	58:  "can't send after socket shutdown",
	59:  "too many references: can't splice",
	60:  "operation timed out",
	61:  "connection refused",
	62:  "too many levels of symbolic links",
	63:  "file name too long",
	64:  "host is down",
	65:  "no route to host",
	66:  "directory not empty",
	67:  "too many processes",
	68:  "too many users",
	69:  "disc quota exceeded",
	70:  "stale NFS file handle",
	71:  "too many levels of remote in path",
	72:  "RPC struct is bad",
	73:  "RPC version wrong",
	74:  "RPC prog. not avail",
	75:  "program version wrong",
	76:  "bad procedure for program",
	77:  "no locks available",
	78:  "function not implemented",
	79:  "inappropriate file type or format",
	80:  "authentication error",
	81:  "need authenticator",
	82:  "device power is off",
	83:  "device error",
	84:  "value too large to be stored in data type",
	85:  "bad executable",
	86:  "bad CPU type in executable",
	87:  "shared library version mismatch",
	88:  "malformed Mach-o file",
	89:  "operation canceled",
	90:  "identifier removed",
	91:  "no message of desired type",
	92:  "illegal byte sequence",
	93:  "attribute not found",
	94:  "bad message",
	95:  "reserved",
	96:  "no message available on STREAM",
	97:  "reserved",
	98:  "no STREAM resources",
	99:  "not a STREAM",
	100: "protocol error",
	101: "STREAM ioctl timeout",
	102: "operation not supported on socket",
	103: "policy not found",
	104: "state not recoverable",
	105: "previous owner died",
	106: "interface output queue is full",
}
