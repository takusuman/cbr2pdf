/* unzip - unzip function for cbr2pdf in Google Go
 *
 * Copyright (c) 2023 Luiz Antônio Rangel
 * 
 * SPDX-Licence-Identifier: BSD 3-Clause */

package lib

//import "errors"
import "log"
import "os"
import "io"
import "path/filepath"
import "archive/zip"

func Unzip(mode rune, zipfile, destination string) (error) {
    var jflag, xflag bool

    switch mode {
	case 'j':
		jflag = true
	case 'x':
		xflag = true
	default:
    	log.Fatalf("ERROR: Unzip(): mode was not specified.\n")
	}
    
	// jflag -> doesn't create directory structure
	// xflag -> extract it all    
    if xflag {
        if err := extract_directories(zipfile, destination); err != nil {
            return err
        }
    } else if jflag {
        if err := extract_nodir(zipfile, destination); err != nil {
            return err
        }
    }
    
    return nil   
}

func extract_directories(zipfile, destination string) (error) {
	var xpath string
    rfile, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer rfile.Close()

	// Create a directory that only the user can read,
	// write and execute (access) and that his/her group
	// (and the other users in general) can only read
	// and execute.
	// No need to catch errors, I think.
	_ = os.MkdirAll(destination, 0755)
    
    // Extract each file inside the Zip file.
	for _, file := range rfile.File {
		current_file, err := file.Open()
        if err != nil {
            return err
        }
        defer current_file.Close()
        
		// Path for extracting, it will be the destination path plus
        // the file name.
		xpath = filepath.Join(destination, file.Name)
		
        if file.FileInfo().IsDir() {
                continue
        } else {
            os.MkdirAll(filepath.Dir(xpath), file.Mode())
        }
        
        destfile, err := os.Create(xpath)
        if err != nil {
            return err
        }
        defer destfile.Close()

        // Now, in fact, extract the files, copying from current_file, which was
        // declared as the contents of file.Open() before, to destfile, which
        // is the new file created by os.Create(). 
        // As always, catch errors with err != nil.
        written_bytes, err := io.Copy(destfile, current_file)
        if err != nil {
            return err
        }
        log.Printf("INFO: %v bytes copied from %s to %v...\n",
                written_bytes, file.Name, xpath)
    }
    return nil
}

func extract_nodir(zipfile, destination string) (error) {
	var xpath string
    rfile, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer rfile.Close()

    // Extract each file inside the Zip file.
	for _, file := range rfile.File {
        if !file.FileInfo().IsDir() {
            current_file, err := file.Open()
            if err != nil {
                return err
            }
            defer current_file.Close()
        
            xpath = filepath.Join(destination, filepath.Base(file.Name))
            
            destfile, err := os.Create(xpath)
            if err != nil {
                return err
            }
            defer destfile.Close()

            written_bytes, err := io.Copy(destfile, current_file)
            if err != nil {
                return err
            }
            log.Printf("INFO: %v bytes copied from %s to %v...\n",
                    written_bytes, file.Name, xpath)
        }
    }
    return nil
}
