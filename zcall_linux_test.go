// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

//go:build linux

package zcall

import (
	"reflect"
	"testing"
	"unsafe"
)

func encodeAttr(attrType uint16, value []byte) []byte {
	length := int(unsafe.Sizeof(rtAttr{})) + len(value)
	aligned := rtaAlignOf(length)
	buf := make([]byte, aligned)
	attr := (*rtAttr)(unsafe.Pointer(&buf[0]))
	attr.Len = uint16(length)
	attr.Type = attrType
	copy(buf[int(unsafe.Sizeof(rtAttr{})):length], value)
	return buf
}

func TestParseLinkInfoParsesAttributes(t *testing.T) {
	payload := make([]byte, int(unsafe.Sizeof(ifInfomsg{})))
	ifim := (*ifInfomsg)(unsafe.Pointer(&payload[0]))
	ifim.Index = 7
	ifim.Flags = IFF_UP | IFF_RUNNING | IFF_MULTICAST

	mtu := uint32(9000)
	payload = append(payload,
		encodeAttr(IFLA_IFNAME, append([]byte("eth0"), 0))...,
	)
	payload = append(payload,
		encodeAttr(IFLA_MTU, unsafe.Slice((*byte)(unsafe.Pointer(&mtu)), 4))...,
	)
	payload = append(payload,
		encodeAttr(IFLA_ADDRESS, []byte{0x02, 0x42, 0xac, 0x11, 0x00, 0x02})...,
	)

	link, ok, err := parseLinkInfo(payload)
	if err != nil {
		t.Fatalf("parseLinkInfo() error = %v", err)
	}
	if !ok {
		t.Fatal("parseLinkInfo() ok = false, want true")
	}
	if link.Index != 7 {
		t.Fatalf("parseLinkInfo() index = %d, want 7", link.Index)
	}
	if link.Name != "eth0" {
		t.Fatalf("parseLinkInfo() name = %q, want %q", link.Name, "eth0")
	}
	if link.MTU != int(mtu) {
		t.Fatalf("parseLinkInfo() mtu = %d, want %d", link.MTU, mtu)
	}
	if link.Flags != ifim.Flags {
		t.Fatalf("parseLinkInfo() flags = %#x, want %#x", link.Flags, ifim.Flags)
	}
	wantAddr := []byte{0x02, 0x42, 0xac, 0x11, 0x00, 0x02}
	if !reflect.DeepEqual(link.HardwareAddr, wantAddr) {
		t.Fatalf("parseLinkInfo() hardware addr = %v, want %v", link.HardwareAddr, wantAddr)
	}
}

func TestParseLinkInfoRejectsMalformedAttribute(t *testing.T) {
	payload := make([]byte, int(unsafe.Sizeof(ifInfomsg{}))+int(unsafe.Sizeof(rtAttr{})))
	ifim := (*ifInfomsg)(unsafe.Pointer(&payload[0]))
	ifim.Index = 1
	attr := (*rtAttr)(unsafe.Pointer(&payload[unsafe.Sizeof(ifInfomsg{})]))
	attr.Len = uint16(unsafe.Sizeof(rtAttr{}) - 1)
	attr.Type = IFLA_IFNAME

	_, ok, err := parseLinkInfo(payload)
	if err == nil {
		t.Fatal("parseLinkInfo() error = nil, want error")
	}
	if ok {
		t.Fatal("parseLinkInfo() ok = true, want false")
	}
	if err != Errno(EINVAL) {
		t.Fatalf("parseLinkInfo() error = %v, want %v", err, Errno(EINVAL))
	}
}

func TestInterfaceLookupValidation(t *testing.T) {
	if _, err := InterfaceByName(""); err != Errno(EINVAL) {
		t.Fatalf("InterfaceByName(\"\") error = %v, want %v", err, Errno(EINVAL))
	}
	if _, err := InterfaceByName("0123456789abcdef"); err != Errno(EINVAL) {
		t.Fatalf("InterfaceByName(long) error = %v, want %v", err, Errno(EINVAL))
	}
	if _, err := InterfaceByIndex(0); err != Errno(EINVAL) {
		t.Fatalf("InterfaceByIndex(0) error = %v, want %v", err, Errno(EINVAL))
	}
}

func TestInterfacesLookupRoundTrip(t *testing.T) {
	links, err := Interfaces()
	if err != nil {
		t.Fatalf("Interfaces() error = %v", err)
	}
	if len(links) == 0 {
		t.Fatal("Interfaces() returned no links")
	}

	link := links[0]
	if link.Name == "" {
		t.Fatal("Interfaces() returned link with empty name")
	}
	if link.Index <= 0 {
		t.Fatalf("Interfaces() returned invalid index %d", link.Index)
	}

	byName, err := InterfaceByName(link.Name)
	if err != nil {
		t.Fatalf("InterfaceByName(%q) error = %v", link.Name, err)
	}
	if byName.Index != link.Index {
		t.Fatalf("InterfaceByName(%q) index = %d, want %d", link.Name, byName.Index, link.Index)
	}

	byIndex, err := InterfaceByIndex(link.Index)
	if err != nil {
		t.Fatalf("InterfaceByIndex(%d) error = %v", link.Index, err)
	}
	if byIndex.Name != link.Name {
		t.Fatalf("InterfaceByIndex(%d) name = %q, want %q", link.Index, byIndex.Name, link.Name)
	}

	index, err := IfNameToIndex(link.Name)
	if err != nil {
		t.Fatalf("IfNameToIndex(%q) error = %v", link.Name, err)
	}
	if index != uint32(link.Index) {
		t.Fatalf("IfNameToIndex(%q) = %d, want %d", link.Name, index, link.Index)
	}
}