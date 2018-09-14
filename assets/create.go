// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type file struct {
	name  string
	data  string
	mime  string
	mtime time.Time
	size  int64 // If 0, it means the data is uncompressed
	hash  []byte
}

type fileSlice []*file

func localName(filename, root string) string {
	if strings.HasPrefix(filename, root) {
		name := filename[len(root):]
		parts := strings.Split(name, string(os.PathSeparator))
		return path.Join(parts...)
	}
	return filename
}

func process(filenames []string, root string) (fileSlice, error) {
	var b bytes.Buffer
	var b2 bytes.Buffer
	hash := sha256.New()

	files := make(fileSlice, 0, len(filenames))
	for _, asset := range filenames {
		f, err := os.Open(asset)
		if err != nil {
			return nil, err
		}
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if _, err := b.ReadFrom(f); err != nil {
			return nil, err
		}
		f.Close()
		compressedWriter, _ := gzip.NewWriterLevel(&b2, gzip.BestCompression)
		writer := io.MultiWriter(compressedWriter, hash)
		if _, err := writer.Write(b.Bytes()); err != nil {
			return nil, err
		}
		compressedWriter.Close()

		one := file{
			name:  localName(asset, root),
			data:  b.String(),
			mime:  mime.TypeByExtension(filepath.Ext(asset)),
			mtime: stat.ModTime(),
			hash:  hash.Sum(nil),
		}
		if b2.Len() < b.Len() {
			one.data = b2.String()
			one.size = stat.Size()
		}
		files = append(files, &one)

		b.Reset()
		b2.Reset()
		hash.Reset()
	}
	return files, nil
}

func generate(from fileSlice, media io.Writer, pack string) {
	fmt.Fprintf(media, "// *** DO NOT EDIT ***\n// This file was generated by github.com/seamia/tools/assets/cmd/assets\n")
	fmt.Fprintf(media, "// on %s\n", time.Now().Format(time.RFC850))
	fmt.Fprintf(media, "package %v\n", pack)
	fmt.Fprintln(media, "")
	fmt.Fprintln(media, "import \"github.com/seamia/tools/assets\"\n")

	fmt.Fprintf(media, "var staticAssets = assets.AssetRoot{\n")
	for _, one := range from {

		fmt.Fprintf(media, "\t%q: {\n", one.name)
		fmt.Fprintf(media, "\t\tData:  %q,\n", one.data)
		if len(one.mime) != 0 {
			fmt.Fprintf(media, "\t\tMime:  %q,\n", one.mime)
		}
		fmt.Fprintf(media, "\t\tMtime: %v,\n", one.mtime.Unix())
		if one.size != 0 {
			fmt.Fprintf(media, "\t\tSize:  %v,\n", one.size)
		}
		fmt.Fprintf(media, "\t\tHash:  %q,\n", hex.EncodeToString(one.hash))
		fmt.Fprintf(media, "\t},\n")
	}
	fmt.Fprintf(media, "}\n")
	fmt.Fprintf(media, "\nfunc init() {\n\tassets.Assign(staticAssets)\n}\n")
}

func Generate(list []string, root string, destination string, pack string) error {

	back, err := process(list, root)
	if err != nil {
		return err
	}
	out, err := os.Create(destination)
	if err == nil {
		defer out.Close()
		generate(back, out, pack)
	}
	return err
}
