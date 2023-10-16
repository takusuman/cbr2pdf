// Trying to do a cb7/cba/cbr/cbt/cbz extractor
package main

import "fmt"
import "github.com/Projeto-Pindorama/cbr2pdf/lib"
import "github.com/liamg/magic"
import "os"
import "path/filepath"
import "golang.org/x/term"

var program string
var outfile *os.File
var formatAssociation = map[string]string{
	".cb7": "7z",
	".cba": "ace",
	".cbr": "rar",
	".cbt": "tar",
	".cbz": "zip",
}

func main() {
	program = os.Args[0]

	if len(os.Args) <= 1 {
		usage()
	}

	if term.IsTerminal(int(os.Stderr.Fd())) {
		outfile = os.Stderr
	} else {
		// We will be writing the log to a file in the same directory as the
		// .cbr itself. This may be disabled and/or sent to another file when
		// we got this program SaaS'ed.
		cbfile_directory := filepath.Dir(os.Args[1])
		outfile, _ = os.Create(fmt.Sprintf("%s/log.txt", cbfile_directory))
	}

	Unpack_cbfile(os.Args[1])
}

func Unpack_cbfile(cbfile string) {
    // For our (un)compressor boilerplate later on, also declare the variables
    // that will contain file information.
    var c func(string, string) error
    var file_extension_per_Ext, file_extension, file_description string
    
    // Gather file information.
	file_extension_per_Ext, file_extension, file_description,
        err := get_cbfile_info(cbfile)
    if err != nil {
        fmt.Fprintf(outfile, 
            "ERROR: could not gather file information.\nERROR: %s\n", err)
    }
    
	// If the file extension doesn't match the expected compression algorithm,
	// just output a warning message.
	if file_extension != formatAssociation[file_extension_per_Ext] {
		fmt.Fprintf(outfile,
			"WARNING: %s extension seems to be out of place.\n", cbfile)
		fmt.Fprintf(outfile,
			"WARNING: The '%s' extension is usually found in %s-compressed files, not in %s.\n",
			file_extension_per_Ext, formatAssociation[file_extension_per_Ext],
			file_extension)
	}
	fmt.Fprintf(outfile, "INFO: %s: %s\n", filepath.Base(cbfile), file_description)
    
    tempdir, _ := os.MkdirTemp("", fmt.Sprintf("%s_teste", filepath.Base(program)))
   	fmt.Fprintf(outfile, "INFO: Created temporary directory %s\n", tempdir)
    fmt.Fprintf(outfile, "INFO: Uncompressing %s\nGentlemen, start your engines!\n", cbfile)
    
    c = get_uncompressor(cbfile)
    
    c(cbfile, tempdir)
}

func get_cbfile_info(file string) (string, string, string, error) {
	cbfile, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(outfile, "ERROR: os.ReadFile failed", err)
	}

	// magic.Lookup(file []byte) will return a FileType struct, from
	// which we'll be only using Types.Extension and Types.Description.
	// We'll be returing the error from it, not from os.ReadFile()
	info, errUnknownMagic := magic.Lookup(cbfile)

	// Also get the file extension by filepath.Ext() to have a reference of
	// what was used as the extension, so we can log if one have used the
	// wrong extension for the compression type.
	file_extension_by_regex := filepath.Ext(file)

	return file_extension_by_regex, info.Extension,
		info.Description, errUnknownMagic
}

func get_uncompressor(archive string) func(string, string) error {
	switch _, filetype, _, _ := get_cbfile_info(archive); filetype {
		case "7z":
		case "ace":
		case "tar":
			fmt.Fprintf(outfile, "ERROR: %s compression not implemented yet.\n",
				filetype)
		case "rar":
			c := func(rarfile, destination string) error {
				return extract.Unrar('e', rarfile, destination)
			}
			return c
		case "zip":
			// Just a boilerplate to the lib.Unzip() function, so we doesn't
			// need to have the other decompressors' functions with the same
			// input.
			c := func(zipfile, destination string) error {
				return extract.Unzip('j', zipfile, destination)
		    	}
			return c
    		default:
			fmt.Fprintf(outfile, "ERROR: %s unsupported format.\n", filetype)
	}
    
    return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s: file.cb[7|r|t|z]\n", program)
	os.Exit(1)
}
