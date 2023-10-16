/* zip - unzip (and possible zip later on) function for cbr2pdf in Google Go
 *
 * Copyright (c) 2023 Luiz AntÃ´nio Rangel
 *
 * SPDX-Licence-Identifier: BSD 3-Clause */

package extract

import "log"
import "os"
import "io"
import "path/filepath"
import "archive/zip"

func Unzip(mode rune, zipfile, destination string) error {
	var jflag, xflag bool
	var written_bytes int64
	var xpath string
	var err error
	var rfile *zip.ReadCloser
	var file *zip.File
	var files []*zip.File
	var current_file io.ReadCloser
	var destfile *os.File

	// "jflag": Doesn't extract files following the Zip file hierarchy.
	// "xflag": Thoroughly follow Zip file hierarchy, creating directory per
	// directory.
	switch mode {
	case 'j':
		jflag = true
	case 'x':
		xflag = true
	default:
		log.Fatalf("ERROR: Unzip(): mode was not specified.\n")
	}

	rfile, err = zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer rfile.Close()

	// Create a directory that only the user can read,
	// write and execute (access) and that his/her group
	// (and the other users in general) can only read
	// and execute.
	// No need to catch errors, I think, since os.MkdirAll won't be returning
	// errors so motiveless.
	_ = os.MkdirAll(destination, 0755)

	// Extract each file inside the Zip file.
	// Assign rfile.File struct array to "files", so we can iter over it.
	files = rfile.File
	for f := 0; f < len(files); f++ {
		file = files[f]
		current_file, err = file.Open()
		if err != nil {
			return err
		}
		defer current_file.Close()

		// Path for extracting, it will be the destination path plus
		// the file name (or, in the case of jflag, the base of the
		// file name).
		if xflag {
			xpath = filepath.Join(destination, file.Name)
		} else if jflag {
			xpath = filepath.Join(destination, filepath.Base(file.Name))
		}

		// Create a directory for the file to come.
		// This applies both to jflag and xflag in pratice.
		if file.FileInfo().IsDir() {
			continue
		} else {
			os.MkdirAll(filepath.Dir(xpath), 0755)
		}

		destfile, err = os.Create(xpath)
		if err != nil {
			return err
		}
		defer destfile.Close()

		// Now, in fact, extract the files, copying from current_file, which was
		// declared as the contents of file.Open() before, to destfile, which
		// is the new file created by os.Create().
		// As always, catch errors with err != nil.
		written_bytes, err = io.Copy(destfile, current_file)
		if err != nil {
			return err
		}
		log.Printf("INFO: %v bytes copied from %s to %s...\n",
			written_bytes, file.Name, xpath)
	}

	return nil
}
