// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux

package zcall_test

import (
	"errors"
	"net"
	"os"
	"testing"
	"time"
	"unsafe"
	_ "unsafe" // for go:linkname

	"code.hybscloud.com/zcall"
)

func TestIfNameToIndex(t *testing.T) {
	iface, err := net.InterfaceByName("lo")
	if err != nil {
		t.Fatalf("net.InterfaceByName(lo) failed: %v", err)
	}

	index, err := zcall.IfNameToIndex("lo")
	if err != nil {
		t.Fatalf("IfNameToIndex(lo) failed: %v", err)
	}
	if index != uint32(iface.Index) {
		t.Fatalf("IfNameToIndex(lo)=%d, want %d", index, iface.Index)
	}
}

func TestLinksAndLinkLookups(t *testing.T) {
	iface, err := net.InterfaceByName("lo")
	if err != nil {
		t.Fatalf("net.InterfaceByName(lo) failed: %v", err)
	}

	links, err := zcall.Links()
	if err != nil {
		t.Fatalf("Links() failed: %v", err)
	}
	if len(links) == 0 {
		t.Fatal("Links() returned no links")
	}

	var found bool
	for _, link := range links {
		if link.Name == "lo" {
			found = true
			if link.Index != iface.Index {
				t.Fatalf("Links()[lo].Index=%d, want %d", link.Index, iface.Index)
			}
			break
		}
	}
	if !found {
		t.Fatal("Links() did not include loopback")
	}

	byName, err := zcall.LinkByName("lo")
	if err != nil {
		t.Fatalf("LinkByName(lo) failed: %v", err)
	}
	if byName.Index != iface.Index {
		t.Fatalf("LinkByName(lo).Index=%d, want %d", byName.Index, iface.Index)
	}

	byIndex, err := zcall.LinkByIndex(iface.Index)
	if err != nil {
		t.Fatalf("LinkByIndex(%d) failed: %v", iface.Index, err)
	}
	if byIndex.Name != iface.Name {
		t.Fatalf("LinkByIndex(%d).Name=%q, want %q", iface.Index, byIndex.Name, iface.Name)
	}
}

func TestIfNameToIndexRejectsInvalidName(t *testing.T) {
	if _, err := zcall.IfNameToIndex(""); !errors.Is(err, zcall.Errno(zcall.EINVAL)) {
		t.Fatalf("IfNameToIndex(%q) error = %v, want EINVAL", "", err)
	}

	tooLong := "0123456789abcdef"
	if _, err := zcall.IfNameToIndex(tooLong); !errors.Is(err, zcall.Errno(zcall.EINVAL)) {
		t.Fatalf("IfNameToIndex(%q) error = %v, want EINVAL", tooLong, err)
	}
}

func TestLinkLookupsRejectInvalidInput(t *testing.T) {
	if _, err := zcall.LinkByName(""); !errors.Is(err, zcall.Errno(zcall.EINVAL)) {
		t.Fatalf("LinkByName(empty) error = %v, want EINVAL", err)
	}
	if _, err := zcall.LinkByIndex(0); !errors.Is(err, zcall.Errno(zcall.EINVAL)) {
		t.Fatalf("LinkByIndex(0) error = %v, want EINVAL", err)
	}
}

// noescape hides a pointer from escape analysis.
// Used to safely convert uintptr (from mmap) to unsafe.Pointer.
//
//go:linkname noescape runtime.noescape
//go:noescape
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer

//go:linkname openRouteNetlink code.hybscloud.com/zcall.openRouteNetlink
func openRouteNetlink() (uintptr, uint32, error)

//go:linkname sendNetlinkGetLinkDump code.hybscloud.com/zcall.sendNetlinkGetLinkDump
func sendNetlinkGetLinkDump(fd uintptr, seq uint32) error

//go:linkname recvLinks code.hybscloud.com/zcall.recvLinks
func recvLinks(fd uintptr, pid, seq uint32) ([]zcall.Link, error)

//go:linkname parseLink code.hybscloud.com/zcall.parseLink
func parseLink(payload []byte) (zcall.Link, bool, error)

//go:linkname nlmsgAlignOf code.hybscloud.com/zcall.nlmsgAlignOf
func nlmsgAlignOf(length int) int

//go:linkname rtaAlignOf code.hybscloud.com/zcall.rtaAlignOf
func rtaAlignOf(length int) int

type ifreq struct {
	name [zcall.IFNAMSIZ]byte
	data [24]byte
}

type sockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

type nlMsghdr struct {
	Len   uint32
	Type  uint16
	Flags uint16
	Seq   uint32
	Pid   uint32
}

type ifInfomsg struct {
	Family uint8
	_      uint8
	Type   uint16
	Index  int32
	Flags  uint32
	Change uint32
}

type rtAttr struct {
	Len  uint16
	Type uint16
}

func TestOpenRouteNetlink(t *testing.T) {
	fd, pid, err := openRouteNetlink()
	if err != nil {
		t.Fatalf("openRouteNetlink() failed: %v", err)
	}
	defer zcall.Close(fd)

	if fd == 0 || fd == ^uintptr(0) {
		t.Fatalf("openRouteNetlink() fd = %d, want valid descriptor", fd)
	}
	if pid == 0 {
		t.Fatal("openRouteNetlink() pid = 0, want assigned port ID")
	}
}

func TestSendNetlinkGetLinkDumpRejectsInvalidFD(t *testing.T) {
	if err := sendNetlinkGetLinkDump(^uintptr(0), 3); !errors.Is(err, zcall.Errno(zcall.EBADF)) {
		t.Fatalf("sendNetlinkGetLinkDump() error = %v, want EBADF", err)
	}
}

func TestLinkLookupsNotFound(t *testing.T) {
	if _, err := zcall.LinkByName("missing0"); !errors.Is(err, zcall.Errno(zcall.ENODEV)) {
		t.Fatalf("LinkByName(missing) error = %v, want ENODEV", err)
	}
	if _, err := zcall.LinkByIndex(1 << 30); !errors.Is(err, zcall.Errno(zcall.ENODEV)) {
		t.Fatalf("LinkByIndex(missing) error = %v, want ENODEV", err)
	}
}

func TestRecvLinksRejectsInvalidFD(t *testing.T) {
	_, err := recvLinks(^uintptr(0), 1, 1)
	if !errors.Is(err, zcall.Errno(zcall.EBADF)) {
		t.Fatalf("recvLinks() error = %v, want EBADF", err)
	}
}

func TestIoctlSIOCGIFINDEX(t *testing.T) {
	fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_DGRAM|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket() errno = %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	var req ifreq
	copy(req.name[:], "lo")

	if errno := zcall.Ioctl(fd, zcall.SIOCGIFINDEX, unsafe.Pointer(&req)); errno != 0 {
		t.Fatalf("Ioctl(SIOCGIFINDEX) errno = %v", zcall.Errno(errno))
	}

	if index := *(*int32)(unsafe.Pointer(&req.data[0])); index <= 0 {
		t.Fatalf("Ioctl(SIOCGIFINDEX) index = %d, want > 0", index)
	}
}

func TestRecvLinksParsesMultipartReply(t *testing.T) {
	reader, writer := mustSocketpairDatagram(t)
	defer zcall.Close(reader)
	defer zcall.Close(writer)

	const (
		pid = 77
		seq = 9
	)

	payload := appendIfInfoPayload(ifInfomsg{Index: 5, Flags: zcall.IFF_UP | zcall.IFF_MULTICAST},
		appendRtAttr(zcall.IFLA_IFNAME, append([]byte("eth-test"), 0)),
		appendRtAttr(zcall.IFLA_MTU, uint32Bytes(1500)),
	)
	msg := appendNetlinkMessage(nil, nlMsghdr{Type: zcall.RTM_NEWLINK, Seq: seq, Pid: pid}, payload)
	msg = appendNetlinkMessage(msg, nlMsghdr{Type: zcall.NLMSG_DONE, Seq: seq, Pid: pid}, nil)
	writeAll(t, writer, msg)

	links, err := recvLinks(reader, pid, seq)
	if err != nil {
		t.Fatalf("recvLinks() error = %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("recvLinks() len = %d, want 1", len(links))
	}
	if links[0].Name != "eth-test" || links[0].Index != 5 || links[0].MTU != 1500 {
		t.Fatalf("recvLinks() link = %+v, want parsed attributes", links[0])
	}
}

func TestRecvLinksRejectsMalformedMessages(t *testing.T) {
	tests := []struct {
		name    string
		message []byte
		pid     uint32
		seq     uint32
		wantErr zcall.Errno
	}{
		{
			name:    "short read",
			message: []byte{0x01},
			pid:     1,
			seq:     1,
			wantErr: zcall.EINVAL,
		},
		{
			name:    "sequence mismatch",
			message: appendNetlinkMessage(nil, nlMsghdr{Type: zcall.NLMSG_DONE, Seq: 2, Pid: 1}, nil),
			pid:     1,
			seq:     1,
			wantErr: zcall.EINVAL,
		},
		{
			name:    "invalid message length",
			message: rawNetlinkMessage(nlMsghdr{Len: uint32(unsafe.Sizeof(nlMsghdr{})) - 1, Type: zcall.NLMSG_DONE, Seq: 1, Pid: 1}, nil),
			pid:     1,
			seq:     1,
			wantErr: zcall.EINVAL,
		},
		{
			name:    "short nlmsg error payload",
			message: appendNetlinkMessage(nil, nlMsghdr{Type: zcall.NLMSG_ERROR, Seq: 1, Pid: 1}, []byte{0x01, 0x02, 0x03}),
			pid:     1,
			seq:     1,
			wantErr: zcall.EINVAL,
		},
		{
			name:    "positive nlmsg error payload",
			message: appendNetlinkMessage(nil, nlMsghdr{Type: zcall.NLMSG_ERROR, Seq: 1, Pid: 1}, int32Bytes(0)),
			pid:     1,
			seq:     1,
			wantErr: zcall.EINVAL,
		},
		{
			name:    "negative nlmsg error payload",
			message: appendNetlinkMessage(nil, nlMsghdr{Type: zcall.NLMSG_ERROR, Seq: 1, Pid: 1}, int32Bytes(-int32(zcall.ENODEV))),
			pid:     1,
			seq:     1,
			wantErr: zcall.ENODEV,
		},
		{
			name: "newlink parse error",
			message: appendNetlinkMessage(nil, nlMsghdr{Type: zcall.RTM_NEWLINK, Seq: 1, Pid: 1},
				appendIfInfoPayload(ifInfomsg{Index: 1}, malformedRtAttr(2))),
			pid:     1,
			seq:     1,
			wantErr: zcall.EINVAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, writer := mustSocketpairDatagram(t)
			defer zcall.Close(reader)
			defer zcall.Close(writer)

			writeAll(t, writer, tt.message)

			_, err := recvLinks(reader, tt.pid, tt.seq)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("recvLinks() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseLinkParsesCoreAttributes(t *testing.T) {
	payload := appendIfInfoPayload(ifInfomsg{Index: 7, Flags: zcall.IFF_UP | zcall.IFF_MULTICAST},
		appendRtAttr(zcall.IFLA_IFNAME, append([]byte("eth0"), 0)),
		appendRtAttr(zcall.IFLA_MTU, uint32Bytes(9000)),
		appendRtAttr(zcall.IFLA_ADDRESS, []byte{0x02, 0x42, 0xac, 0x11, 0x00, 0x02}),
	)

	link, ok, err := parseLink(payload)
	if err != nil {
		t.Fatalf("parseLink() error = %v", err)
	}
	if !ok {
		t.Fatal("parseLink() reported !ok")
	}
	if link.Index != 7 {
		t.Fatalf("parseLink().Index = %d, want 7", link.Index)
	}
	if link.Flags != zcall.IFF_UP|zcall.IFF_MULTICAST {
		t.Fatalf("parseLink().Flags = %#x, want %#x", link.Flags, zcall.IFF_UP|zcall.IFF_MULTICAST)
	}
	if link.Name != "eth0" {
		t.Fatalf("parseLink().Name = %q, want %q", link.Name, "eth0")
	}
	if link.MTU != 9000 {
		t.Fatalf("parseLink().MTU = %d, want 9000", link.MTU)
	}
	if got := len(link.HardwareAddr); got != 6 {
		t.Fatalf("parseLink().HardwareAddr len = %d, want 6", got)
	}
	if link.HardwareAddr[0] != 0x02 || link.HardwareAddr[5] != 0x02 {
		t.Fatalf("parseLink().HardwareAddr = %v, want preserved bytes", link.HardwareAddr)
	}
}

func TestParseLinkRejectsMalformedAttribute(t *testing.T) {
	payload := appendIfInfoPayload(ifInfomsg{Index: 1}, malformedRtAttr(2))

	_, _, err := parseLink(payload)
	if !errors.Is(err, zcall.EINVAL) {
		t.Fatalf("parseLink() error = %v, want EINVAL", err)
	}
}

func TestParseLinkSkipsUnnamedOrInvalidLinks(t *testing.T) {
	t.Run("short payload", func(t *testing.T) {
		_, ok, err := parseLink(make([]byte, unsafe.Sizeof(ifInfomsg{})-1))
		if err != nil {
			t.Fatalf("parseLink(short) error = %v", err)
		}
		if ok {
			t.Fatal("parseLink(short) ok = true, want false")
		}
	})

	t.Run("invalid index", func(t *testing.T) {
		payload := appendIfInfoPayload(ifInfomsg{Index: 0}, appendRtAttr(zcall.IFLA_IFNAME, append([]byte("lo"), 0)))
		_, ok, err := parseLink(payload)
		if err != nil {
			t.Fatalf("parseLink(invalid index) error = %v", err)
		}
		if ok {
			t.Fatal("parseLink(invalid index) ok = true, want false")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		payload := appendIfInfoPayload(ifInfomsg{Index: 2}, appendRtAttr(zcall.IFLA_MTU, uint32Bytes(1500)))
		_, ok, err := parseLink(payload)
		if err != nil {
			t.Fatalf("parseLink(missing name) error = %v", err)
		}
		if ok {
			t.Fatal("parseLink(missing name) ok = true, want false")
		}
	})
}

func TestAlignmentHelpers(t *testing.T) {
	if got := nlmsgAlignOf(17); got != 20 {
		t.Fatalf("nlmsgAlignOf(17) = %d, want 20", got)
	}
	if got := rtaAlignOf(5); got != 8 {
		t.Fatalf("rtaAlignOf(5) = %d, want 8", got)
	}
}

func appendIfInfoPayload(info ifInfomsg, attrs ...[]byte) []byte {
	payload := make([]byte, unsafe.Sizeof(info))
	*(*ifInfomsg)(unsafe.Pointer(&payload[0])) = info
	for _, attr := range attrs {
		payload = append(payload, attr...)
	}
	return payload
}

func appendRtAttr(typ uint16, value []byte) []byte {
	const hdrSize = int(unsafe.Sizeof(rtAttr{}))
	attrLen := hdrSize + len(value)
	alignedLen := rtaAlignOf(attrLen)
	buf := make([]byte, alignedLen)
	attr := (*rtAttr)(unsafe.Pointer(&buf[0]))
	attr.Len = uint16(attrLen)
	attr.Type = typ
	copy(buf[hdrSize:hdrSize+len(value)], value)
	return buf
}

func malformedRtAttr(length uint16) []byte {
	buf := make([]byte, unsafe.Sizeof(rtAttr{}))
	attr := (*rtAttr)(unsafe.Pointer(&buf[0]))
	attr.Len = length
	attr.Type = zcall.IFLA_IFNAME
	return buf
}

func uint32Bytes(v uint32) []byte {
	buf := make([]byte, 4)
	*(*uint32)(unsafe.Pointer(&buf[0])) = v
	return buf
}

func int32Bytes(v int32) []byte {
	buf := make([]byte, 4)
	*(*int32)(unsafe.Pointer(&buf[0])) = v
	return buf
}

func mustSocketpairDatagram(t *testing.T) (uintptr, uintptr) {
	t.Helper()
	var fds [2]int32
	if errno := zcall.Socketpair(zcall.AF_UNIX, zcall.SOCK_DGRAM|zcall.SOCK_CLOEXEC, 0, &fds); errno != 0 {
		t.Fatalf("Socketpair() errno = %v", zcall.Errno(errno))
	}
	return uintptr(fds[0]), uintptr(fds[1])
}

func writeAll(t *testing.T, fd uintptr, msg []byte) {
	t.Helper()
	n, errno := zcall.Write(fd, msg)
	if errno != 0 {
		t.Fatalf("Write() errno = %v", zcall.Errno(errno))
	}
	if n != uintptr(len(msg)) {
		t.Fatalf("Write() = %d, want %d", n, len(msg))
	}
}

func appendNetlinkMessage(dst []byte, hdr nlMsghdr, payload []byte) []byte {
	const hdrSize = int(unsafe.Sizeof(nlMsghdr{}))
	hdr.Len = uint32(hdrSize + len(payload))
	alignedLen := nlmsgAlignOf(int(hdr.Len))
	start := len(dst)
	dst = append(dst, make([]byte, alignedLen)...)
	msg := dst[start : start+alignedLen]
	*(*nlMsghdr)(unsafe.Pointer(&msg[0])) = hdr
	copy(msg[hdrSize:hdrSize+len(payload)], payload)
	return dst
}

func rawNetlinkMessage(hdr nlMsghdr, payload []byte) []byte {
	buf := make([]byte, int(unsafe.Sizeof(nlMsghdr{}))+len(payload))
	*(*nlMsghdr)(unsafe.Pointer(&buf[0])) = hdr
	copy(buf[int(unsafe.Sizeof(nlMsghdr{})):], payload)
	return buf
}

func TestEventfd2(t *testing.T) {
	fd, errno := zcall.Eventfd2(0, zcall.EFD_NONBLOCK|zcall.EFD_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Eventfd2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	if fd == 0 || fd == ^uintptr(0) {
		t.Fatalf("Eventfd2 returned invalid fd: %d", fd)
	}

	// Write to eventfd
	val := uint64(1)
	n, errno := zcall.Write(fd, (*[8]byte)(unsafe.Pointer(&val))[:])
	if errno != 0 {
		t.Fatalf("Write to eventfd failed: %v", zcall.Errno(errno))
	}
	if n != 8 {
		t.Fatalf("Write returned %d, expected 8", n)
	}

	// Read from eventfd
	var readVal uint64
	n, errno = zcall.Read(fd, (*[8]byte)(unsafe.Pointer(&readVal))[:])
	if errno != 0 {
		t.Fatalf("Read from eventfd failed: %v", zcall.Errno(errno))
	}
	if n != 8 {
		t.Fatalf("Read returned %d, expected 8", n)
	}
	if readVal != 1 {
		t.Fatalf("Read value %d, expected 1", readVal)
	}
}

func TestTimerfdCreate(t *testing.T) {
	fd, errno := zcall.TimerfdCreate(zcall.CLOCK_MONOTONIC, zcall.TFD_NONBLOCK|zcall.TFD_CLOEXEC)
	if errno != 0 {
		t.Fatalf("TimerfdCreate failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	if fd == 0 || fd == ^uintptr(0) {
		t.Fatalf("TimerfdCreate returned invalid fd: %d", fd)
	}
}

func TestSocket(t *testing.T) {
	fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	if fd == 0 || fd == ^uintptr(0) {
		t.Fatalf("Socket returned invalid fd: %d", fd)
	}
}

func TestSocketpair(t *testing.T) {
	var fds [2]int32
	errno := zcall.Socketpair(zcall.AF_UNIX, zcall.SOCK_STREAM|zcall.SOCK_CLOEXEC, 0, &fds)
	if errno != 0 {
		t.Fatalf("Socketpair failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	// Write to one end
	msg := []byte("hello")
	n, errno := zcall.Write(uintptr(fds[0]), msg)
	if errno != 0 {
		t.Fatalf("Write failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(msg)) {
		t.Fatalf("Write returned %d, expected %d", n, len(msg))
	}

	// Read from other end
	buf := make([]byte, 16)
	n, errno = zcall.Read(uintptr(fds[1]), buf)
	if errno != 0 {
		t.Fatalf("Read failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(msg)) {
		t.Fatalf("Read returned %d, expected %d", n, len(msg))
	}
	if string(buf[:n]) != "hello" {
		t.Fatalf("Read data mismatch: got %q, expected %q", buf[:n], msg)
	}
}

func TestPipe2(t *testing.T) {
	var fds [2]int32
	errno := zcall.Pipe2(&fds, zcall.O_NONBLOCK|zcall.O_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	// Write to write end (fds[1])
	msg := []byte("pipe test")
	n, errno := zcall.Write(uintptr(fds[1]), msg)
	if errno != 0 {
		t.Fatalf("Write to pipe failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(msg)) {
		t.Fatalf("Write returned %d, expected %d", n, len(msg))
	}

	// Read from read end (fds[0])
	buf := make([]byte, 32)
	n, errno = zcall.Read(uintptr(fds[0]), buf)
	if errno != 0 {
		t.Fatalf("Read from pipe failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(msg)) {
		t.Fatalf("Read returned %d, expected %d", n, len(msg))
	}
	if string(buf[:n]) != "pipe test" {
		t.Fatalf("Read data mismatch: got %q, expected %q", buf[:n], msg)
	}
}

func TestReadvWritev(t *testing.T) {
	var fds [2]int32
	errno := zcall.Socketpair(zcall.AF_UNIX, zcall.SOCK_STREAM|zcall.SOCK_CLOEXEC, 0, &fds)
	if errno != 0 {
		t.Fatalf("Socketpair failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	// Prepare iovec for writev
	buf1 := []byte("hello")
	buf2 := []byte("world")
	iovWrite := []zcall.Iovec{
		{Base: &buf1[0], Len: uint64(len(buf1))},
		{Base: &buf2[0], Len: uint64(len(buf2))},
	}

	n, errno := zcall.Writev(uintptr(fds[0]), unsafe.Pointer(&iovWrite[0]), uintptr(len(iovWrite)))
	if errno != 0 {
		t.Fatalf("Writev failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(buf1)+len(buf2)) {
		t.Fatalf("Writev returned %d, expected %d", n, len(buf1)+len(buf2))
	}

	// Prepare iovec for readv
	rbuf1 := make([]byte, 5)
	rbuf2 := make([]byte, 5)
	iovRead := []zcall.Iovec{
		{Base: &rbuf1[0], Len: uint64(len(rbuf1))},
		{Base: &rbuf2[0], Len: uint64(len(rbuf2))},
	}

	n, errno = zcall.Readv(uintptr(fds[1]), unsafe.Pointer(&iovRead[0]), uintptr(len(iovRead)))
	if errno != 0 {
		t.Fatalf("Readv failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(buf1)+len(buf2)) {
		t.Fatalf("Readv returned %d, expected %d", n, len(buf1)+len(buf2))
	}
	if string(rbuf1) != "hello" || string(rbuf2) != "world" {
		t.Fatalf("Readv data mismatch: got %q %q, expected %q %q", rbuf1, rbuf2, buf1, buf2)
	}
}

func TestSplice(t *testing.T) {
	// Create two pipes
	var pipe1 [2]int32
	errno := zcall.Pipe2(&pipe1, zcall.O_NONBLOCK|zcall.O_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(pipe1[0]))
	defer zcall.Close(uintptr(pipe1[1]))

	var pipe2 [2]int32
	errno = zcall.Pipe2(&pipe2, zcall.O_NONBLOCK|zcall.O_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(pipe2[0]))
	defer zcall.Close(uintptr(pipe2[1]))

	// Write data to first pipe
	msg := []byte("splice test data")
	n, errno := zcall.Write(uintptr(pipe1[1]), msg)
	if errno != 0 {
		t.Fatalf("Write failed: %v", zcall.Errno(errno))
	}

	// Splice from pipe1 read end to pipe2 write end
	spliced, errno := zcall.Splice(uintptr(pipe1[0]), nil, uintptr(pipe2[1]), nil, uintptr(n), zcall.SPLICE_F_NONBLOCK)
	if errno != 0 {
		t.Fatalf("Splice failed: %v", zcall.Errno(errno))
	}
	if spliced != n {
		t.Fatalf("Splice returned %d, expected %d", spliced, n)
	}

	// Read from pipe2
	buf := make([]byte, 32)
	n, errno = zcall.Read(uintptr(pipe2[0]), buf)
	if errno != 0 {
		t.Fatalf("Read failed: %v", zcall.Errno(errno))
	}
	if string(buf[:n]) != "splice test data" {
		t.Fatalf("Read data mismatch: got %q, expected %q", buf[:n], msg)
	}
}

func TestTee(t *testing.T) {
	// Create two pipes
	var pipe1 [2]int32
	errno := zcall.Pipe2(&pipe1, zcall.O_NONBLOCK|zcall.O_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(pipe1[0]))
	defer zcall.Close(uintptr(pipe1[1]))

	var pipe2 [2]int32
	errno = zcall.Pipe2(&pipe2, zcall.O_NONBLOCK|zcall.O_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(pipe2[0]))
	defer zcall.Close(uintptr(pipe2[1]))

	// Write data to first pipe
	msg := []byte("tee test")
	n, errno := zcall.Write(uintptr(pipe1[1]), msg)
	if errno != 0 {
		t.Fatalf("Write failed: %v", zcall.Errno(errno))
	}

	// Tee from pipe1 to pipe2 (duplicates without consuming)
	teed, errno := zcall.Tee(uintptr(pipe1[0]), uintptr(pipe2[1]), uintptr(n), zcall.SPLICE_F_NONBLOCK)
	if errno != 0 {
		t.Fatalf("Tee failed: %v", zcall.Errno(errno))
	}
	if teed != n {
		t.Fatalf("Tee returned %d, expected %d", teed, n)
	}

	// Read from pipe2 (should have copy)
	buf := make([]byte, 32)
	n2, errno := zcall.Read(uintptr(pipe2[0]), buf)
	if errno != 0 {
		t.Fatalf("Read from pipe2 failed: %v", zcall.Errno(errno))
	}
	if string(buf[:n2]) != "tee test" {
		t.Fatalf("Read data mismatch: got %q, expected %q", buf[:n2], msg)
	}

	// Read from pipe1 (data should still be there)
	n1, errno := zcall.Read(uintptr(pipe1[0]), buf)
	if errno != 0 {
		t.Fatalf("Read from pipe1 failed: %v", zcall.Errno(errno))
	}
	if string(buf[:n1]) != "tee test" {
		t.Fatalf("Read data mismatch: got %q, expected %q", buf[:n1], msg)
	}
}

func TestSetsockoptGetsockopt(t *testing.T) {
	fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	// Set SO_REUSEADDR
	val := int32(1)
	errno = zcall.Setsockopt(fd, zcall.SOL_SOCKET, zcall.SO_REUSEADDR, unsafe.Pointer(&val), 4)
	if errno != 0 {
		t.Fatalf("Setsockopt SO_REUSEADDR failed: %v", zcall.Errno(errno))
	}

	// Get SO_REUSEADDR
	var getVal int32
	optlen := uint32(4)
	errno = zcall.Getsockopt(fd, zcall.SOL_SOCKET, zcall.SO_REUSEADDR, unsafe.Pointer(&getVal), unsafe.Pointer(&optlen))
	if errno != 0 {
		t.Fatalf("Getsockopt SO_REUSEADDR failed: %v", zcall.Errno(errno))
	}
	if getVal != 1 {
		t.Fatalf("Getsockopt returned %d, expected 1", getVal)
	}
}

func TestTimerfdSettime(t *testing.T) {
	fd, errno := zcall.TimerfdCreate(zcall.CLOCK_MONOTONIC, zcall.TFD_NONBLOCK|zcall.TFD_CLOEXEC)
	if errno != 0 {
		t.Fatalf("TimerfdCreate failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	// Set timer to expire in 100ms
	newValue := zcall.Itimerspec{
		Value: zcall.Timespec{Sec: 0, Nsec: 100000000},
	}
	var oldValue zcall.Itimerspec

	errno = zcall.TimerfdSettime(fd, 0, unsafe.Pointer(&newValue), unsafe.Pointer(&oldValue))
	if errno != 0 {
		t.Fatalf("TimerfdSettime failed: %v", zcall.Errno(errno))
	}

	// Get current timer value
	var currValue zcall.Itimerspec
	errno = zcall.TimerfdGettime(fd, unsafe.Pointer(&currValue))
	if errno != 0 {
		t.Fatalf("TimerfdGettime failed: %v", zcall.Errno(errno))
	}

	// Timer should be armed (value should be non-zero)
	if currValue.Value.Sec == 0 && currValue.Value.Nsec == 0 {
		t.Fatal("Timer should be armed but value is zero")
	}
}

func TestMmapMunmap(t *testing.T) {
	// Map anonymous memory
	size := uintptr(4096)
	ptr, errno := zcall.Mmap(nil, size, zcall.PROT_READ|zcall.PROT_WRITE, zcall.MAP_PRIVATE|zcall.MAP_ANONYMOUS, ^uintptr(0), 0)
	if errno != 0 {
		t.Fatalf("Mmap failed: %v", zcall.Errno(errno))
	}

	// Verify mmap returned a valid address (not NULL or MAP_FAILED)
	addr := uintptr(ptr)
	if addr == 0 || addr == ^uintptr(0) {
		t.Fatalf("Mmap returned invalid address: %x", addr)
	}

	// Verify the address is page-aligned (4KB alignment)
	if addr&0xFFF != 0 {
		t.Fatalf("Mmap returned non-page-aligned address: %x", addr)
	}

	// Unmap the memory
	errno = zcall.Munmap(ptr, size)
	if errno != 0 {
		t.Fatalf("Munmap failed: %v", zcall.Errno(errno))
	}
}

func TestGetsocknameGetpeername(t *testing.T) {
	// Create a listening socket
	listenFd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(listenFd)

	// Bind to localhost with any port
	addr := [16]byte{
		2, 0, // AF_INET
		0, 0, // port 0 (any)
		127, 0, 0, 1, // 127.0.0.1
	}
	errno = zcall.Bind(listenFd, unsafe.Pointer(&addr), 16)
	if errno != 0 {
		t.Fatalf("Bind failed: %v", zcall.Errno(errno))
	}

	// Get socket name to find assigned port
	var boundAddr [16]byte
	addrLen := uint32(16)
	errno = zcall.Getsockname(listenFd, unsafe.Pointer(&boundAddr), unsafe.Pointer(&addrLen))
	if errno != 0 {
		t.Fatalf("Getsockname failed: %v", zcall.Errno(errno))
	}

	// Verify address family
	if boundAddr[0] != 2 { // AF_INET
		t.Fatalf("Getsockname returned wrong family: %d", boundAddr[0])
	}
}

func TestShutdown(t *testing.T) {
	var fds [2]int32
	errno := zcall.Socketpair(zcall.AF_UNIX, zcall.SOCK_STREAM|zcall.SOCK_CLOEXEC, 0, &fds)
	if errno != 0 {
		t.Fatalf("Socketpair failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	// Shutdown write end
	errno = zcall.Shutdown(uintptr(fds[0]), zcall.SHUT_WR)
	if errno != 0 {
		t.Fatalf("Shutdown failed: %v", zcall.Errno(errno))
	}

	// Read should return 0 (EOF)
	buf := make([]byte, 16)
	n, errno := zcall.Read(uintptr(fds[1]), buf)
	if errno != 0 {
		t.Fatalf("Read after shutdown failed: %v", zcall.Errno(errno))
	}
	if n != 0 {
		t.Fatalf("Read after shutdown returned %d, expected 0", n)
	}
}

func TestSendtoRecvfrom(t *testing.T) {
	// Create UDP socket pair using socketpair isn't possible for UDP,
	// so we'll use a single socket bound to localhost
	fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_DGRAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	// Bind to localhost
	addr := [16]byte{
		2, 0, // AF_INET
		0, 0, // port 0 (any)
		127, 0, 0, 1, // 127.0.0.1
	}
	errno = zcall.Bind(fd, unsafe.Pointer(&addr), 16)
	if errno != 0 {
		t.Fatalf("Bind failed: %v", zcall.Errno(errno))
	}

	// Get bound address
	var boundAddr [16]byte
	addrLen := uint32(16)
	errno = zcall.Getsockname(fd, unsafe.Pointer(&boundAddr), unsafe.Pointer(&addrLen))
	if errno != 0 {
		t.Fatalf("Getsockname failed: %v", zcall.Errno(errno))
	}

	// Send to self
	msg := []byte("udp test")
	n, errno := zcall.Sendto(fd, msg, 0, unsafe.Pointer(&boundAddr), 16)
	if errno != 0 {
		t.Fatalf("Sendto failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(msg)) {
		t.Fatalf("Sendto returned %d, expected %d", n, len(msg))
	}

	// Receive
	buf := make([]byte, 32)
	var fromAddr [16]byte
	fromLen := uint32(16)
	n, errno = zcall.Recvfrom(fd, buf, 0, unsafe.Pointer(&fromAddr), unsafe.Pointer(&fromLen))
	if errno != 0 {
		t.Fatalf("Recvfrom failed: %v", zcall.Errno(errno))
	}
	if string(buf[:n]) != "udp test" {
		t.Fatalf("Recvfrom data mismatch: got %q, expected %q", buf[:n], msg)
	}
}

func TestVmsplice(t *testing.T) {
	// Create a pipe
	var fds [2]int32
	errno := zcall.Pipe2(&fds, zcall.O_NONBLOCK|zcall.O_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	// Prepare data
	data := []byte("vmsplice test")
	iov := zcall.Iovec{
		Base: &data[0],
		Len:  uint64(len(data)),
	}

	// Vmsplice into pipe write end
	n, errno := zcall.Vmsplice(uintptr(fds[1]), unsafe.Pointer(&iov), 1, zcall.SPLICE_F_GIFT)
	if errno != 0 {
		t.Fatalf("Vmsplice failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(data)) {
		t.Fatalf("Vmsplice returned %d, expected %d", n, len(data))
	}

	// Read from pipe
	buf := make([]byte, 32)
	n, errno = zcall.Read(uintptr(fds[0]), buf)
	if errno != 0 {
		t.Fatalf("Read failed: %v", zcall.Errno(errno))
	}
	if string(buf[:n]) != "vmsplice test" {
		t.Fatalf("Read data mismatch: got %q, expected %q", buf[:n], data)
	}
}

func TestSendmsgRecvmsg(t *testing.T) {
	var fds [2]int32
	errno := zcall.Socketpair(zcall.AF_UNIX, zcall.SOCK_STREAM|zcall.SOCK_CLOEXEC, 0, &fds)
	if errno != 0 {
		t.Fatalf("Socketpair failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	// Prepare message
	data := []byte("sendmsg test")
	iov := zcall.Iovec{
		Base: &data[0],
		Len:  uint64(len(data)),
	}
	msg := zcall.Msghdr{
		Iov:    &iov,
		Iovlen: 1,
	}

	// Send
	n, errno := zcall.Sendmsg(uintptr(fds[0]), unsafe.Pointer(&msg), 0)
	if errno != 0 {
		t.Fatalf("Sendmsg failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(data)) {
		t.Fatalf("Sendmsg returned %d, expected %d", n, len(data))
	}

	// Receive
	buf := make([]byte, 32)
	recvIov := zcall.Iovec{
		Base: &buf[0],
		Len:  uint64(len(buf)),
	}
	recvMsg := zcall.Msghdr{
		Iov:    &recvIov,
		Iovlen: 1,
	}

	n, errno = zcall.Recvmsg(uintptr(fds[1]), unsafe.Pointer(&recvMsg), 0)
	if errno != 0 {
		t.Fatalf("Recvmsg failed: %v", zcall.Errno(errno))
	}
	if string(buf[:n]) != "sendmsg test" {
		t.Fatalf("Recvmsg data mismatch: got %q, expected %q", buf[:n], data)
	}
}

func TestBindListenAccept(t *testing.T) {
	// Create listening socket
	listenFd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(listenFd)

	// Bind to localhost
	addr := [16]byte{
		2, 0, // AF_INET
		0, 0, // port 0 (any)
		127, 0, 0, 1, // 127.0.0.1
	}
	errno = zcall.Bind(listenFd, unsafe.Pointer(&addr), 16)
	if errno != 0 {
		t.Fatalf("Bind failed: %v", zcall.Errno(errno))
	}

	// Listen
	errno = zcall.Listen(listenFd, 128)
	if errno != 0 {
		t.Fatalf("Listen failed: %v", zcall.Errno(errno))
	}

	// Get bound address
	var boundAddr [16]byte
	addrLen := uint32(16)
	errno = zcall.Getsockname(listenFd, unsafe.Pointer(&boundAddr), unsafe.Pointer(&addrLen))
	if errno != 0 {
		t.Fatalf("Getsockname failed: %v", zcall.Errno(errno))
	}

	// Create client socket
	clientFd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(clientFd)

	// Connect (will return EINPROGRESS for non-blocking)
	errno = zcall.Connect(clientFd, unsafe.Pointer(&boundAddr), 16)
	if errno != 0 && zcall.Errno(errno) != zcall.EINPROGRESS {
		t.Fatalf("Connect failed: %v", zcall.Errno(errno))
	}

	// Accept (may return EAGAIN if connection not ready)
	var clientAddr [16]byte
	clientAddrLen := uint32(16)
	acceptFd, errno := zcall.Accept4(listenFd, unsafe.Pointer(&clientAddr), unsafe.Pointer(&clientAddrLen), zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC)
	if errno != 0 && zcall.Errno(errno) != zcall.EAGAIN {
		t.Fatalf("Accept4 failed: %v", zcall.Errno(errno))
	}
	if errno == 0 {
		zcall.Close(acceptFd)
	}
}

func TestAcceptAndGetpeername(t *testing.T) {
	// Create non-blocking listening socket
	// Non-blocking is required because zcall bypasses Go's netpoller
	listenFd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(listenFd)

	// Bind to localhost with ephemeral port
	addr := [16]byte{2, 0, 0, 0, 127, 0, 0, 1}
	errno = zcall.Bind(listenFd, unsafe.Pointer(&addr), 16)
	if errno != 0 {
		t.Fatalf("Bind failed: %v", zcall.Errno(errno))
	}

	errno = zcall.Listen(listenFd, 128)
	if errno != 0 {
		t.Fatalf("Listen failed: %v", zcall.Errno(errno))
	}

	// Get bound address (includes assigned port)
	var boundAddr [16]byte
	addrLen := uint32(16)
	errno = zcall.Getsockname(listenFd, unsafe.Pointer(&boundAddr), unsafe.Pointer(&addrLen))
	if errno != 0 {
		t.Fatalf("Getsockname failed: %v", zcall.Errno(errno))
	}

	// Create non-blocking client socket
	clientFd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(clientFd)

	// Non-blocking connect returns EINPROGRESS
	errno = zcall.Connect(clientFd, unsafe.Pointer(&boundAddr), 16)
	if errno != 0 && errno != uintptr(zcall.EINPROGRESS) {
		t.Fatalf("Connect failed: %v", zcall.Errno(errno))
	}

	// Spin-wait for Accept (non-blocking returns EAGAIN until connection ready)
	var acceptFd uintptr
	var peerAddr [16]byte
	peerAddrLen := uint32(16)
	for range 1000 {
		acceptFd, errno = zcall.Accept(listenFd, unsafe.Pointer(&peerAddr), unsafe.Pointer(&peerAddrLen))
		if errno == 0 {
			break
		}
		if errno != uintptr(zcall.EAGAIN) && errno != uintptr(zcall.EWOULDBLOCK) {
			t.Fatalf("Accept failed: %v", zcall.Errno(errno))
		}
		// Brief yield to allow connection to complete
		for range 1000 {
			// spin
		}
	}
	if errno != 0 {
		t.Fatalf("Accept timed out after spinning")
	}
	defer zcall.Close(acceptFd)

	// Test Getpeername on accepted socket
	var remotePeer [16]byte
	remotePeerLen := uint32(16)
	errno = zcall.Getpeername(acceptFd, unsafe.Pointer(&remotePeer), unsafe.Pointer(&remotePeerLen))
	if errno != 0 {
		t.Fatalf("Getpeername failed: %v", zcall.Errno(errno))
	}

	// Verify it's an IPv4 address (family = AF_INET = 2)
	if remotePeer[0] != 2 {
		t.Fatalf("Getpeername returned wrong family: %d", remotePeer[0])
	}
}

func TestPreadvPwritev(t *testing.T) {
	// Preadv/Pwritev require seekable file descriptors (not sockets/pipes)
	// Test by verifying the syscall interface works - sockets return ESPIPE
	var fds [2]int32
	errno := zcall.Socketpair(zcall.AF_UNIX, zcall.SOCK_STREAM|zcall.SOCK_CLOEXEC, 0, &fds)
	if errno != 0 {
		t.Fatalf("Socketpair failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	buf := []byte("preadv pwritev test")
	iov := []zcall.Iovec{
		{Base: &buf[0], Len: uint64(len(buf))},
	}

	// Pwritev on socket should return ESPIPE (illegal seek)
	_, errno = zcall.Pwritev(uintptr(fds[0]), unsafe.Pointer(&iov[0]), 1, 0)
	if errno == 0 {
		t.Fatal("Pwritev on socket should fail with ESPIPE")
	}
	if zcall.Errno(errno) != zcall.ESPIPE {
		t.Fatalf("Pwritev errno = %v, want ESPIPE", zcall.Errno(errno))
	}

	// Preadv on socket should also return ESPIPE
	rbuf := make([]byte, 32)
	rIov := []zcall.Iovec{
		{Base: &rbuf[0], Len: uint64(len(rbuf))},
	}
	_, errno = zcall.Preadv(uintptr(fds[1]), unsafe.Pointer(&rIov[0]), 1, 0)
	if errno == 0 {
		t.Fatal("Preadv on socket should fail with ESPIPE")
	}
	if zcall.Errno(errno) != zcall.ESPIPE {
		t.Fatalf("Preadv errno = %v, want ESPIPE", zcall.Errno(errno))
	}
}

func TestPreadvPwritevHighOffset(t *testing.T) {
	f, err := os.CreateTemp("", "zcall-preadv-pwritev-*")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	fd := uintptr(f.Fd())
	offset := int64(1<<33 + 123)
	want := []byte("preadv-pwritev high offset")
	iov := []zcall.Iovec{
		{Base: &want[0], Len: uint64(len(want))},
	}

	n, errno := zcall.Pwritev(fd, unsafe.Pointer(&iov[0]), 1, offset)
	if errno != 0 {
		t.Fatalf("Pwritev high offset failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(want)) {
		t.Fatalf("Pwritev high offset n=%d, want %d", n, len(want))
	}
	st, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if st.Size() != offset+int64(len(want)) {
		t.Fatalf("Pwritev high offset size=%d, want %d", st.Size(), offset+int64(len(want)))
	}

	got := make([]byte, len(want))
	rIov := []zcall.Iovec{
		{Base: &got[0], Len: uint64(len(got))},
	}
	n, errno = zcall.Preadv(fd, unsafe.Pointer(&rIov[0]), 1, offset)
	if errno != 0 {
		t.Fatalf("Preadv high offset failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(got)) {
		t.Fatalf("Preadv high offset n=%d, want %d", n, len(got))
	}
	if string(got) != string(want) {
		t.Fatalf("Preadv high offset read %q, want %q", got, want)
	}
}

func TestPreadv2Pwritev2(t *testing.T) {
	// Preadv2/Pwritev2 require seekable file descriptors (not sockets/pipes)
	// Test by verifying the syscall interface works - sockets return ESPIPE
	var fds [2]int32
	errno := zcall.Socketpair(zcall.AF_UNIX, zcall.SOCK_STREAM|zcall.SOCK_CLOEXEC, 0, &fds)
	if errno != 0 {
		t.Fatalf("Socketpair failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	buf := []byte("preadv2 pwritev2 test")
	iov := []zcall.Iovec{
		{Base: &buf[0], Len: uint64(len(buf))},
	}

	// Pwritev2 on socket should return ESPIPE (illegal seek)
	_, errno = zcall.Pwritev2(uintptr(fds[0]), unsafe.Pointer(&iov[0]), 1, 0, 0)
	if errno == 0 {
		t.Fatal("Pwritev2 on socket should fail with ESPIPE")
	}
	if zcall.Errno(errno) != zcall.ESPIPE {
		t.Fatalf("Pwritev2 errno = %v, want ESPIPE", zcall.Errno(errno))
	}

	// Preadv2 on socket should also return ESPIPE
	rbuf := make([]byte, 32)
	rIov := []zcall.Iovec{
		{Base: &rbuf[0], Len: uint64(len(rbuf))},
	}
	_, errno = zcall.Preadv2(uintptr(fds[1]), unsafe.Pointer(&rIov[0]), 1, 0, 0)
	if errno == 0 {
		t.Fatal("Preadv2 on socket should fail with ESPIPE")
	}
	if zcall.Errno(errno) != zcall.ESPIPE {
		t.Fatalf("Preadv2 errno = %v, want ESPIPE", zcall.Errno(errno))
	}
}

func TestPreadv2Pwritev2HighOffsetAndFlags(t *testing.T) {
	f, err := os.CreateTemp("", "zcall-preadv2-pwritev2-*")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	fd := uintptr(f.Fd())
	offset := int64(1<<33 + 456)
	want := []byte("preadv2-pwritev2 high offset")
	iov := []zcall.Iovec{
		{Base: &want[0], Len: uint64(len(want))},
	}

	n, errno := zcall.Pwritev2(fd, unsafe.Pointer(&iov[0]), 1, offset, 0)
	if errno != 0 {
		t.Fatalf("Pwritev2 high offset failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(want)) {
		t.Fatalf("Pwritev2 high offset n=%d, want %d", n, len(want))
	}
	st, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	wantSize := offset + int64(len(want))
	if st.Size() != wantSize {
		t.Fatalf("Pwritev2 high offset size=%d, want %d", st.Size(), wantSize)
	}

	got := make([]byte, len(want))
	rIov := []zcall.Iovec{
		{Base: &got[0], Len: uint64(len(got))},
	}
	n, errno = zcall.Preadv2(fd, unsafe.Pointer(&rIov[0]), 1, offset, 0)
	if errno != 0 {
		t.Fatalf("Preadv2 high offset failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(got)) {
		t.Fatalf("Preadv2 high offset n=%d, want %d", n, len(got))
	}
	if string(got) != string(want) {
		t.Fatalf("Preadv2 high offset read %q, want %q", got, want)
	}

	if _, err := f.Seek(0, 2); err != nil {
		t.Fatalf("Seek end failed: %v", err)
	}
	current := []byte(" current")
	currentIov := []zcall.Iovec{
		{Base: &current[0], Len: uint64(len(current))},
	}
	n, errno = zcall.Pwritev2(fd, unsafe.Pointer(&currentIov[0]), 1, -1, 0)
	if errno != 0 {
		t.Fatalf("Pwritev2 current offset failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(current)) {
		t.Fatalf("Pwritev2 current offset n=%d, want %d", n, len(current))
	}
	wantSize += int64(len(current))
	st, err = f.Stat()
	if err != nil {
		t.Fatalf("Stat after current offset failed: %v", err)
	}
	if st.Size() != wantSize {
		t.Fatalf("Pwritev2 current offset size=%d, want %d", st.Size(), wantSize)
	}

	appendData := []byte(" append")
	appendIov := []zcall.Iovec{
		{Base: &appendData[0], Len: uint64(len(appendData))},
	}
	n, errno = zcall.Pwritev2(fd, unsafe.Pointer(&appendIov[0]), 1, 0, zcall.RWF_APPEND)
	if errno != 0 {
		t.Fatalf("Pwritev2 RWF_APPEND failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(appendData)) {
		t.Fatalf("Pwritev2 RWF_APPEND n=%d, want %d", n, len(appendData))
	}
	wantSize += int64(len(appendData))
	st, err = f.Stat()
	if err != nil {
		t.Fatalf("Stat after RWF_APPEND failed: %v", err)
	}
	if st.Size() != wantSize {
		t.Fatalf("Pwritev2 RWF_APPEND size=%d, want %d", st.Size(), wantSize)
	}
}

func TestSpliceWithOffsets(t *testing.T) {
	// Create two pipes
	var pipe1 [2]int32
	errno := zcall.Pipe2(&pipe1, zcall.O_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(pipe1[0]))
	defer zcall.Close(uintptr(pipe1[1]))

	var pipe2 [2]int32
	errno = zcall.Pipe2(&pipe2, zcall.O_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(pipe2[0]))
	defer zcall.Close(uintptr(pipe2[1]))

	// Write data
	msg := []byte("splice offset test")
	n, errno := zcall.Write(uintptr(pipe1[1]), msg)
	if errno != 0 {
		t.Fatalf("Write failed: %v", zcall.Errno(errno))
	}

	// Splice with nil offsets (pipes don't support offsets)
	var offIn, offOut int64
	spliced, errno := zcall.Splice(uintptr(pipe1[0]), &offIn, uintptr(pipe2[1]), &offOut, uintptr(n), 0)
	// Pipes don't support offsets, so this may fail with ESPIPE - that's expected
	if errno != 0 && zcall.Errno(errno) != zcall.ESPIPE {
		t.Fatalf("Splice with offsets failed unexpectedly: %v", zcall.Errno(errno))
	}
	if errno == 0 && spliced != n {
		t.Fatalf("Splice returned %d, expected %d", spliced, n)
	}
}

func TestIoUringSetup(t *testing.T) {
	// io_uring_setup requires a params struct
	type ioUringParams struct {
		sqEntries    uint32
		cqEntries    uint32
		flags        uint32
		sqThreadCpu  uint32
		sqThreadIdle uint32
		features     uint32
		wqFd         uint32
		resv         [3]uint32
		sqOff        [10]uint32
		cqOff        [10]uint32
	}

	var params ioUringParams
	fd, errno := zcall.IoUringSetup(8, unsafe.Pointer(&params))
	if errno != 0 {
		// io_uring may not be available on all systems
		if zcall.Errno(errno) == zcall.ENOSYS {
			t.Skip("io_uring not supported on this kernel")
		}
		t.Fatalf("IoUringSetup failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	if fd == 0 || fd == ^uintptr(0) {
		t.Fatalf("IoUringSetup returned invalid fd: %d", fd)
	}

	// Test IoUringRegister - register eventfd
	efd, errno := zcall.Eventfd2(0, zcall.EFD_NONBLOCK|zcall.EFD_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Eventfd2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(efd)

	_, errno = zcall.IoUringRegister(fd, zcall.IORING_REGISTER_EVENTFD, unsafe.Pointer(&efd), 1)
	if errno != 0 {
		t.Fatalf("IoUringRegister failed: %v", zcall.Errno(errno))
	}

	// Test IoUringEnter - just submit 0 entries
	_, errno = zcall.IoUringEnter(fd, 0, 0, 0, nil, 0)
	if errno != 0 {
		t.Fatalf("IoUringEnter failed: %v", zcall.Errno(errno))
	}
}

func TestErrnoUnknown(t *testing.T) {
	// Test unknown errno value
	unknownErrno := zcall.Errno(9999)
	errStr := unknownErrno.Error()
	if errStr != "errno 9999" {
		t.Fatalf("Unknown errno string = %q, want %q", errStr, "errno 9999")
	}
}

func TestSendmmsgRecvmmsg(t *testing.T) {
	// Create UDP socket for multi-message test
	fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_DGRAM|zcall.SOCK_CLOEXEC, 0)
	if errno != 0 {
		t.Fatalf("Socket failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	// Bind to localhost
	addr := [16]byte{2, 0, 0, 0, 127, 0, 0, 1}
	errno = zcall.Bind(fd, unsafe.Pointer(&addr), 16)
	if errno != 0 {
		t.Fatalf("Bind failed: %v", zcall.Errno(errno))
	}

	// Get bound address
	var boundAddr [16]byte
	addrLen := uint32(16)
	errno = zcall.Getsockname(fd, unsafe.Pointer(&boundAddr), unsafe.Pointer(&addrLen))
	if errno != 0 {
		t.Fatalf("Getsockname failed: %v", zcall.Errno(errno))
	}

	// Prepare messages for sendmmsg
	msg1 := []byte("msg1")
	msg2 := []byte("msg2")

	iov1 := zcall.Iovec{Base: &msg1[0], Len: uint64(len(msg1))}
	iov2 := zcall.Iovec{Base: &msg2[0], Len: uint64(len(msg2))}

	mmsg := []zcall.Mmsghdr{
		{Hdr: zcall.Msghdr{Name: &boundAddr[0], Namelen: 16, Iov: &iov1, Iovlen: 1}},
		{Hdr: zcall.Msghdr{Name: &boundAddr[0], Namelen: 16, Iov: &iov2, Iovlen: 1}},
	}

	// Send multiple messages
	n, errno := zcall.Sendmmsg(fd, unsafe.Pointer(&mmsg[0]), 2, 0)
	if errno != 0 {
		t.Fatalf("Sendmmsg failed: %v", zcall.Errno(errno))
	}
	if n != 2 {
		t.Fatalf("Sendmmsg returned %d, expected 2", n)
	}

	// Prepare buffers for recvmmsg
	rbuf1 := make([]byte, 16)
	rbuf2 := make([]byte, 16)

	rIov1 := zcall.Iovec{Base: &rbuf1[0], Len: uint64(len(rbuf1))}
	rIov2 := zcall.Iovec{Base: &rbuf2[0], Len: uint64(len(rbuf2))}

	rmmsg := []zcall.Mmsghdr{
		{Hdr: zcall.Msghdr{Iov: &rIov1, Iovlen: 1}},
		{Hdr: zcall.Msghdr{Iov: &rIov2, Iovlen: 1}},
	}

	// Receive multiple messages
	n, errno = zcall.Recvmmsg(fd, unsafe.Pointer(&rmmsg[0]), 2, zcall.MSG_DONTWAIT, nil)
	if errno != 0 {
		t.Fatalf("Recvmmsg failed: %v", zcall.Errno(errno))
	}
	if n < 1 {
		t.Fatalf("Recvmmsg returned %d, expected at least 1", n)
	}
}

func TestErrno(t *testing.T) {
	tests := []struct {
		errno zcall.Errno
		want  string
	}{
		{0, "success"},                                     // zero errno
		{zcall.EPERM, "operation not permitted"},           // known errno
		{zcall.ENOENT, "no such file or directory"},        // known errno
		{zcall.EAGAIN, "resource temporarily unavailable"}, // known errno
		{zcall.EINVAL, "invalid argument"},                 // known errno
		{zcall.ECONNREFUSED, "connection refused"},         // known errno
	}

	for _, tt := range tests {
		if got := tt.errno.Error(); got != tt.want {
			t.Errorf("Errno(%d).Error() = %q, want %q", tt.errno, got, tt.want)
		}
	}
}

func TestErrnoLargeValues(t *testing.T) {
	// Test large errno values to fully exercise uitoa digit conversion
	tests := []struct {
		errno zcall.Errno
		want  string
	}{
		{200, "errno 200"},           // beyond errnoStrings array (3 digits)
		{999, "errno 999"},           // 3 digits
		{12345, "errno 12345"},       // 5 digits
		{100000, "errno 100000"},     // 6 digits
		{9999999, "errno 9999999"},   // 7 digits, all 9s
		{10000000, "errno 10000000"}, // 8 digits
	}

	for _, tt := range tests {
		got := tt.errno.Error()
		if got != tt.want {
			t.Errorf("Errno(%d).Error() = %q, want %q", tt.errno, got, tt.want)
		}
	}
}

func TestErrnoIs(t *testing.T) {
	err := zcall.Errno(zcall.EAGAIN)
	if !errors.Is(err, zcall.EAGAIN) {
		t.Error("errors.Is failed for EAGAIN")
	}
	if errors.Is(err, zcall.EINVAL) {
		t.Error("errors.Is should not match EINVAL")
	}
}

func TestErrnoTemporary(t *testing.T) {
	if !zcall.EAGAIN.Temporary() {
		t.Error("EAGAIN should be temporary")
	}
	if !zcall.EINTR.Temporary() {
		t.Error("EINTR should be temporary")
	}
	if !zcall.EINPROGRESS.Temporary() {
		t.Error("EINPROGRESS should be temporary")
	}
	if zcall.EINVAL.Temporary() {
		t.Error("EINVAL should not be temporary")
	}
}

func TestErrnoTimeout(t *testing.T) {
	if !zcall.EAGAIN.Timeout() {
		t.Error("EAGAIN should be timeout")
	}
	if !zcall.ETIMEDOUT.Timeout() {
		t.Error("ETIMEDOUT should be timeout")
	}
	if zcall.EINVAL.Timeout() {
		t.Error("EINVAL should not be timeout")
	}
}

func TestSyscallError(t *testing.T) {
	// Try to close an invalid fd
	errno := zcall.Close(^uintptr(0))
	if errno == 0 {
		t.Fatal("Close(-1) should fail")
	}
	if zcall.Errno(errno) != zcall.EBADF {
		t.Fatalf("Close(-1) errno = %v, want EBADF", zcall.Errno(errno))
	}
}

func TestReadWriteEmptyBuffer(t *testing.T) {
	fd, errno := zcall.Eventfd2(0, zcall.EFD_NONBLOCK|zcall.EFD_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Eventfd2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	// Empty write should return 0
	n, errno := zcall.Write(fd, nil)
	if errno != 0 {
		t.Fatalf("Write(nil) failed: %v", zcall.Errno(errno))
	}
	if n != 0 {
		t.Fatalf("Write(nil) returned %d, expected 0", n)
	}

	// Empty read should return 0
	n, errno = zcall.Read(fd, nil)
	if errno != 0 {
		t.Fatalf("Read(nil) failed: %v", zcall.Errno(errno))
	}
	if n != 0 {
		t.Fatalf("Read(nil) returned %d, expected 0", n)
	}
}

func BenchmarkSyscall4(b *testing.B) {
	fd, errno := zcall.Eventfd2(0, zcall.EFD_NONBLOCK|zcall.EFD_CLOEXEC)
	if errno != 0 {
		b.Fatalf("Eventfd2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	val := uint64(1)
	buf := (*[8]byte)(unsafe.Pointer(&val))[:]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zcall.Write(fd, buf)
		zcall.Read(fd, buf)
	}
}

func BenchmarkSyscall4Latency(b *testing.B) {
	fd, errno := zcall.Eventfd2(0, zcall.EFD_NONBLOCK|zcall.EFD_CLOEXEC)
	if errno != 0 {
		b.Fatalf("Eventfd2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	val := uint64(1)
	buf := (*[8]byte)(unsafe.Pointer(&val))[:]

	var totalNs int64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		zcall.Write(fd, buf)
		zcall.Read(fd, buf)
		totalNs += time.Since(start).Nanoseconds()
	}
	b.ReportMetric(float64(totalNs)/float64(b.N), "ns/op")
}

func BenchmarkEventfd2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fd, errno := zcall.Eventfd2(0, zcall.EFD_NONBLOCK|zcall.EFD_CLOEXEC)
		if errno != 0 {
			b.Fatalf("Eventfd2 failed: %v", zcall.Errno(errno))
		}
		zcall.Close(fd)
	}
}

func BenchmarkSocket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
		if errno != 0 {
			b.Fatalf("Socket failed: %v", zcall.Errno(errno))
		}
		zcall.Close(fd)
	}
}

// BenchmarkZcallWrite benchmarks zcall.Write to a pipe.
// This measures the raw syscall latency without Go runtime hooks.
func BenchmarkZcallWrite(b *testing.B) {
	var fds [2]int32
	errno := zcall.Pipe2(&fds, zcall.O_NONBLOCK|zcall.O_CLOEXEC)
	if errno != 0 {
		b.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	buf := make([]byte, 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zcall.Write(uintptr(fds[1]), buf)
		zcall.Read(uintptr(fds[0]), buf)
	}
}

// BenchmarkZcallWriteLatency measures per-operation latency for zcall.Write.
func BenchmarkZcallWriteLatency(b *testing.B) {
	var fds [2]int32
	errno := zcall.Pipe2(&fds, zcall.O_NONBLOCK|zcall.O_CLOEXEC)
	if errno != 0 {
		b.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(fds[0]))
	defer zcall.Close(uintptr(fds[1]))

	buf := make([]byte, 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		zcall.Write(uintptr(fds[1]), buf)
		zcall.Read(uintptr(fds[0]), buf)
		b.ReportMetric(float64(time.Since(start).Nanoseconds()), "ns/roundtrip")
	}
}

func TestSignalfd4(t *testing.T) {
	// Create a signal mask (64-bit mask for signals)
	var mask uint64 = 1 << (10 - 1) // SIGUSR1 (signal 10)

	// Create a new signalfd
	fd, errno := zcall.Signalfd4(^uintptr(0), unsafe.Pointer(&mask), 8, zcall.SFD_NONBLOCK|zcall.SFD_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Signalfd4 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	if fd == 0 || fd == ^uintptr(0) {
		t.Fatalf("Signalfd4 returned invalid fd: %d", fd)
	}

	// Try to read - should return EAGAIN since no signal is pending
	var sigInfo [128]byte // signalfd_siginfo is 128 bytes
	_, errno = zcall.Read(fd, sigInfo[:])
	if errno == 0 {
		t.Fatal("Read should have returned EAGAIN (no signal pending)")
	}
	if zcall.Errno(errno) != zcall.EAGAIN {
		t.Fatalf("Read errno = %v, want EAGAIN", zcall.Errno(errno))
	}

	// Modify the signalfd with a new mask
	var newMask uint64 = 1 << (12 - 1) // SIGUSR2 (signal 12)
	newFd, errno := zcall.Signalfd4(fd, unsafe.Pointer(&newMask), 8, 0)
	if errno != 0 {
		t.Fatalf("Signalfd4 (modify) failed: %v", zcall.Errno(errno))
	}
	if newFd != fd {
		t.Fatalf("Signalfd4 (modify) returned different fd: %d vs %d", newFd, fd)
	}
}

func TestPidfdOpen(t *testing.T) {
	// Open a pidfd for the current process
	pid := uintptr(1) // init process (always exists)
	fd, errno := zcall.PidfdOpen(pid, 0)
	if errno != 0 {
		// May fail without CAP_SYS_PTRACE for other processes
		if zcall.Errno(errno) == zcall.EPERM || zcall.Errno(errno) == zcall.ENOSYS {
			t.Skipf("PidfdOpen not permitted or not supported: %v", zcall.Errno(errno))
		}
		t.Fatalf("PidfdOpen failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	if fd == 0 || fd == ^uintptr(0) {
		t.Fatalf("PidfdOpen returned invalid fd: %d", fd)
	}
}

func TestPidfdOpenInvalidPid(t *testing.T) {
	// Try to open a pidfd for an invalid PID
	_, errno := zcall.PidfdOpen(^uintptr(0), 0)
	if errno == 0 {
		t.Fatal("PidfdOpen should fail for invalid PID")
	}
	// Should return ESRCH (no such process) or EINVAL
	e := zcall.Errno(errno)
	if e != zcall.ESRCH && e != zcall.EINVAL {
		t.Fatalf("PidfdOpen errno = %v, want ESRCH or EINVAL", e)
	}
}

func TestPidfdGetfd(t *testing.T) {
	// Create a pipe to have a valid fd
	var pipeFds [2]int32
	errno := zcall.Pipe2(&pipeFds, zcall.O_CLOEXEC)
	if errno != 0 {
		t.Fatalf("Pipe2 failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(uintptr(pipeFds[0]))
	defer zcall.Close(uintptr(pipeFds[1]))

	// Open pidfd for init process
	pidfd, errno := zcall.PidfdOpen(1, 0)
	if errno != 0 {
		if zcall.Errno(errno) == zcall.EPERM || zcall.Errno(errno) == zcall.ENOSYS {
			t.Skipf("PidfdOpen not permitted or not supported: %v", zcall.Errno(errno))
		}
		t.Fatalf("PidfdOpen failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(pidfd)

	// Try to get fd from another process - should fail without CAP_SYS_PTRACE
	_, errno = zcall.PidfdGetfd(pidfd, 0, 0)
	if errno == 0 {
		t.Log("PidfdGetfd succeeded (may have CAP_SYS_PTRACE)")
	} else {
		// Expected to fail with EPERM
		e := zcall.Errno(errno)
		if e != zcall.EPERM && e != zcall.ENOSYS && e != zcall.EBADF {
			t.Fatalf("PidfdGetfd errno = %v, want EPERM, ENOSYS, or EBADF", e)
		}
	}
}

func TestPidfdGetfdInvalidFd(t *testing.T) {
	// Try with invalid pidfd
	_, errno := zcall.PidfdGetfd(^uintptr(0), 0, 0)
	if errno == 0 {
		t.Fatal("PidfdGetfd should fail with invalid pidfd")
	}
	if zcall.Errno(errno) != zcall.EBADF {
		t.Fatalf("PidfdGetfd errno = %v, want EBADF", zcall.Errno(errno))
	}
}

func TestPidfdSendSignal(t *testing.T) {
	// Open pidfd for init process
	pidfd, errno := zcall.PidfdOpen(1, 0)
	if errno != 0 {
		if zcall.Errno(errno) == zcall.EPERM || zcall.Errno(errno) == zcall.ENOSYS {
			t.Skipf("PidfdOpen not permitted or not supported: %v", zcall.Errno(errno))
		}
		t.Fatalf("PidfdOpen failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(pidfd)

	// Try to send signal 0 (null signal - checks if process exists)
	errno = zcall.PidfdSendSignal(pidfd, 0, nil, 0)
	if errno != 0 {
		// May fail with EPERM if we don't have permission to signal init
		e := zcall.Errno(errno)
		if e != zcall.EPERM {
			t.Fatalf("PidfdSendSignal errno = %v, want success or EPERM", e)
		}
		t.Logf("PidfdSendSignal returned EPERM (expected without CAP_KILL)")
	}
}

func TestPidfdSendSignalInvalidFd(t *testing.T) {
	// Try with invalid pidfd
	errno := zcall.PidfdSendSignal(^uintptr(0), 0, nil, 0)
	if errno == 0 {
		t.Fatal("PidfdSendSignal should fail with invalid pidfd")
	}
	if zcall.Errno(errno) != zcall.EBADF {
		t.Fatalf("PidfdSendSignal errno = %v, want EBADF", zcall.Errno(errno))
	}
}

func TestMemfdCreate(t *testing.T) {
	// Create a memfd with a name
	name := []byte("test_memfd\x00") // null-terminated name
	fd, errno := zcall.MemfdCreate(unsafe.Pointer(&name[0]), zcall.MFD_CLOEXEC)
	if errno != 0 {
		if zcall.Errno(errno) == zcall.ENOSYS {
			t.Skip("memfd_create not supported on this kernel")
		}
		t.Fatalf("MemfdCreate failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	if fd == 0 || fd == ^uintptr(0) {
		t.Fatalf("MemfdCreate returned invalid fd: %d", fd)
	}

	// Write to memfd
	data := []byte("hello memfd")
	n, errno := zcall.Write(fd, data)
	if errno != 0 {
		t.Fatalf("Write to memfd failed: %v", zcall.Errno(errno))
	}
	if n != uintptr(len(data)) {
		t.Fatalf("Write returned %d, expected %d", n, len(data))
	}
}

func TestMemfdCreateWithSealing(t *testing.T) {
	// Create a memfd with sealing support
	name := []byte("test_sealed\x00")
	fd, errno := zcall.MemfdCreate(unsafe.Pointer(&name[0]), zcall.MFD_CLOEXEC|zcall.MFD_ALLOW_SEALING)
	if errno != 0 {
		if zcall.Errno(errno) == zcall.ENOSYS {
			t.Skip("memfd_create not supported on this kernel")
		}
		t.Fatalf("MemfdCreate failed: %v", zcall.Errno(errno))
	}
	defer zcall.Close(fd)

	if fd == 0 || fd == ^uintptr(0) {
		t.Fatalf("MemfdCreate returned invalid fd: %d", fd)
	}
}

func TestErrnoZero(t *testing.T) {
	// Test that uitoa handles 0 correctly (covers the val == 0 branch)
	e := zcall.Errno(0)
	if e.Error() != "success" {
		t.Errorf("Errno(0).Error() = %q, want %q", e.Error(), "success")
	}
}

func TestErrnoEmptyString(t *testing.T) {
	// Test errno values that may have empty strings in the array
	// These should fall through to the "errno N" format
	// Find an errno that's in range but has no string defined
	for i := 1; i < 200; i++ {
		e := zcall.Errno(i)
		str := e.Error()
		// Verify it returns something reasonable
		if str == "" {
			t.Errorf("Errno(%d).Error() returned empty string", i)
		}
	}
}

func BenchmarkMemfdCreate(b *testing.B) {
	name := []byte("bench_memfd\x00")
	for i := 0; i < b.N; i++ {
		fd, errno := zcall.MemfdCreate(unsafe.Pointer(&name[0]), zcall.MFD_CLOEXEC)
		if errno != 0 {
			if zcall.Errno(errno) == zcall.ENOSYS {
				b.Skip("memfd_create not supported")
			}
			b.Fatalf("MemfdCreate failed: %v", zcall.Errno(errno))
		}
		zcall.Close(fd)
	}
}

func BenchmarkPidfdOpen(b *testing.B) {
	// First check if pidfd_open is supported
	fd, errno := zcall.PidfdOpen(1, 0)
	if errno != 0 {
		b.Skipf("PidfdOpen not available: %v", zcall.Errno(errno))
	}
	zcall.Close(fd)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fd, _ := zcall.PidfdOpen(1, 0)
		zcall.Close(fd)
	}
}

func BenchmarkSignalfd4(b *testing.B) {
	var mask uint64 = 1 << 9 // SIGUSR1
	for i := 0; i < b.N; i++ {
		fd, errno := zcall.Signalfd4(^uintptr(0), unsafe.Pointer(&mask), 8, zcall.SFD_NONBLOCK|zcall.SFD_CLOEXEC)
		if errno != 0 {
			b.Fatalf("Signalfd4 failed: %v", zcall.Errno(errno))
		}
		zcall.Close(fd)
	}
}
