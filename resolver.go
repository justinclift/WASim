package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

// * Originally from Xe's land, then modified *
func funcResolver(name string) (*wasm.Module, error) {
	m := wasm.NewModule()
	switch name {
	case "env":
		m.Types = &wasm.SectionTypes{
			Entries: []wasm.FunctionSig{
				{
					Form:        0,
					ParamTypes:  []wasm.ValueType{},
					ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
				},
				{
					Form:        1,
					ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
					ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
				},
				{
					Form:        2,
					ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
					ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
				},
				{
					Form:        3,
					ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
					ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
				},
				{
					Form:        4,
					ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
					ReturnTypes: []wasm.ValueType{},
				},
				{
					Form:        5,
					ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
					ReturnTypes: []wasm.ValueType{},
				},
				{
					Form:        6,
					ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
					ReturnTypes: []wasm.ValueType{},
				},
				{
					Form:        7,
					ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
					ReturnTypes: []wasm.ValueType{},
				},
				{
					Form:        8,
					ParamTypes:  []wasm.ValueType{},
					ReturnTypes: []wasm.ValueType{wasm.ValueTypeF64},
				},
			},
		}
		m.FunctionIndexSpace = []wasm.Function{
			{
				Sig:  &m.Types.Entries[0],
				Host: reflect.ValueOf(ioGetStderr),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[0],
				Host: reflect.ValueOf(ioGetStdout),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[2],
				Host: reflect.ValueOf(resourceOpen),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[3],
				Host: reflect.ValueOf(resourceRead),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[3],
				Host: reflect.ValueOf(resourceWrite),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[7],
				Host: reflect.ValueOf(syscallJSValueCall),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[6],
				Host: reflect.ValueOf(syscallJSValueGet),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[4],
				Host: reflect.ValueOf(syscallJSValuePrepareString),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[6],
				Host: reflect.ValueOf(syscallJSValueLoadString),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[5],
				Host: reflect.ValueOf(syscallJSStringVal),
				Body: &wasm.FunctionBody{},
			},
			{
				Sig:  &m.Types.Entries[8],
				Host: reflect.ValueOf(runtimeTicks),
				Body: &wasm.FunctionBody{},
			},
		}
		m.Export = &wasm.SectionExports{
			Entries: map[string]wasm.ExportEntry{
				"io_get_stderr": {
					FieldStr: "io_get_stdout",
					Kind:     wasm.ExternalFunction,
					Index:    0,
				},
				"io_get_stdout": {
					FieldStr: "io_get_stdout",
					Kind:     wasm.ExternalFunction,
					Index:    1,
				},
				"resource_open": {
					FieldStr: "resource_open",
					Kind:     wasm.ExternalFunction,
					Index:    2,
				},
				"resource_read": {
					FieldStr: "resource_read",
					Kind:     wasm.ExternalFunction,
					Index:    3,
				},
				"resource_write": {
					FieldStr: "resource_write",
					Kind:     wasm.ExternalFunction,
					Index:    4,
				},
				"syscall/js.valueCall": {
					FieldStr: "syscall/js.valueCall",
					Kind:     wasm.ExternalFunction,
					Index:    5,
				},
				"syscall/js.valueGet": {
					FieldStr: "syscall/js.valueGet",
					Kind:     wasm.ExternalFunction,
					Index:    6,
				},
				"syscall/js.valuePrepareString": {
					FieldStr: "syscall/js.valuePrepareString",
					Kind:     wasm.ExternalFunction,
					Index:    7,
				},
				"syscall/js.valueLoadString": {
					FieldStr: "syscall/js.valueLoadString",
					Kind:     wasm.ExternalFunction,
					Index:    8,
				},
				"syscall/js.stringVal": {
					FieldStr: "syscall/js.stringVal",
					Kind:     wasm.ExternalFunction,
					Index:    9,
				},
				"runtime.ticks": {
					FieldStr: "runtime.ticks",
					Kind:     wasm.ExternalFunction,
					Index:    10,
				},
			},
		}
		return m, nil

	case "imports": // For debugging wagon custom_section.wasm test data
		m.Types = &wasm.SectionTypes{
			Entries: []wasm.FunctionSig{
				{
					Form:        0,
					ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
					ReturnTypes: []wasm.ValueType{},
				},
			},
		}
		m.FunctionIndexSpace = []wasm.Function{
			{
				Sig:  &m.Types.Entries[0],
				Host: reflect.ValueOf(wagonImportStub),
				Body: &wasm.FunctionBody{},
			},
		}
		m.Export = &wasm.SectionExports{
			Entries: map[string]wasm.ExportEntry{
				"imported_func": {
					FieldStr: "imported_func",
					Kind:     wasm.ExternalFunction,
					Index:    0,
				},
			},
		}
		return m, nil

	default:
		// To keep things simple for now, only allow the above functions
		return nil, fmt.Errorf("unknown function requested")
	}
}

// * JS host function calls *
func syscallJSValueCall(proc *exec.Process, a int32, b int32, c int32, d int32, e int32, f int32, g int32, h int32, i int32) {
	fmt.Print("In syscallJSValueCall()")
	return
}

func syscallJSValueGet(proc *exec.Process, returnPtr int32, paramB int32, propertyNamePtr int32, propertyNameLen int32, paramE int32, paramF int32) {
	propertyName := make([]byte, propertyNameLen)
	_, err := proc.ReadAt(propertyName, int64(propertyNamePtr))
	if err != nil {
		log.Print(err)
	}

	// Write JS Object ID to memory at the return pointer location
	var endianess = binary.LittleEndian
	var val uint64
	switch string(propertyName) { // Known good object IDs captured from Go wasm running in FF
	case "Object":
		val = 0x7FF8000300000008
	case "Array":
		val = 0x7FF8000300000009
	case "Int8Array":
		val = 0x7FF800030000000A
	case "Int16Array":
		val = 0x7FF800030000000B
	case "Int32Array":
		val = 0x7FF800030000000C
	case "Uint8Array":
		val = 0x7FF800030000000D
	case "Uint16Array":
		val = 0x7FF800030000000E
	case "Uint32Array":
		val = 0x7FF800030000000F
	case "Float32Array":
		val = 0x7FF8000300000010
	case "Float64Array":
		val = 0x7FF8000300000011
	case "document":
		val = 0x7FF8000300000012
	}
	buf := make([]byte, 8)
	endianess.PutUint64(buf, val)
	_, err = proc.WriteAt(buf, int64(returnPtr))
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("Returned JS object ID %#x for js.Global().Get(\"%s\")", val, propertyName)
	return
}

func syscallJSValuePrepareString(proc *exec.Process, a int32, b int32, c int32, d int32, e int32, f int32) {
	fmt.Print("In syscallJSValuePrepareString()")
	return
}

func syscallJSValueLoadString(proc *exec.Process, a int32, b int32, c int32, d int32, e int32, f int32) {
	fmt.Print("In syscallJSValueLoadString()")
	return
}

func syscallJSStringVal(proc *exec.Process, a int32, b int32, c int32, d int32, e int32, f int32) {
	fmt.Print("In syscallJSStringVal()")
	return
}

// * Other host function calls *
func wagonImportStub(proc *exec.Process, x int32) {
	return
}

func runtimeTicks(proc *exec.Process) int32 {
	return FILE_STDOUT
}

func ioGetStderr(proc *exec.Process) int32 {
	return FILE_STDERR
}

func ioGetStdout(proc *exec.Process) int32 {
	return FILE_STDOUT
}

func resourceOpen(proc *exec.Process, urlPtr uint32, urlLen uint32) int32 {

	// Read a section of the WASM vm's memory
	data := make([]byte, urlLen)
	bytesRead, err := proc.ReadAt(data, int64(urlPtr))
	if err != nil {
		log.Print(err)
		return int32(bytesRead)
	}
	u := string(data)

	_, err = url.Parse(u)
	// uu, err := url.Parse(u)
	if err != nil {
		log.Printf("can't parse url %s: %v", u, err)
		return 0
	}

	// q := uu.Query()
	// switch uu.Scheme {
	// case "log":
	// 	prefix := q.Get("prefix")
	// 	file = fileresolver.Log(os.Stdout, p.name+": "+prefix, log.LstdFlags)
	// case "random":
	// 	file = fileresolver.Random()
	// case "null":
	// 	file = fileresolver.Null()
	// case "zero":
	// 	file = fileresolver.Zero()
	// case "http", "https":
	// 	var err error
	// 	file, err = fileresolver.HTTP(p.hc, uu)
	// 	if err != nil {
	// 		p.logger.Printf("can't resource_open(%q): %v", u, err)
	// 		return 0, UnknownError
	// 	}
	// default:
	// 	return 0, fmt.Errorf("unknown url: %s", u)
	// }

	// fid := rand.Int31()
	// FileHandles[fid] = file

	// Return a file handle
	return FILE_UNKNOWN
}

// Host function call "resource_read"
// Just a stub for now
func resourceRead(proc *exec.Process, fid int32, dataPtr int32, dataLen int32) int32 {

	// TODO: This function seems like it should be reading bytes from the given file (eg os.Stdin), then writing
	//       them to the given spot in the VM's memory, up to dataLen in length

	// data := make([]byte, dataLen)
	// bytesRead, err := proc.ReadAt(data, int64(dataPtr))  // TODO: proc.Write() is probably the call to use here
	// if err != nil {
	// 	log.Print(err)
	// 	return 1
	// }
	//
	// if bytesRead != int(dataLen) {
	// 	log.Printf("Incorrect # of bytes read.  Requested %d, but read %d\n", dataLen, bytesRead)
	// 	return 2
	// }

	// fmt.Printf("%s", string(data))
	return int32(0)
}

// Host function call "resource_write"
func resourceWrite(proc *exec.Process, fid int32, dataPtr int32, dataLen int32) int32 {

	// Determine the output file to write to
	var outTarget *os.File
	switch fid {
	case FILE_STDERR:
		outTarget = os.Stderr
	case FILE_STDOUT:
		outTarget = os.Stdout
	}

	// Read the data from the VM's memory
	data := make([]byte, dataLen)
	bytesRead, err := proc.ReadAt(data, int64(dataPtr))
	if err != nil {
		log.Print(err)
		return 1 // TODO: Find out if there are meaningful error codes defined in the spec that should be returned
	}

	if bytesRead != int(dataLen) {
		log.Printf("Incorrect # of bytes read.  Requested %d, but read %d\n", dataLen, bytesRead)
		return 2
	}

	// Write the data to the requested output
	_, err = fmt.Fprintf(outTarget, "%s", string(data))
	if err != nil {
		log.Print(err)
	}
	return 0
}
