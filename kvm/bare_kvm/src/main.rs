use std::io::Write;
use std::ptr::null_mut;
use std::slice;

use kvm_ioctls::{Kvm, VcpuExit};

fn main() {
    let kvm = Kvm::new().unwrap();

    let vm = kvm.create_vm().unwrap();

    let mem_size = 0x4000;
    let load_addr: *mut u8 = unsafe {
        libc::mmap(
            null_mut(),
            mem_size,
            libc::PROT_READ | libc::PROT_WRITE,
            libc::MAP_ANONYMOUS | libc::MAP_SHARED | libc::MAP_NORESERVE,
            -1,
            0,
        ) as *mut u8
    };

    let slot = 0;
    let guest_addr = 0x1000;
    let mem_region = kvm_bindings::kvm_userspace_memory_region {
        slot,
        guest_phys_addr: guest_addr,
        memory_size: mem_size as u64,
        userspace_addr: load_addr as u64,
        flags: kvm_bindings::KVM_MEM_LOG_DIRTY_PAGES,
    };
    unsafe { vm.set_user_memory_region(mem_region).unwrap() };

    let vcpu_fd = vm.create_vcpu(0).unwrap();

    let mut vcpu_sregs = vcpu_fd.get_sregs().unwrap();
    vcpu_sregs.cs.base = 0;
    vcpu_sregs.cs.selector = 0;
    vcpu_fd.set_sregs(&vcpu_sregs).unwrap();

    let mut vcpu_regs = vcpu_fd.get_regs().unwrap();
    vcpu_regs.rip = guest_addr;
    vcpu_regs.rax = 2;
    vcpu_regs.rbx = 3;
    vcpu_regs.rflags = 2;
    vcpu_fd.set_regs(&vcpu_regs).unwrap();

    // Calulate 0 + 1 + 2 + ... + 20 = 210
    let asm_code = &[
        0x31, 0xc0, /* xor    %ax,    %ax */
        0x31, 0xdb, /* xor    %bx,    %bx */
        // LOOP:
        0x01, 0xd8, /* add    %bx,    %ax */
        0x83, 0xc3, 0x01, /* add    $1,     %bx */
        0x83, 0xfb, 0x14, /* cmp    $20,    %bx */
        0x7e, 0xf6, /* jle    LOOP        */
        0xba, 0xff, 0x0e, /* mov    $0x217, %dx */
        0xee, /* out    %al,    %dx */
        0xf4, /* hlt                */
    ];

    unsafe {
        let mut slice = slice::from_raw_parts_mut(load_addr, mem_size);
        let _ = slice.write(asm_code).unwrap();
    }

    loop {
        match vcpu_fd.run().expect("run failed") {
            VcpuExit::Hlt => {
                println!("VM halts");
                break;
            }
            VcpuExit::IoOut(addr, data) => {
                println!("Address: {:#x}, Data: {:?}", addr, data);
            }
            r => {
                println!("Unknown halts: {:?}", r);
                break;
            }
        }
    }
}
