package main

import (
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
			},
		}
		return m, nil

	case "syscall":
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
			},
		}
		m.FunctionIndexSpace = []wasm.Function{
			{
				Sig:  &m.Types.Entries[0],
				Host: reflect.ValueOf(syscallJSStub),
				Body: &wasm.FunctionBody{},
			},
		}
		m.Export = &wasm.SectionExports{
			Entries: map[string]wasm.ExportEntry{
				"js": {
					FieldStr: "js",
					Kind:     wasm.ExternalFunction,
					Index:    0,
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

// * Stub host function calls *
func wagonImportStub(proc *exec.Process, x int32) {
	return
}

func syscallJSStub(proc *exec.Process) int32 {
	return 0
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
