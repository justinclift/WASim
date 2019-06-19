(module
  (type $t0 (func (param i32)))
  (type $t1 (func (result i32)))
  (type $t2 (func (param i32 i32 i32) (result i32)))
  (type $t3 (func))
  (type $t4 (func (param i32 i32 i32)))
  (type $t5 (func (param i32 i32)))
  (import "env" "io_get_stdout" (func $io_get_stdout (type $t1)))
  (import "env" "resource_write" (func $resource_write (type $t2)))
  (func $__wasm_call_ctors (type $t3))
  (func $_start (type $t3)
    i32.const 0
    call $io_get_stdout
    i32.store offset=1024)
  (func $runtime.activateTask (type $t4) (param $p0 i32) (param $p1 i32) (param $p2 i32)
    (local $l3 i32)
    block $B0
      block $B1
        block $B2
          local.get $p0
          i32.eqz
          br_if $B2
          local.get $p0
          i32.load
          i32.eqz
          br_if $B1
          i32.const 0
          i32.load offset=1032
          local.tee $l3
          i32.eqz
          br_if $B0
          i32.const 0
          local.get $p0
          i32.store offset=1032
          local.get $l3
          local.get $p0
          i32.store offset=8
        end
        return
      end
      local.get $p0
      local.get $p0
      i32.load offset=4
      call_indirect (type $t0) $T0
      return
    end
    i32.const 0
    local.get $p0
    i32.store offset=1028
    i32.const 0
    local.get $p0
    i32.store offset=1032)
  (func $cwa_main (type $t3)
    i32.const 0
    call $io_get_stdout
    i32.store offset=1024
    i32.const 1094
    i32.const 12
    call $runtime.printstring
    call $runtime.printnl)
  (func $runtime.printstring (type $t5) (param $p0 i32) (param $p1 i32)
    (local $l2 i32)
    i32.const 0
    local.set $l2
    block $B0
      loop $L1
        local.get $l2
        local.get $p1
        i32.ge_s
        br_if $B0
        local.get $p0
        local.get $l2
        i32.add
        i32.load8_u
        call $runtime.putchar
        local.get $l2
        i32.const 1
        i32.add
        local.set $l2
        br $L1
      end
    end)
  (func $runtime.printnl (type $t3)
    i32.const 13
    call $runtime.putchar
    i32.const 10
    call $runtime.putchar)
  (func $runtime.getFuncPtr (type $t3)
    call $runtime.nilPanic
    unreachable)
  (func $runtime.nilPanic (type $t3)
    call $runtime.runtimePanic
    unreachable)
  (func $memset (type $t2) (param $p0 i32) (param $p1 i32) (param $p2 i32) (result i32)
    (local $l3 i32) (local $l4 i32)
    i32.const 0
    local.set $l3
    block $B0
      block $B1
        loop $L2
          local.get $l3
          local.get $p2
          i32.ge_u
          br_if $B1
          local.get $p0
          local.get $l3
          i32.add
          local.tee $l4
          i32.eqz
          br_if $B0
          local.get $l4
          local.get $p1
          i32.store8
          local.get $l3
          i32.const 1
          i32.add
          local.set $l3
          br $L2
        end
      end
      local.get $p0
      return
    end
    call $runtime.nilPanic
    unreachable)
  (func $runtime.runtimePanic (type $t3)
    i32.const 1072
    i32.const 22
    call $runtime.printstring
    i32.const 1040
    i32.const 23
    call $runtime.printstring
    call $runtime.printnl
    unreachable
    unreachable)
  (func $runtime.putchar (type $t0) (param $p0 i32)
    (local $l1 i32)
    global.get $g0
    i32.const 16
    i32.sub
    local.tee $l1
    global.set $g0
    local.get $l1
    i32.const 0
    i32.store offset=12
    local.get $l1
    local.get $p0
    i32.store8 offset=12
    i32.const 0
    i32.load offset=1024
    local.get $l1
    i32.const 12
    i32.add
    i32.const 1
    call $resource_write
    drop
    local.get $l1
    i32.const 16
    i32.add
    global.set $g0)
  (func $resume (type $t3)
    call $runtime.getFuncPtr
    unreachable)
  (table $T0 1 1 funcref)
  (memory $memory 2)
  (global $g0 (mut i32) (i32.const 66656))
  (global $__heap_base i32 (i32.const 66656))
  (global $__data_end i32 (i32.const 1106))
  (global $__dso_handle i32 (i32.const 1024))
  (export "memory" (memory 0))
  (export "__wasm_call_ctors" (func $__wasm_call_ctors))
  (export "__heap_base" (global 1))
  (export "__data_end" (global 2))
  (export "__dso_handle" (global 3))
  (export "_start" (func $_start))
  (export "runtime.activateTask" (func $runtime.activateTask))
  (export "cwa_main" (func $cwa_main))
  (export "memset" (func $memset))
  (export "resume" (func $resume))
  (data $d0 (i32.const 1024) "\00\00\00\00\00\00\00\00\00\00\00\00")
  (data $d1 (i32.const 1040) "nil pointer dereference\00\00\00\00\00\00\00\00\00panic: runtime error: Hello world!"))
