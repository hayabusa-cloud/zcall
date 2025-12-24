// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux

package zcall

// Common Linux error numbers.
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
	EAGAIN          Errno = 11
	EWOULDBLOCK     Errno = EAGAIN
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
	EDEADLK         Errno = 35
	ENAMETOOLONG    Errno = 36
	ENOLCK          Errno = 37
	ENOSYS          Errno = 38
	ENOTEMPTY       Errno = 39
	ELOOP           Errno = 40
	ENOMSG          Errno = 42
	EIDRM           Errno = 43
	ECHRNG          Errno = 44
	EL2NSYNC        Errno = 45
	EL3HLT          Errno = 46
	EL3RST          Errno = 47
	ELNRNG          Errno = 48
	EUNATCH         Errno = 49
	ENOCSI          Errno = 50
	EL2HLT          Errno = 51
	EBADE           Errno = 52
	EBADR           Errno = 53
	EXFULL          Errno = 54
	ENOANO          Errno = 55
	EBADRQC         Errno = 56
	EBADSLT         Errno = 57
	EBFONT          Errno = 59
	ENOSTR          Errno = 60
	ENODATA         Errno = 61
	ETIME           Errno = 62
	ENOSR           Errno = 63
	ENONET          Errno = 64
	ENOPKG          Errno = 65
	EREMOTE         Errno = 66
	ENOLINK         Errno = 67
	EADV            Errno = 68
	ESRMNT          Errno = 69
	ECOMM           Errno = 70
	EPROTO          Errno = 71
	EMULTIHOP       Errno = 72
	EDOTDOT         Errno = 73
	EBADMSG         Errno = 74
	EOVERFLOW       Errno = 75
	ENOTUNIQ        Errno = 76
	EBADFD          Errno = 77
	EREMCHG         Errno = 78
	ELIBACC         Errno = 79
	ELIBBAD         Errno = 80
	ELIBSCN         Errno = 81
	ELIBMAX         Errno = 82
	ELIBEXEC        Errno = 83
	EILSEQ          Errno = 84
	ERESTART        Errno = 85
	ESTRPIPE        Errno = 86
	EUSERS          Errno = 87
	ENOTSOCK        Errno = 88
	EDESTADDRREQ    Errno = 89
	EMSGSIZE        Errno = 90
	EPROTOTYPE      Errno = 91
	ENOPROTOOPT     Errno = 92
	EPROTONOSUPPORT Errno = 93
	ESOCKTNOSUPPORT Errno = 94
	EOPNOTSUPP      Errno = 95
	ENOTSUP         Errno = EOPNOTSUPP
	EPFNOSUPPORT    Errno = 96
	EAFNOSUPPORT    Errno = 97
	EADDRINUSE      Errno = 98
	EADDRNOTAVAIL   Errno = 99
	ENETDOWN        Errno = 100
	ENETUNREACH     Errno = 101
	ENETRESET       Errno = 102
	ECONNABORTED    Errno = 103
	ECONNRESET      Errno = 104
	ENOBUFS         Errno = 105
	EISCONN         Errno = 106
	ENOTCONN        Errno = 107
	ESHUTDOWN       Errno = 108
	ETOOMANYREFS    Errno = 109
	ETIMEDOUT       Errno = 110
	ECONNREFUSED    Errno = 111
	EHOSTDOWN       Errno = 112
	EHOSTUNREACH    Errno = 113
	EALREADY        Errno = 114
	EINPROGRESS     Errno = 115
	ESTALE          Errno = 116
	EUCLEAN         Errno = 117
	ENOTNAM         Errno = 118
	ENAVAIL         Errno = 119
	EISNAM          Errno = 120
	EREMOTEIO       Errno = 121
	EDQUOT          Errno = 122
	ENOMEDIUM       Errno = 123
	EMEDIUMTYPE     Errno = 124
	ECANCELED       Errno = 125
	ENOKEY          Errno = 126
	EKEYEXPIRED     Errno = 127
	EKEYREVOKED     Errno = 128
	EKEYREJECTED    Errno = 129
	EOWNERDEAD      Errno = 130
	ENOTRECOVERABLE Errno = 131
	ERFKILL         Errno = 132
	EHWPOISON       Errno = 133
)

// errnoStrings maps errno values to their string representations.
var errnoStrings = [...]string{
	1:   "operation not permitted",
	2:   "no such file or directory",
	3:   "no such process",
	4:   "interrupted system call",
	5:   "input/output error",
	6:   "no such device or address",
	7:   "argument list too long",
	8:   "exec format error",
	9:   "bad file descriptor",
	10:  "no child processes",
	11:  "resource temporarily unavailable",
	12:  "cannot allocate memory",
	13:  "permission denied",
	14:  "bad address",
	15:  "block device required",
	16:  "device or resource busy",
	17:  "file exists",
	18:  "invalid cross-device link",
	19:  "no such device",
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
	34:  "numerical result out of range",
	35:  "resource deadlock avoided",
	36:  "file name too long",
	37:  "no locks available",
	38:  "function not implemented",
	39:  "directory not empty",
	40:  "too many levels of symbolic links",
	42:  "no message of desired type",
	43:  "identifier removed",
	44:  "channel number out of range",
	45:  "level 2 not synchronized",
	46:  "level 3 halted",
	47:  "level 3 reset",
	48:  "link number out of range",
	49:  "protocol driver not attached",
	50:  "no CSI structure available",
	51:  "level 2 halted",
	52:  "invalid exchange",
	53:  "invalid request descriptor",
	54:  "exchange full",
	55:  "no anode",
	56:  "invalid request code",
	57:  "invalid slot",
	59:  "bad font file format",
	60:  "device not a stream",
	61:  "no data available",
	62:  "timer expired",
	63:  "out of streams resources",
	64:  "machine is not on the network",
	65:  "package not installed",
	66:  "object is remote",
	67:  "link has been severed",
	68:  "advertise error",
	69:  "srmount error",
	70:  "communication error on send",
	71:  "protocol error",
	72:  "multihop attempted",
	73:  "RFS specific error",
	74:  "bad message",
	75:  "value too large for defined data type",
	76:  "name not unique on network",
	77:  "file descriptor in bad state",
	78:  "remote address changed",
	79:  "can not access a needed shared library",
	80:  "accessing a corrupted shared library",
	81:  ".lib section in a.out corrupted",
	82:  "attempting to link in too many shared libraries",
	83:  "cannot exec a shared library directly",
	84:  "invalid or incomplete multibyte or wide character",
	85:  "interrupted system call should be restarted",
	86:  "streams pipe error",
	87:  "too many users",
	88:  "socket operation on non-socket",
	89:  "destination address required",
	90:  "message too long",
	91:  "protocol wrong type for socket",
	92:  "protocol not available",
	93:  "protocol not supported",
	94:  "socket type not supported",
	95:  "operation not supported",
	96:  "protocol family not supported",
	97:  "address family not supported by protocol",
	98:  "address already in use",
	99:  "cannot assign requested address",
	100: "network is down",
	101: "network is unreachable",
	102: "network dropped connection on reset",
	103: "software caused connection abort",
	104: "connection reset by peer",
	105: "no buffer space available",
	106: "transport endpoint is already connected",
	107: "transport endpoint is not connected",
	108: "cannot send after transport endpoint shutdown",
	109: "too many references: cannot splice",
	110: "connection timed out",
	111: "connection refused",
	112: "host is down",
	113: "no route to host",
	114: "operation already in progress",
	115: "operation now in progress",
	116: "stale file handle",
	117: "structure needs cleaning",
	118: "not a XENIX named type file",
	119: "no XENIX semaphores available",
	120: "is a named type file",
	121: "remote I/O error",
	122: "disk quota exceeded",
	123: "no medium found",
	124: "wrong medium type",
	125: "operation canceled",
	126: "required key not available",
	127: "key has expired",
	128: "key has been revoked",
	129: "key was rejected by service",
	130: "owner died",
	131: "state not recoverable",
	132: "operation not possible due to RF-kill",
	133: "memory page has hardware error",
}
