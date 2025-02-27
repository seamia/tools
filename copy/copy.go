package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// tmpPermissionForDirectory makes the destination directory writable,
	// so that stuff can be copied recursively even if any original directory is NOT writable.
	// See https://github.com/otiai10/copy/pull/9 for more information.
	tmpPermissionForDirectory = os.FileMode(0755)

	flagCopyFiles        = true
	flagCheckBeforeWrite = true
	flagRemoveOldFiles   = true
	flagVerbal           = false

	cutoffFileSize = 1 * 1024 * 1024

	allowOnlyGoFiles = true
)

var (
	bytesCopied    int64
	bytesDiscarded int64
	bytesNotCopied int64
)

// Copy copies src to dest, doesn't matter if src is a directory or a file
func Copy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		fmt.Println("error10", err)
		return err
	}

	if flagRemoveOldFiles {
		srcPrefix = src
		dstPrefix = dest
	}

	err = copy(src, dest, info)
	if err == nil {
		if info, err := os.Lstat(dest); err == nil {
			removeOldFiles(src, dest, info)
		}
	}
	return err
}

func removeOldFiles(src, dest string, info os.FileInfo) {
	if flagRemoveOldFiles {
		if info.Mode()&os.ModeSymlink != 0 {
			return
		}
		if info.IsDir() {
			// return dcopy(src, dest, info)
		}
		// return fcopy(src, dest, info)
	}
}

// copy dispatches copy-funcs according to the mode.
// Because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
func copy(src, dest string, info os.FileInfo) error {

	if disregard(info) {
		bytesDiscarded += info.Size()
		return nil
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return lcopy(src, dest, info)
	}
	if info.IsDir() {
		return dcopy(src, dest, info)
	}
	return fcopy(src, dest, info)
}

// fcopy is for just a file,
// with considering existence of parent directory
// and file permission.
func fcopy(src, dest string, info os.FileInfo) error {

	if flagCopyFiles {
		if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
			fmt.Println("error1", err)
			return err
		}

		if flagRemoveOldFiles {
			rememberFileExistence(src)
		}

		if flagCheckBeforeWrite {
			if dinfo, err := os.Lstat(dest); err == nil {
				if info.Size() == dinfo.Size() {
					if info.ModTime() == dinfo.ModTime() {
						bytesNotCopied += info.Size()
						return nil
					}
				}
			}
		}

		f, err := os.Create(dest)
		if err != nil {
			fmt.Println("error2", err)
			return err
		}
		defer f.Close()

		if err = os.Chmod(f.Name(), info.Mode()); err != nil {
			fmt.Println("error3", err)
			return err
		}

		s, err := os.Open(src)
		if err != nil {
			fmt.Println("error4", err)
			return err
		}
		defer s.Close()

		_, err = io.Copy(f, s)
		_ = s.Close()

		if err == nil {
			bytesCopied += info.Size()
			if flagCheckBeforeWrite {
				f.Close()
				err := os.Chtimes(dest, info.ModTime(), info.ModTime())

				if err != nil {
					fmt.Println(err)
				}
			}
		}

		return err
	} else {
		return nil
	}
}

// dcopy is for a directory,
// with scanning contents inside the directory
// and pass everything to "copy" recursively.
func dcopy(srcdir, destdir string, info os.FileInfo) error {

	originalMode := info.Mode()

	// Make dest dir with 0755 so that everything writable.
	if err := os.MkdirAll(destdir, tmpPermissionForDirectory); err != nil {
		fmt.Println("error5", err)
		return err
	}
	// Recover dir mode with original one.
	defer os.Chmod(destdir, originalMode)

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		fmt.Println("error6", err)

		return nil // todo: is this an appropriate ??
		// return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())
		if err := copy(cs, cd, content); err != nil {
			// If any error, exit immediately
			fmt.Println("error7", err)
			return err
		}
	}

	return nil
}

// lcopy is for a symlink,
// with just creating a new symlink by replicating src symlink.
func lcopy(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		fmt.Println("error8", err)
		return err
	}

	if flagCheckBeforeWrite {
		points2, err := os.Readlink(dest)
		if err == nil {
			if points2 == src {
				return nil // no need to copy
			}
		}
	}

	return os.Symlink(src, dest)
}

func disregard(info os.FileInfo) bool {
	name := strings.ToLower(info.Name())
	if info.IsDir() {
		for _, one := range []string{"vendor", ".git", ".idea", ".terraform", "node_modules", "tmp", "linux_amd64", "darwin_amd64", "windows_amd64", ".vscode"} {
			if name == one {
				return true
			}
		}

		if strings.HasPrefix(name, ".") {
			if flagVerbal {
				fmt.Println("==============", info.Name())
			}
		}
		// fmt.Println("==============", info.Name())
	} else {

		if allowOnlyGoFiles {
			if strings.HasSuffix(name, ".go") {
				return false
			}

			// exact match
			for _, one := range []string{"go.mod"} {
				if name == one {
					return false
				}
			}

			return true
		}

		// exact match
		for _, one := range []string{".ds_store"} {
			if name == one {
				return true
			}
		}

		// suffix match
		for _, one := range []string{".pack", ".wav", ".gif", ".png", ".jpg", ".tga", ".glb", ".fbx", ".a", ".mpk",
			".psd", ".gz", ".zip", ".bz2", ".jar", ".bmp", ".exe", ".deb", ".obj", ".dll", ".diorama", ".asset", ".svg",
			".dylib", ".tif", ".tiff", ".so", ".pdf", ".bundle", ".bin", ".blend", ".lock"} {
			if strings.HasSuffix(name, one) {
				return true
			}
		}
		// fmt.Println("----------------------------------", info.Name())

		if info.Size() > cutoffFileSize {
			// this is an oversized file
			for _, one := range []string{".sql", ".json", ".js", ".py", ".go"} {
				if strings.HasSuffix(name, one) {
					// this file has a legit extension
					return false
				}
			}

			if flagVerbal {
				fmt.Println("\t\t\tdiscarding ", info.Name(), " due to size: ", info.Size())
			}
			return true
		}
	}
	return false
}

var (
	existing  map[string]bool
	srcPrefix = ""
	dstPrefix = ""
)

func rememberFileExistence(name string) {
	if existing == nil {
		existing = make(map[string]bool)
	}

	// remove prefix (if specified)
	l := len(srcPrefix)
	if l > 0 {
		if strings.HasPrefix(name, srcPrefix) {
			name = name[l:]
		}
	}

	existing[name] = true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "missing required arguments: app from to\n")
		return
	}
	src := os.Args[1]
	dst := os.Args[2]

	err := Copy(src, dst)
	if err != nil {
		fmt.Println("FAILED:", err)
	} else {
		fmt.Println("... copied:", bytesCopied, "; discarded:", bytesDiscarded, "; not-copied: ", bytesNotCopied)
	}
}
