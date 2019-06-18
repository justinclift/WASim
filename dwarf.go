package main

import (
	"bytes"
	"debug/dwarf"
	"fmt"
	"io"

	dwfRdr "github.com/go-delve/delve/pkg/dwarf/reader"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/go-interpreter/wagon/wasm/leb128"
)

// TODO: This should probably be merged upstream into wagons' ReadModule() and/or readSection()

// Parse a wasm file, turning the DWARF data into something useful
func parseDwarf(m *wasm.Module) error {
	// Check for the custom sections generated by LLVM
	functionNames := make(map[int]string)
	data := make(map[string][]byte)
	for _, sec := range m.Customs {
		switch sec.Name {
		case "name":
			extractFunctionNames(sec, functionNames)

		case ".debug_info":
			data["info"] = extractDwarf(sec.Name, sec.RawSection.Bytes)

		case ".debug_macinfo":
			data["macinfo"] = extractDwarf(sec.Name, sec.RawSection.Bytes)

		case ".debug_pubtypes":
			data["pubtypes"] = extractDwarf(sec.Name, sec.RawSection.Bytes)

		case ".debug_ranges":
			data["ranges"] = extractDwarf(sec.Name, sec.RawSection.Bytes)

		case ".debug_abbrev":
			data["abbrev"] = extractDwarf(sec.Name, sec.RawSection.Bytes)

		case ".debug_line":
			data["line"] = extractDwarf(sec.Name, sec.RawSection.Bytes)

		case ".debug_str":
			data["str"] = extractDwarf(sec.Name, sec.RawSection.Bytes)

		case ".debug_pubnames":
			data["pubnames"] = extractDwarf(sec.Name, sec.RawSection.Bytes)
		}
	}

	// If there wasn't any dwarf data, there's nothing to do
	if len(data) == 0 {
		return nil
	}

	// Construct a DWARF object from the section data
	d, err := dwarf.New(data["abbrev"], nil, nil, data["info"], data["line"], data["pubnames"], data["ranges"], data["str"])
	if err != nil {
		panic(err)
	}
	reader := dwfRdr.New(d)

	var compiler, fileName string
	var cuRanges [][2]uint64
	var proL PrologueInfo
	srcFiles := make(map[string]int)

	// Read through the DWARF data
outer:
	for tag, err := reader.Next(); tag != nil; tag, err = reader.Next() {
		if err != nil {
			panic(err)
		}

		// Info on DWARF 4 fields:
		//   https://github.com/golang/go/blob/d97bd5d07ac4e7b342053b335428ff9c97212f9f/src/debug/dwarf/entry.go#L217-L244

		switch tag.Tag {
		case dwarf.TagCompileUnit:

			// Store the compiler name (eg "TinyGo") and name of the main compiled .go file
			compiler = tag.Val(dwarf.AttrProducer).(string)
			fileName = tag.Val(dwarf.AttrName).(string)

			// Grab the Program Counter ranges in this compile unit
			cuRanges, err = d.Ranges(tag)
			if err != nil {
				panic(err)
			}

			fmt.Print("Compile unit")

			if fileName != "" {
				srcFiles[fileName] = 0
				fmt.Printf("  * Original filename: %s\n", fileName)
			}
			if compiler != "" {
				fmt.Printf("  * Compiled using: %s\n", compiler)
			}
			compileDir, ok := tag.Val(dwarf.AttrCompDir).(string)
			if ok {
				fmt.Printf("  * Compile dir: %s\n", compileDir)
			}
			lineInfoOffset, ok := tag.Val(dwarf.AttrStmtList).(int64)
			if ok {
				fmt.Printf("  * Line info offset: %v\n", lineInfoOffset)
			}

			if len(cuRanges) > 0 {
				fmt.Printf("  * Number of ranges: %d\n", len(cuRanges))
				for _, j := range cuRanges {
					fmt.Printf("    * Low: %v  High: %v\n", j[0], j[1])
				}
			}

			for i, j := range tag.Field {
				fmt.Printf("    * Attribute %d - Attr: '%v'  Class: '%v'  Value: '%v'\n", i, j.Attr, j.Class, j.Val)
			}
			lines, err := d.LineReader(tag)
			if err != nil {
				panic(err)
			}
			var lEntry dwarf.LineEntry
			i := 0
			for {
				err = lines.Next(&lEntry)
				if err != nil {
					break
				}
				fmt.Printf("Entry %d\n", i)
				fmt.Printf("  * File: %v\n", lEntry.File.Name)
				fmt.Printf("  * Address: %v (0x%x)\n", lEntry.Address, lEntry.Address)
				fmt.Printf("  * Basic Block: %v\n", lEntry.BasicBlock)
				fmt.Printf("  * Is Statement: %v\n", lEntry.IsStmt)
				fmt.Printf("  * Line: %v\n", lEntry.Line)
				fmt.Printf("  * Column: %v\n", lEntry.Column)
				fmt.Printf("  * Prologue End: %v\n", lEntry.PrologueEnd)
				fmt.Printf("  * Epilogue Begin: %v\n", lEntry.EpilogueBegin)
				fmt.Printf("  * End Sequence: %v\n", lEntry.EndSequence)
				i++

				// Add the file name to the list of source files
				var src int
				var ok bool
				if src, ok = srcFiles[lEntry.File.Name]; !ok {
					n := len(srcFiles) + 1
					srcFiles[lEntry.File.Name] = n
					src = n
				}

				// Record where the prologue ends
				proL.File = src
				proL.Line = lEntry.Line
				proL.Col = lEntry.Column
			}

		case dwarf.TagVariable:
			fmt.Println("Variable present")

		case dwarf.TagSubprogram:

			fmt.Println("In sub program")

			// Grab the Program Counter ranges in this compile unit
			cuRanges, err = d.Ranges(tag)
			if err != nil {
				panic(err)
			}

			if len(cuRanges) > 0 {
				fmt.Printf("Number of ranges: %d\n", len(cuRanges))
				for _, j := range cuRanges {
					fmt.Printf("  * Low: %v  High: %v\n", j[0], j[1])
				}
			}

			for i, j := range tag.Field {
				fmt.Printf("    * Attribute %d - Attr: '%v'  Class: '%v'  Value: '%v'\n", i, j.Attr, j.Class, j.Val)
			}

			if tag.Children {
				abstractOriginNameTable := make(map[dwarf.Offset]string)
				fmt.Print("  * Has children")
				ok1 := false
				// TODO: Figure out which bits to keep, remove, etc.
				for {

					tag, err = reader.Next()
					if err != nil {
						break outer
					}
					if tag.Tag == 0 {
						break
					}
					if tag.Tag == dwarf.TagInlinedSubroutine {
						originOffset := tag.Val(dwarf.AttrAbstractOrigin).(dwarf.Offset)
						name := abstractOriginNameTable[originOffset]
						fmt.Printf("Name: %v\n", name)
						if ranges, _ := d.Ranges(tag); len(ranges) == 1 {
							ok1 = true
							// lowpc = ranges[0][0]
							// highpc = ranges[0][1]
						}
						callfileidx, ok1 := tag.Val(dwarf.AttrCallFile).(int64)
						callline, ok2 := tag.Val(dwarf.AttrCallLine).(int64)
						if ok1 && ok2 {
							fmt.Printf("callfileidx: %v\n", callfileidx)
							fmt.Printf("callline: %v\n", callline)
						}
					}
					reader.SkipChildren()
				}
				if ok1 {
					fmt.Print("Ok1 is true")
				}
			}

		case dwarf.TagBaseType:

			if name, ok := tag.Val(dwarf.AttrName).(string); ok {
				fmt.Printf("Base type name: %v\n", name)
			}
			reader.SkipChildren()

		default:
			// Display the human readable name for the DWARF tag
			fmt.Println(tag.Tag.String())

			// Display each of the attributes for the DWARF tags
			for i, j := range tag.Field {
				fmt.Printf("  * Attribute %d - Attr: '%v'  Class: '%v'  Value: '%v'\n", i, j.Attr, j.Class, j.Val)
			}
		}

		fmt.Println()
	}
	return nil
}

// Parses the "name" custom section, extracting the function names present in the file
func extractFunctionNames(sec *wasm.SectionCustom, functionNames map[int]string) {
	r := bytes.NewReader(sec.RawSection.Bytes)
	for {
		// LLVM generated "name" section starts with 4, being the length of the word "name"
		m, err := leb128.ReadVarUint32(r)
		if err != nil || m != 4 {
			break
		}
		b := make([]byte, 5) // Length of "name" plus the byte before hand giving its length
		if _, err = io.ReadFull(r, b); err != nil {
			panic(err)
		}

		// Length of the remaining data in this section
		payloadLen, err := leb128.ReadVarUint32(r)
		if err != nil {
			panic(err)
		}
		data := make([]byte, int(payloadLen))
		n, err := r.Read(data)
		if err != nil {
			panic(err)
		}
		if n != len(data) {
			panic("len mismatch")
		}
		r := bytes.NewReader(data)
		for {
			// The first value contains the number of functions
			count, err := leb128.ReadVarUint32(r)
			if err != nil {
				break
			}
			for i := 0; i < int(count); i++ {
				// Each function name entry contains:
				//   * Its function number (eg 0, 1, 2, etc)
				//   * The length of the name in bytes. eg 10 for "wasm_stuff"
				//   * The name of the function
				index, err := leb128.ReadVarUint32(r)
				if err != nil {
					panic(err)
				}
				nameLen, err := leb128.ReadVarUint32(r)
				if err != nil {
					panic(err)
				}
				name := make([]byte, int(nameLen))
				n, err := r.Read(name)
				if err != nil {
					panic(err)
				}
				if n != len(name) {
					panic("len mismatch")
				}
				functionNames[int(index)] = string(name)
			}
		}
	}
}