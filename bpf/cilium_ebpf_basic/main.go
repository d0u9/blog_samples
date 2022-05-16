package main

import (
	"log"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go tc bpf/bpf.c -- -I./bpf

const ETH_NAME = "eth0"

func main() {
	var err error

	objs := tcObjects{}
	if err := loadTcObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	progFd := objs.tcPrograms.TcMain.FD()

	eth0, err := netlink.LinkByName(ETH_NAME)
	if err != nil {
		log.Fatalf("cannot find %s: %v", ETH_NAME, err)
	}

	attrs := netlink.QdiscAttrs{
		LinkIndex: eth0.Attrs().Index,
		Handle:    netlink.MakeHandle(0xffff, 0),
		Parent:    netlink.HANDLE_CLSACT,
	}

	qdisc := &netlink.GenericQdisc{
		QdiscAttrs: attrs,
		QdiscType:  "clsact",
	}

	if err := netlink.QdiscAdd(qdisc); err != nil {
		log.Fatalf("cannot add clsact qdisc: %v", err)
	}

	filterAttrs := netlink.FilterAttrs{
		LinkIndex: eth0.Attrs().Index,
		Parent:    netlink.HANDLE_MIN_INGRESS,
		Handle:    netlink.MakeHandle(0, 1),
		Protocol:  unix.ETH_P_ALL,
		Priority:  1,
	}

	filter := &netlink.BpfFilter{
		FilterAttrs:  filterAttrs,
		Fd:           progFd,
		Name:         "hi-tc",
		DirectAction: true,
	}

	if err := netlink.FilterAdd(filter); err != nil {
		log.Fatalf("cannot attach bpf object to filter: %v", err)
	}
}
