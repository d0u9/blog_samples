package main

import (
	"fmt"
	"log"
	"path"
	"github.com/cilium/ebpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go prog bpf/bpf.c -- -I../include

const ETH_NAME = "eth0"
const BPF_FS    = "/sys/fs/bpf"

type Map3Spec struct {
    inner   *ebpf.MapSpec  `ebpf:"map_3"`
}

func LoadMap3FromFile(specs *ebpf.CollectionSpec) *ebpf.Map {
	opt := ebpf.MapOptions{
			PinPath: path.Join(BPF_FS, "map3_dir"),
	}
    mp, err := ebpf.NewMapWithOptions(specs.Maps["map_3"], opt);
    if err != nil {
        log.Fatalf("cannot load map3: %v", err)
    }
    return mp
}

func LoadMap4FromFile(specs *ebpf.CollectionSpec) *ebpf.Map {
	opt := ebpf.MapOptions{
			PinPath: path.Join(BPF_FS, "map4_dir"),
	}
    mp, err := ebpf.NewMapWithOptions(specs.Maps["map_4"], opt);
    if err != nil {
        log.Fatalf("cannot load map4: %v", err)
    }
    return mp
}

func main() {
    specs, err := loadProg()
    if err != nil {
        log.Fatalf("load prog spec failed: %v", err)
    }

    map3 := LoadMap3FromFile(specs)
    map4 := LoadMap4FromFile(specs)

    fmt.Println(map3, map4)
    rewriteMap := map[string]*ebpf.Map {
        "map_3": map3,
        "map_4": map4,
    }

    fmt.Printf("%v\n", specs);
    if err := specs.RewriteMaps(rewriteMap); err != nil {
        log.Fatalf("rewrite map failed: %v\n", err)
    }

    fmt.Printf("%v\n", specs);


    objs := progPrograms{}
    opts := ebpf.CollectionOptions {
        Maps: ebpf.MapOptions {
			PinPath: BPF_FS,
        },
    }
    if err := specs.LoadAndAssign(&objs, &opts); err != nil {
        log.Fatalf("LoadAndAssign Failed: %v\n", err)
    }

	var buffer string
	fmt.Println("Press any key to continue")
	fmt.Scanf("%s", &buffer)
}
