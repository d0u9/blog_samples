#include <linux/bpf.h>
#include <linux/pkt_cls.h>
#include "include/helpers.h"

SEC("tc_prog")
int tc_main(struct __sk_buff *skb)
{
	char hello_str[] = "hello egress pkt";
	bpf_trace_printk(hello_str, sizeof(hello_str));
	return TC_ACT_OK;
}

char __license[] SEC("license") = "GPL";
