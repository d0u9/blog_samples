#include <linux/bpf.h>
#include <linux/pkt_cls.h>

#include "helpers.h"
#include "maps/map_1.h"
#include "maps/map_3.h"
#include "maps/map_4.h"

SEC("kprobe/sys_execve")
int prog_a(struct __sk_buff *skb)
{
	__u32 key = 2;
	__u64 initval = 1, *valp = NULL;
	valp = (__u64 *)bpf_map_lookup_elem(&map_1, &key);
	if (!valp) {
		bpf_map_update_elem(&map_1, &key, &initval, BPF_ANY);
	} else {
		__sync_fetch_and_add(valp, 1);
	}

	valp = (__u64 *)bpf_map_lookup_elem(&map_3, &key);
	if (!valp) {
		bpf_map_update_elem(&map_3, &key, &initval, BPF_ANY);
	} else {
		__sync_fetch_and_add(valp, 1);
	}

	valp = (__u64 *)bpf_map_lookup_elem(&map_4, &key);
	if (!valp) {
		bpf_map_update_elem(&map_4, &key, &initval, BPF_ANY);
	} else {
		__sync_fetch_and_add(valp, 1);
	}

	return 0;
}

char __license[] SEC("license") = "GPL";


