package socks

// #include "socks.h"
// #include <stdio.h>
// #include <stdlib.h>
// #include <arpa/inet.h>
import "C"

import (
	"errors"
	"fmt"
	"net"
	"unsafe"
)

func CreateDeviceBoundUDPSocket(addr *net.IP, port uint16, device string) (fd int, err error) {
	devcstr := C.CString(device)
    //fmt.Printf("creating device for %s", device)
	defer C.free(unsafe.Pointer(devcstr))
	fd = int(C.createDeviceBoundUDPSocket(C.uint32_t(ipToInt(addr)), C.uint16_t(port), devcstr))
	if fd < 0 {
		return fd, errors.New(fmt.Sprintf("Got return code %d when creating bound socket.", fd))
	}

	return
}

// not thread safe
func WriteToUDP(fd int, buf []byte, dest *net.UDPAddr) (int, error) {
    C.putchar(89);
    //fmt.Printf("WriteToUDP (go): fd=%d, buf len=%d, dest=%s\n", fd, len(buf), dest);
	res := C.writeToUDP(C.int(fd), unsafe.Pointer(&buf[0]),
		C.size_t(len(buf)), C.uint32_t(ipToInt(&dest.IP)),
		C.uint16_t(dest.Port))
	if res < 0 {
		return int(res), errors.New(fmt.Sprintf("Got error code %d from sendto", res))
	}
	return int(res), nil
}

// helper function for converting ipv4 addr to uint32_t
// doesn't work for ipv6
func ipToInt(ip *net.IP) C.uint32_t {
	addr4 := ip.To4()
	var intaddr uint32
	for i := 0; i < 4; i++ {
		intaddr = intaddr << 8
		intaddr = intaddr + uint32(addr4[i])
	}
    //fmt.Printf("Conv %s->%x\n", ip, intaddr)
	return C.uint32_t(intaddr)
}
