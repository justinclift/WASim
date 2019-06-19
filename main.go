package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/jackc/pgx"
)

type PrologueInfo struct {
	File int
	Line int
	Col  int
}

var (
	// Just while developing, to allow skipping past the DWARF debug info disassembly step
	disassembleDwarf = false

	// Yes, using globals for this is ugly.  But it's also super simple, so suitable for learning. ;)
	vm *exec.VM
)

const (
	FILE_UNKNOWN = 9999999 - iota
	FILE_STDIN
	FILE_STDOUT
	FILE_STDERR
)

func main() {
	// Connect to the database
	var err error
	cfg := pgx.ConnConfig{
		Host:      "/tmp",
		User:      "jc",
		Database:  "wasim",
		TLSConfig: nil,
	}

	pgPoolConfig := pgx.ConnPoolConfig{cfg, 45, nil, 5 * time.Second}
	pg, err := pgx.NewConnPool(pgPoolConfig)
	if err != nil {
		panic(err)
	}

	// Grab the next available execution_run number
	var dbRun int
	dbQuery := `SELECT nextval('execution_runs_seq')`
	err = pg.QueryRow(dbQuery).Scan(&dbRun)
	if err != nil {
		log.Fatalf("retrieving next execution run number failed: %v\n", err)
	}
	log.Printf("opLog execution run: %d\n", dbRun)

	// Load the wasm file containing DWARF debug info
	// TODO: Pass the file to load via command line arguments
	raw, err := ioutil.ReadFile("testdata/hello-world-simplified.wasm")
	if err != nil {
		panic(err)
	}

	// Parse the wasm file
	m, err := wasm.ReadModule(bytes.NewReader(raw), funcResolver)
	if err != nil {
		panic(err)
	}

	// NOTE: Much of this was initially copied from Delve, then modified

	// Construct a DWARF object from the section data
	if disassembleDwarf {
		err = parseDwarf(m)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Construct the wasm VM
	vm, err = exec.NewVM(m, exec.PGConnPool(pg), exec.PGDBRun(dbRun))
	if err != nil {
		log.Fatalf("could not create wasm vm: %v", err)
	}

	// Locate the main function to run
	foundMain := false
	mainID := uint32(0)
	for name, entry := range m.Export.Entries {
		if (name == "main" && entry.FieldStr == "main") || (name == "cwa_main" && entry.FieldStr == "cwa_main") {
			mainID = entry.Index
			foundMain = true
			break
		}
	}
	if !foundMain {
		panic("no main function found")
	}

	// Run the main function
	_, err = vm.ExecCode(int64(mainID))
	if err != nil {
		panic(err)
	}
}

// Returns the DWARF data contained in a given custom section
func extractDwarf(name string, data []byte) []byte {
	// Skip past the section name string at the start
	r := bytes.NewReader(data)
	var err error
	b := make([]byte, len(name)+1)
	if _, err = io.ReadFull(r, b); err != nil {
		panic(err)
	}

	// The remaining data should be the DWARF info
	var z bytes.Buffer
	if _, err = io.Copy(&z, r); err != nil {
		panic(err)
	}
	return z.Bytes()
}
