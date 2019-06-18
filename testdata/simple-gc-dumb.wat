(module
  (type $t0 (func (result i32)))
  (type $t1 (func (param i32 i32 i32) (result i32)))
  (type $t2 (func))
  (type $t3 (func (param i32)))
  (import "env" "io_get_stdout" (func $io_get_stdout (type $t0)))
  (import "env" "resource_write" (func $resource_write (type $t1)))
  (func $__wasm_call_ctors (type $t2))
  (func $_start (type $t2)
    call $runtime.initAll)
  (func $runtime.initAll (type $t2)
    i32.const 0
    call $io_get_stdout
    i32.store offset=1024)
  (func $cwa_main (type $t2)
    call $runtime.initAll
    call $simple.go.main)
  (func $simple.go.main (type $t2)
    call $simple.go.hello)
  (func $runtime.printnl (type $t2)
    i32.const 13
    call $runtime.putchar
    i32.const 10
    call $runtime.putchar)
  (func $runtime.putchar (type $t3) (param $p0 i32)
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
  (func $runtime.printstring (type $t2)
    i32.const 72
    call $runtime.putchar
    i32.const 101
    call $runtime.putchar
    i32.const 108
    call $runtime.putchar
    i32.const 108
    call $runtime.putchar
    i32.const 111
    call $runtime.putchar
    i32.const 32
    call $runtime.putchar
    i32.const 102
    call $runtime.putchar
    i32.const 114
    call $runtime.putchar
    i32.const 111
    call $runtime.putchar
    i32.const 109
    call $runtime.putchar
    i32.const 32
    call $runtime.putchar
    i32.const 97
    call $runtime.putchar
    i32.const 32
    call $runtime.putchar
    i32.const 102
    call $runtime.putchar
    i32.const 117
    call $runtime.putchar
    i32.const 110
    call $runtime.putchar
    i32.const 99
    call $runtime.putchar
    i32.const 116
    call $runtime.putchar
    i32.const 105
    call $runtime.putchar
    i32.const 111
    call $runtime.putchar
    i32.const 110
    call $runtime.putchar
    i32.const 32
    call $runtime.putchar
    i32.const 99
    call $runtime.putchar
    i32.const 97
    call $runtime.putchar
    i32.const 108
    call $runtime.putchar
    i32.const 108
    call $runtime.putchar)
  (func $simple.go.hello (type $t2)
    call $runtime.printstring
    call $runtime.printnl)
  (table $T0 1 1 funcref)
  (memory $memory 2)
  (global $g0 (mut i32) (i32.const 66576))
  (global $__heap_base i32 (i32.const 66576))
  (global $__data_end i32 (i32.const 1028))
  (global $__dso_handle i32 (i32.const 1024))
  (export "memory" (memory 0))
  (export "__wasm_call_ctors" (func $__wasm_call_ctors))
  (export "__heap_base" (global 1))
  (export "__data_end" (global 2))
  (export "__dso_handle" (global 3))
  (export "_start" (func $_start))
  (export "cwa_main" (func $cwa_main))
  (data $d0 (i32.const 1024) "\00\00\00\00"))
