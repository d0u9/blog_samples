package main

import (
	"fmt"
	"log"

	"github.com/cilium/ebpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go prog bpf/bpf.c -- -I../include

const ETH_NAME = "eth0"
const PIN_PATH = "/sys/fs/bpf/"

func main() {
	objs_a := progObjects{}
	opta := ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: PIN_PATH,
		},
	}
	if err := loadProgObjects(&objs_a, &opta); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs_a.Close()

	var buffer string
	fmt.Println("Press any key to continue")
	fmt.Scanf("%s", &buffer)
}
