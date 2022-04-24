#ifndef __MAP_4_H
#define __MAP_4_H

#include "helpers.h"

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__type(key, __u32);
	__type(value, __u32);
	__uint(max_entries, 10);
	__uint(pinning, LIBBPF_PIN_BY_NAME);
} map_4 SEC(".maps");

#endif /* __MAP_4_H */
