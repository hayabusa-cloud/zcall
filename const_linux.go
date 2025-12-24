// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux

package zcall

// Syscall numbers for Linux.
// These are architecture-specific and defined in const_linux_*.go files.
// The constants below are common definitions used across architectures.

// Socket address families.
const (
	AF_UNIX   = 1
	AF_LOCAL  = AF_UNIX
	AF_INET   = 2
	AF_INET6  = 10
	AF_PACKET = 17
)

// Socket types.
const (
	SOCK_STREAM    = 1
	SOCK_DGRAM     = 2
	SOCK_RAW       = 3
	SOCK_RDM       = 4
	SOCK_SEQPACKET = 5
	SOCK_NONBLOCK  = 0x800
	SOCK_CLOEXEC   = 0x80000
)

// IP protocols.
const (
	IPPROTO_IP   = 0
	IPPROTO_ICMP = 1
	IPPROTO_TCP  = 6
	IPPROTO_UDP  = 17
	IPPROTO_IPV6 = 41
	IPPROTO_SCTP = 132
	IPPROTO_RAW  = 255
)

// Socket options levels.
const (
	SOL_SOCKET = 1
	SOL_IP     = 0
	SOL_TCP    = 6
	SOL_UDP    = 17
	SOL_IPV6   = 41
	SOL_SCTP   = 132
)

// Socket options (SOL_SOCKET level).
const (
	SO_DEBUG        = 1
	SO_REUSEADDR    = 2
	SO_TYPE         = 3
	SO_ERROR        = 4
	SO_DONTROUTE    = 5
	SO_BROADCAST    = 6
	SO_SNDBUF       = 7
	SO_RCVBUF       = 8
	SO_KEEPALIVE    = 9
	SO_OOBINLINE    = 10
	SO_NO_CHECK     = 11
	SO_PRIORITY     = 12
	SO_LINGER       = 13
	SO_BSDCOMPAT    = 14
	SO_REUSEPORT    = 15
	SO_RCVLOWAT     = 18
	SO_SNDLOWAT     = 19
	SO_RCVTIMEO     = 20
	SO_SNDTIMEO     = 21
	SO_ACCEPTCONN   = 30
	SO_SNDBUFFORCE  = 32
	SO_RCVBUFFORCE  = 33
	SO_PROTOCOL     = 38
	SO_DOMAIN       = 39
	SO_ZEROCOPY     = 60
	SO_INCOMING_CPU = 49
	SO_BUSY_POLL    = 46
)

// TCP options.
const (
	TCP_NODELAY       = 1
	TCP_MAXSEG        = 2
	TCP_CORK          = 3
	TCP_KEEPIDLE      = 4
	TCP_KEEPINTVL     = 5
	TCP_KEEPCNT       = 6
	TCP_SYNCNT        = 7
	TCP_LINGER2       = 8
	TCP_DEFER_ACCEPT  = 9
	TCP_WINDOW_CLAMP  = 10
	TCP_INFO          = 11
	TCP_QUICKACK      = 12
	TCP_CONGESTION    = 13
	TCP_FASTOPEN      = 23
	TCP_NOTSENT_LOWAT = 25
)

// UDP options.
const (
	UDP_CORK    = 1
	UDP_ENCAP   = 100
	UDP_SEGMENT = 103
	UDP_GRO     = 104
)

// IPv6 options.
const (
	IPV6_V6ONLY       = 26
	IPV6_RECVPKTINFO  = 49
	IPV6_PKTINFO      = 50
	IPV6_RECVHOPLIMIT = 51
	IPV6_HOPLIMIT     = 52
)

// File descriptor flags.
const (
	O_RDONLY   = 0x0
	O_WRONLY   = 0x1
	O_RDWR     = 0x2
	O_CREAT    = 0x40
	O_EXCL     = 0x80
	O_NOCTTY   = 0x100
	O_TRUNC    = 0x200
	O_APPEND   = 0x400
	O_NONBLOCK = 0x800
	O_SYNC     = 0x101000
	O_CLOEXEC  = 0x80000
	O_DIRECT   = 0x4000
)

// eventfd flags.
const (
	EFD_SEMAPHORE = 0x1
	EFD_CLOEXEC   = 0x80000
	EFD_NONBLOCK  = 0x800
)

// timerfd flags.
const (
	TFD_CLOEXEC       = 0x80000
	TFD_NONBLOCK      = 0x800
	TFD_TIMER_ABSTIME = 0x1
)

// signalfd flags.
const (
	SFD_CLOEXEC  = 0x80000
	SFD_NONBLOCK = 0x800
)

// pidfd flags.
const (
	PIDFD_NONBLOCK = 0x800
)

// memfd flags.
const (
	MFD_CLOEXEC       = 0x1
	MFD_ALLOW_SEALING = 0x2
	MFD_HUGETLB       = 0x4
	MFD_NOEXEC_SEAL   = 0x8
	MFD_EXEC          = 0x10
)

// Clock IDs for timerfd.
const (
	CLOCK_REALTIME  = 0
	CLOCK_MONOTONIC = 1
)

// io_uring setup flags.
const (
	IORING_SETUP_IOPOLL        = 1 << 0
	IORING_SETUP_SQPOLL        = 1 << 1
	IORING_SETUP_SQ_AFF        = 1 << 2
	IORING_SETUP_CQSIZE        = 1 << 3
	IORING_SETUP_CLAMP         = 1 << 4
	IORING_SETUP_ATTACH_WQ     = 1 << 5
	IORING_SETUP_R_DISABLED    = 1 << 6
	IORING_SETUP_SUBMIT_ALL    = 1 << 7
	IORING_SETUP_COOP_TASKRUN  = 1 << 8
	IORING_SETUP_TASKRUN_FLAG  = 1 << 9
	IORING_SETUP_SQE128        = 1 << 10
	IORING_SETUP_CQE32         = 1 << 11
	IORING_SETUP_SINGLE_ISSUER = 1 << 12
	IORING_SETUP_DEFER_TASKRUN = 1 << 13
)

// io_uring enter flags.
const (
	IORING_ENTER_GETEVENTS       = 1 << 0
	IORING_ENTER_SQ_WAKEUP       = 1 << 1
	IORING_ENTER_SQ_WAIT         = 1 << 2
	IORING_ENTER_EXT_ARG         = 1 << 3
	IORING_ENTER_REGISTERED_RING = 1 << 4
)

// io_uring opcodes.
const (
	IORING_OP_NOP             = 0
	IORING_OP_READV           = 1
	IORING_OP_WRITEV          = 2
	IORING_OP_FSYNC           = 3
	IORING_OP_READ_FIXED      = 4
	IORING_OP_WRITE_FIXED     = 5
	IORING_OP_POLL_ADD        = 6
	IORING_OP_POLL_REMOVE     = 7
	IORING_OP_SYNC_FILE_RANGE = 8
	IORING_OP_SENDMSG         = 9
	IORING_OP_RECVMSG         = 10
	IORING_OP_TIMEOUT         = 11
	IORING_OP_TIMEOUT_REMOVE  = 12
	IORING_OP_ACCEPT          = 13
	IORING_OP_ASYNC_CANCEL    = 14
	IORING_OP_LINK_TIMEOUT    = 15
	IORING_OP_CONNECT         = 16
	IORING_OP_FALLOCATE       = 17
	IORING_OP_OPENAT          = 18
	IORING_OP_CLOSE           = 19
	IORING_OP_FILES_UPDATE    = 20
	IORING_OP_STATX           = 21
	IORING_OP_READ            = 22
	IORING_OP_WRITE           = 23
	IORING_OP_FADVISE         = 24
	IORING_OP_MADVISE         = 25
	IORING_OP_SEND            = 26
	IORING_OP_RECV            = 27
	IORING_OP_OPENAT2         = 28
	IORING_OP_EPOLL_CTL       = 29
	IORING_OP_SPLICE          = 30
	IORING_OP_PROVIDE_BUFFERS = 31
	IORING_OP_REMOVE_BUFFERS  = 32
	IORING_OP_TEE             = 33
	IORING_OP_SHUTDOWN        = 34
	IORING_OP_RENAMEAT        = 35
	IORING_OP_UNLINKAT        = 36
	IORING_OP_MKDIRAT         = 37
	IORING_OP_SYMLINKAT       = 38
	IORING_OP_LINKAT          = 39
	IORING_OP_MSG_RING        = 40
	IORING_OP_FSETXATTR       = 41
	IORING_OP_SETXATTR        = 42
	IORING_OP_FGETXATTR       = 43
	IORING_OP_GETXATTR        = 44
	IORING_OP_SOCKET          = 45
	IORING_OP_URING_CMD       = 46
	IORING_OP_SEND_ZC         = 47
	IORING_OP_SENDMSG_ZC      = 48
)

// io_uring SQE flags.
const (
	IOSQE_FIXED_FILE       = 1 << 0
	IOSQE_IO_DRAIN         = 1 << 1
	IOSQE_IO_LINK          = 1 << 2
	IOSQE_IO_HARDLINK      = 1 << 3
	IOSQE_ASYNC            = 1 << 4
	IOSQE_BUFFER_SELECT    = 1 << 5
	IOSQE_CQE_SKIP_SUCCESS = 1 << 6
)

// io_uring register opcodes.
const (
	IORING_REGISTER_BUFFERS          = 0
	IORING_UNREGISTER_BUFFERS        = 1
	IORING_REGISTER_FILES            = 2
	IORING_UNREGISTER_FILES          = 3
	IORING_REGISTER_EVENTFD          = 4
	IORING_UNREGISTER_EVENTFD        = 5
	IORING_REGISTER_FILES_UPDATE     = 6
	IORING_REGISTER_EVENTFD_ASYNC    = 7
	IORING_REGISTER_PROBE            = 8
	IORING_REGISTER_PERSONALITY      = 9
	IORING_UNREGISTER_PERSONALITY    = 10
	IORING_REGISTER_RESTRICTIONS     = 11
	IORING_REGISTER_ENABLE_RINGS     = 12
	IORING_REGISTER_FILES2           = 13
	IORING_REGISTER_FILES_UPDATE2    = 14
	IORING_REGISTER_BUFFERS2         = 15
	IORING_REGISTER_BUFFERS_UPDATE   = 16
	IORING_REGISTER_IOWQ_AFF         = 17
	IORING_UNREGISTER_IOWQ_AFF       = 18
	IORING_REGISTER_IOWQ_MAX_WORKERS = 19
	IORING_REGISTER_RING_FDS         = 20
	IORING_UNREGISTER_RING_FDS       = 21
	IORING_REGISTER_PBUF_RING        = 22
	IORING_UNREGISTER_PBUF_RING      = 23
	IORING_REGISTER_SYNC_CANCEL      = 24
	IORING_REGISTER_FILE_ALLOC_RANGE = 25
)

// mmap protection flags.
const (
	PROT_NONE  = 0x0
	PROT_READ  = 0x1
	PROT_WRITE = 0x2
	PROT_EXEC  = 0x4
)

// mmap flags.
const (
	MAP_SHARED    = 0x1
	MAP_PRIVATE   = 0x2
	MAP_FIXED     = 0x10
	MAP_ANONYMOUS = 0x20
	MAP_POPULATE  = 0x8000
)

// Poll events.
const (
	POLLIN    = 0x1
	POLLPRI   = 0x2
	POLLOUT   = 0x4
	POLLERR   = 0x8
	POLLHUP   = 0x10
	POLLNVAL  = 0x20
	POLLRDHUP = 0x2000
)

// Shutdown how.
const (
	SHUT_RD   = 0
	SHUT_WR   = 1
	SHUT_RDWR = 2
)

// MSG flags for send/recv.
const (
	MSG_OOB          = 0x1
	MSG_PEEK         = 0x2
	MSG_DONTROUTE    = 0x4
	MSG_CTRUNC       = 0x8
	MSG_PROXY        = 0x10
	MSG_TRUNC        = 0x20
	MSG_DONTWAIT     = 0x40
	MSG_EOR          = 0x80
	MSG_WAITALL      = 0x100
	MSG_FIN          = 0x200
	MSG_SYN          = 0x400
	MSG_CONFIRM      = 0x800
	MSG_RST          = 0x1000
	MSG_ERRQUEUE     = 0x2000
	MSG_NOSIGNAL     = 0x4000
	MSG_MORE         = 0x8000
	MSG_WAITFORONE   = 0x10000
	MSG_BATCH        = 0x40000
	MSG_ZEROCOPY     = 0x4000000
	MSG_FASTOPEN     = 0x20000000
	MSG_CMSG_CLOEXEC = 0x40000000
)

// Splice flags.
const (
	SPLICE_F_MOVE     = 0x1
	SPLICE_F_NONBLOCK = 0x2
	SPLICE_F_MORE     = 0x4
	SPLICE_F_GIFT     = 0x8
)

// RWF flags for preadv2/pwritev2.
const (
	RWF_HIPRI  = 0x1
	RWF_DSYNC  = 0x2
	RWF_SYNC   = 0x4
	RWF_NOWAIT = 0x8
	RWF_APPEND = 0x10
)

// Iovec represents a scatter/gather I/O vector.
// Used by readv, writev, preadv, pwritev, and related syscalls.
type Iovec struct {
	Base *byte
	Len  uint64
}

// Msghdr represents a message header for sendmsg/recvmsg.
type Msghdr struct {
	Name       *byte
	Namelen    uint32
	_          [4]byte // padding on 64-bit
	Iov        *Iovec
	Iovlen     uint64
	Control    *byte
	Controllen uint64
	Flags      int32
	_          [4]byte // padding on 64-bit
}

// Mmsghdr represents a message header for sendmmsg/recvmmsg.
type Mmsghdr struct {
	Hdr Msghdr
	Len uint32
	_   [4]byte // padding on 64-bit
}

// Timespec represents a time value with nanosecond precision.
type Timespec struct {
	Sec  int64
	Nsec int64
}

// Itimerspec represents an interval timer specification.
type Itimerspec struct {
	Interval Timespec
	Value    Timespec
}
