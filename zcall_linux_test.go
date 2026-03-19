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

func TestParseLinkParsesAttributes(t *testing.T) {
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

	link, ok, err := parseLink(payload)
	if err != nil {
		t.Fatalf("parseLink() error = %v", err)
	}
	if !ok {
		t.Fatal("parseLink() ok = false, want true")
	}
	if link.Index != 7 {
		t.Fatalf("parseLink() index = %d, want 7", link.Index)
	}
	if link.Name != "eth0" {
		t.Fatalf("parseLink() name = %q, want %q", link.Name, "eth0")
	}
	if link.MTU != int(mtu) {
		t.Fatalf("parseLink() mtu = %d, want %d", link.MTU, mtu)
	}
	if link.Flags != ifim.Flags {
		t.Fatalf("parseLink() flags = %#x, want %#x", link.Flags, ifim.Flags)
	}
	wantAddr := []byte{0x02, 0x42, 0xac, 0x11, 0x00, 0x02}
	if !reflect.DeepEqual(link.HardwareAddr, wantAddr) {
		t.Fatalf("parseLink() hardware addr = %v, want %v", link.HardwareAddr, wantAddr)
	}
}

func TestParseLinkRejectsMalformedAttribute(t *testing.T) {
	payload := make([]byte, int(unsafe.Sizeof(ifInfomsg{}))+int(unsafe.Sizeof(rtAttr{})))
	ifim := (*ifInfomsg)(unsafe.Pointer(&payload[0]))
	ifim.Index = 1
	attr := (*rtAttr)(unsafe.Pointer(&payload[unsafe.Sizeof(ifInfomsg{})]))
	attr.Len = uint16(unsafe.Sizeof(rtAttr{}) - 1)
	attr.Type = IFLA_IFNAME

	_, ok, err := parseLink(payload)
	if err == nil {
		t.Fatal("parseLink() error = nil, want error")
	}
	if ok {
		t.Fatal("parseLink() ok = true, want false")
	}
	if err != Errno(EINVAL) {
		t.Fatalf("parseLink() error = %v, want %v", err, Errno(EINVAL))
	}
}

func TestLinkLookupValidation(t *testing.T) {
	if _, err := LinkByName(""); err != Errno(EINVAL) {
		t.Fatalf("LinkByName(\"\") error = %v, want %v", err, Errno(EINVAL))
	}
	if _, err := LinkByName("0123456789abcdef"); err != Errno(EINVAL) {
		t.Fatalf("LinkByName(long) error = %v, want %v", err, Errno(EINVAL))
	}
	if _, err := LinkByIndex(0); err != Errno(EINVAL) {
		t.Fatalf("LinkByIndex(0) error = %v, want %v", err, Errno(EINVAL))
	}
}

func TestLinksLookupRoundTrip(t *testing.T) {
	links, err := Links()
	if err != nil {
		t.Fatalf("Links() error = %v", err)
	}
	if len(links) == 0 {
		t.Fatal("Links() returned no links")
	}

	link := links[0]
	if link.Name == "" {
		t.Fatal("Links() returned link with empty name")
	}
	if link.Index <= 0 {
		t.Fatalf("Links() returned invalid index %d", link.Index)
	}

	byName, err := LinkByName(link.Name)
	if err != nil {
		t.Fatalf("LinkByName(%q) error = %v", link.Name, err)
	}
	if byName.Index != link.Index {
		t.Fatalf("LinkByName(%q) index = %d, want %d", link.Name, byName.Index, link.Index)
	}

	byIndex, err := LinkByIndex(link.Index)
	if err != nil {
		t.Fatalf("LinkByIndex(%d) error = %v", link.Index, err)
	}
	if byIndex.Name != link.Name {
		t.Fatalf("LinkByIndex(%d) name = %q, want %q", link.Index, byIndex.Name, link.Name)
	}

	index, err := IfNameToIndex(link.Name)
	if err != nil {
		t.Fatalf("IfNameToIndex(%q) error = %v", link.Name, err)
	}
	if index != uint32(link.Index) {
		t.Fatalf("IfNameToIndex(%q) = %d, want %d", link.Name, index, link.Index)
	}
}
