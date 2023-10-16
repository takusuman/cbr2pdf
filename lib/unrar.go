/* unrar - unrar function for cbr2pdf in Google Go
 *
 * Copyright (c) 2023 Luiz AntÃ´nio Rangel
 *
 * SPDX-Licence-Identifier: BSD 3-Clause */

package extract

import "log"
import "os"
import "io"
import "path/filepath"
import rar "github.com/nwaples/rardecode"

func Unrar(mode rune, rarfile, destination string) error {
	var eflag, xflag bool
	var written_bytes int64
	var xpath string
	var err error
	var rfile *rar.ReadCloser
	var file *rar.FileHeader
	var destfile *os.File

	// "eflag": Doesn't extract files following the Rar file hierarchy.
	// "xflag": Thoroughly follow Rar file hierarchy, creating directory per
	// directory.
	switch mode {
	case 'e':
		eflag = true
	case 'x':
		xflag = true
	default:
		log.Fatalf("ERROR: Unrar(): mode was not specified.\n")
	}

	rfile, err = rar.OpenReader(rarfile, "")
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

	for ;; {
		file, err = rfile.Next()
		if err == io.EOF {
			break
		}

		// Path for extracting, it will be the destination path plus
		// the file name (or, in the case of eflag, the base of the
		// file name).
		if xflag {
			xpath = filepath.Join(destination, file.Name)
		} else if eflag {
			xpath = filepath.Join(destination, filepath.Base(file.Name))
		}

		// Create a directory for the file to come.
		// This applies both to eflag and xflag in pratice.
		if file.IsDir {
			continue
		} else {
			os.MkdirAll(filepath.Dir(xpath), 0755)
		}

		if destfile, err = os.Create(xpath); err != nil {
			return err
		}
		defer destfile.Close()

		if err = destfile.Chmod(file.Mode()); err != nil {
			return err
		}

		// Now, in fact, extract the files, copying from current_file, which was
		// declared as the contents of file.Open() before, to destfile, which
		// is the new file created by os.Create().
		// As always, catch errors with err != nil.
		written_bytes, err = io.Copy(destfile, rfile)
		if err != nil {
			return err
		}
		log.Printf("INFO: %v bytes copied from %s to %s...\n",
			written_bytes, file.Name, xpath)
	}

	return nil
}
