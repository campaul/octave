package main

import (
	"fmt"
	"os"

	"github.com/campaul/octave/assemble"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf("Usage: %v [assembly_file] [output_file]\n", args[0])
		return
	}
	file, err := os.Open(args[1])
	if err != nil {
		fmt.Println("Unable to open file for reading: ", err)
		return
	}
	out, err := assemble.Assemble(file)
	if err != nil {
		fmt.Println("Error while assembling:", err)
		return
	}
	outfile, err := os.Create(args[2])
	if err != nil {
		fmt.Println("Unable to open file for writing: ", err)
		return
	}
	_, err = outfile.Write(out)
	if err != nil {
		fmt.Println("Error while writing to file")
		return
	}
}
