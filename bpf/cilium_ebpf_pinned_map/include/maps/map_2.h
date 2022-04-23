#ifndef __MAP_2_H
#define __MAP_2_H

#include "helpers.h"

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__type(key, __u32);
	__type(value, __u64);
	__uint(max_entries, 10);
	__uint(pinning, LIBBPF_PIN_BY_NAME);
} map_2 SEC(".maps");

#endif /* __MAP_2_H */

