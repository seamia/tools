package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	smallestFileSizeToConsider = 1024
	sizeOfPreviewChunk         = 4096
)

func Hash(raw []byte) string {
	h := sha1.New()
	h.Write(raw)
	return hex.EncodeToString(h.Sum(nil))
}

func hash_file_sha1(filePath string, quick bool) (string, int64, error) {
	var returnSHA1String string
	var fileSize int64 = 0
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String, fileSize, err
	}
	if fi, err := file.Stat(); err == nil {
		fileSize = fi.Size()
	}

	if quick && fileSize > sizeOfPreviewChunk {
		fileSize = sizeOfPreviewChunk
	}

	defer file.Close()
	hash := sha1.New()
	if _, err := io.CopyN(hash, file, fileSize); err != nil {
		return returnSHA1String, fileSize, err
	}
	hashInBytes := hash.Sum(nil)[:20]
	returnSHA1String = hex.EncodeToString(hashInBytes)
	return returnSHA1String, fileSize, nil
}

func nameIsAcceptable(name *string) bool {
	return strings.ToLower(*name)[len(*name)-3:] == ".go"
}

func findDuplicates(root string) {

	var found []string
	var foundSize int64
	bySize := make(map[int64][]string)

	fmt.Println("Discovering all the relevant files from:", root)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if nameIsAcceptable(&path) {

			if !info.IsDir() {
				size := info.Size()
				if size >= smallestFileSizeToConsider {

					/*
						fm := info.Mode()

						if fm & os.ModeDir == os.ModeDir {
							fmt.Println("is a directory")
						}

						ModeAppend                                     // a: append-only
						ModeExclusive                                  // l: exclusive use
						ModeTemporary                                  // T: temporary file; Plan 9 only
						ModeSymlink                                    // L: symbolic link
						ModeDevice                                     // D: device file
						ModeNamedPipe                                  // p: named pipe (FIFO)
						ModeSocket                                     // S: Unix domain socket
						ModeSetuid                                     // u: setuid
						ModeSetgid                                     // g: setgid
						ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
						ModeSticky


						if fm != 438 {
							fmt.Println(fm.String())
						}
					*/

					bySize[size] = append(bySize[size], path)
					found = append(found, path)
					foundSize += size
				}
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Total files:", len(found), "; total size:", foundSize)

	var sizeable []string
	var saved int64
	for size, files := range bySize {
		if len(files) > 1 {
			for _, file := range files {
				sizeable = append(sizeable, file)
			}
		} else {
			saved += size
		}
	}

	fmt.Println("Dups by size files:", len(sizeable), ", saved (bytes):", saved)

	fmt.Println("Calculating prefix hashes of ", len(sizeable), "entries...")
	dups := make(map[string][]string)
	for _, one := range sizeable {
		if hash, _, err := hash_file_sha1(one, true); err == nil {
			dups[hash] = append(dups[hash], one)
		}
	}

	// we got 'quick' (prefix) match -> let's do full match now
	var again []string
	for _, files := range dups {
		if len(files) > 1 {
			again = append(again, files...)
		}
	}
	fmt.Println("there were", len(again), "quick match entries for", len(dups), "unique hashes")

	fmt.Println("Calculating full hashes of ", len(again), "entries...")
	dups = make(map[string][]string)
	for _, one := range again {
		if hash, _, err := hash_file_sha1(one, false); err == nil {
			dups[hash] = append(dups[hash], one)
		}
	}

	for _, files := range dups {
		if len(files) > 1 {
			shortest := len(files[0])
			shortestName := files[0]
			for _, file := range files {
				if len(file) < shortest {
					shortest = len(file)
					shortestName = file
				}
			}

			fmt.Println("rem Keep this one:", shortestName)
			for _, file := range files {
				if file != shortestName {
					fmt.Println("call remove_dupe(", file, ")")
				}
			}
			fmt.Println("")
		}
	}
}

func renameFile(path, hash string) {
	ext := filepath.Ext(path)
	dir := filepath.Dir(path)

	newpath := filepath.Join(dir, strings.ToLower(hash+ext))

	fmt.Println(hash, "   ", path)
	if path != newpath {
		if err := os.Rename(path, newpath); err != nil {
			fmt.Println("*** failed to rename [", path, "] to [", newpath, "] due to:", err)
		}
	}
}

func renameBasedOnContent(root string) {

	var found []string

	fmt.Println("Discovering all the relevant files from:", root)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			if strings.Index(path, "\\xxx\\") < 0 {
				found = append(found, path)
			}

			/*
				if dir := filepath.Dir(path); dir == root {
					found = append(found, path)
				}
			*/

			/*
				if len(filepath.Base(path)) - len(filepath.Ext(path)) != 40 {
					if strings.Index(path, "\\xxx\\xxx\\") < 0 {
						found = append(found, path)
					}
				}
			*/
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Total files:", len(found))

	found = rearrange(found)
	for _, one := range found {
		if hash, _, err := hash_file_sha1(one, false); err == nil {
			renameFile(one, hash)
		}
	}
}

func fileNameLength(full string) int {
	left := strings.LastIndex(full, "\\") // todo: this will break under !windows
	if left == -1 {
		return len(full)
	}
	name := full[left+1:]
	if right := strings.Index(name, "."); right > 0 {
		return right
	}
	return len(name)
}

func rearrange(found []string) []string {
	sort.Slice(found, func(i, j int) bool {
		return fileNameLength(found[i]) < fileNameLength(found[j])
	})

	return found
}

func main() {

	if len(os.Args) > 1 {
		// root := "Y:\\grpc\\proto\\"
		// root := "O:\\play\\src\\github.com"
		root := os.Args[1]

		// findDuplicates(root)
		renameBasedOnContent(root)
	} else {
		fmt.Println("specify the directory you'd like to dedupe")
	}
}
