package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"strings"
)

var builtAt time.Time

func FileSystem() assetFileSystemI {
	return assets
}

type assetFileSystemI interface {
	Open(name string) (http.File, error)
	AssetNames() []string
}

var assets assetFileSystemI

func assetNames() (names []string) {
	return assets.AssetNames()
}

func readAsset(key string) ([]byte, error) {
	r, e := assets.Open(key)
	if e != nil {
		return nil, e
	}
	defer func() {
		_ = r.Close()
	}()

	p, e := ioutil.ReadAll(r)
	if e != nil {
		return nil, e
	}
	return p, nil
}

func mustReadAsset(key string) []byte {
	p, e := readAsset(key)
	if e != nil {
		panic("could not read asset with key " + key + ": " + e.Error())
	}
	return p
}

type assetOsFS struct{ root string }

func (aFS assetOsFS) Open(name string) (http.File, error) {
	return os.Open(filepath.Join(aFS.root, name))
}

func (aFS *assetOsFS) AssetNames() ([]string) {
	names, e := filepath.Glob(aFS.root + "/*")
	if e != nil {
		log.Print(e)
	}
	return names
}

type assetIntFS map[string][]byte

type assetFile struct {
	name string
	*bytes.Reader
}

type assetFileInfo struct {
	*assetFile
}

func (info assetFileInfo) Name() string {
	return info.assetFile.name
}

func (info assetFileInfo) ModTime() time.Time {
	return builtAt
}

func (info assetFileInfo) Mode() os.FileMode {
	return 0644
}

func (info assetFileInfo) Sys() interface{} {
	return nil
}

func (info assetFileInfo) Size() int64 {
	return int64(info.assetFile.Reader.Len())
}

func (info assetFileInfo) IsDir() bool {
	return false
}

func (info assetFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *assetFile) Stat() (os.FileInfo, error) {
	info := assetFileInfo{assetFile: f}
	return info, nil
}

func (afs assetIntFS) AssetNames() (names []string) {
	names = make([]string, 0, len(afs))
	for k, _ := range afs {
		names = append(names, k)
	}
	return names
}

func (afs assetIntFS) Open(name string) (af http.File, e error) {
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		name = "index.html"
	}
	if asset, found := afs[name]; found {
		decomp, e := gzip.NewReader(bytes.NewBuffer(asset))
		if e != nil {
			return nil, e
		}
		defer func() {
			_ = decomp.Close()
		}()
		b, e := ioutil.ReadAll(decomp)
		if e != nil {
			return nil, e
		}
		af = &assetFile{Reader: bytes.NewReader(b), name: name}
		return af, nil
	}
	return nil, os.ErrNotExist
}

func (a *assetFile) Close() error {
	return nil
}

func (a *assetFile) Read(p []byte) (n int, e error) {
	if a.Reader == nil {
		return 0, os.ErrInvalid
	}
	return a.Reader.Read(p)
}

func init() {
	builtAt = time.Now()
	env_name := fmt.Sprintf("GOASSETS_PATH")
	path := os.Getenv(env_name)
	if path != "" {
		stat, e := os.Stat(path)
		if e == nil && stat.IsDir() {
			assets = &assetOsFS{root: path}
			return
		}
	}

	assetsTmp := assetIntFS{}
	assetsTmp["a.js"] = []byte{
		0x1f,0x8b,0x8,0x0,0x0,0x9,0x6e,0x88,0x0,0xff,0x4a,0x2b,
0xcd,0x4b,0x2e,0xc9,0xcc,0xcf,0x53,0xc8,0x48,0xcd,0xc9,0xc9,
0xf,0xcf,0x2f,0xca,0x49,0xd1,0xd0,0x54,0xa8,0xe6,0xe2,0x4c,
0xce,0xcf,0x2b,0xce,0xcf,0x49,0xd5,0xcb,0xc9,0x4f,0xd7,0x50,
0x2,0x4b,0x2a,0x94,0x83,0x64,0x95,0x34,0xad,0xb9,0x6a,0xb9,
0x0,0x0,0x0,0x0,0xff,0xff,0x1,0x0,0x0,0xff,0xff,0x21,
0xd3,0x9d,0xf9,0x37,0x0,0x0,0x0,

	}
	assetsTmp["layout.html"] = []byte{
		0x1f,0x8b,0x8,0x0,0x0,0x9,0x6e,0x88,0x0,0xff,0x54,0x8f,
0x31,0x8f,0xc2,0x30,0xc,0x85,0xf7,0xfe,0xa,0x9f,0xf7,0xbb,
0xa8,0xdb,0xd,0x49,0x97,0x3b,0x36,0x24,0x18,0xba,0x30,0xa1,
0xd0,0xb8,0x4a,0x84,0x4b,0x51,0xe3,0x4a,0xe4,0xdf,0x93,0x26,
0x30,0x90,0xe5,0x29,0x7e,0xcf,0xf6,0x67,0xfd,0xf5,0x7f,0xf8,
0xeb,0x4f,0xc7,0x1d,0x78,0x99,0xb8,0x6b,0xf4,0x5b,0xc8,0xba,
0xae,0x81,0xfc,0xb4,0x4,0x61,0xea,0xf6,0x36,0xcd,0xab,0x40,
0x4f,0x51,0xb4,0xaa,0xa5,0x6a,0x4f,0x24,0x16,0x6,0x6f,0x97,
0x48,0x62,0x70,0x95,0xf1,0xfb,0x17,0x41,0xbd,0x4c,0xe,0xb7,
0x2b,0xf8,0x85,0x46,0x83,0x2a,0x8a,0x95,0x30,0x64,0x49,0x4c,
0x3f,0x43,0x8c,0x8,0x13,0xb9,0x60,0xd,0x5a,0x66,0x84,0x85,
0xd8,0x60,0xf1,0xa2,0x27,0x12,0x4,0x49,0x77,0x32,0x28,0xf4,
0x10,0x55,0xc2,0x79,0xa4,0x56,0x15,0x4b,0x5f,0x66,0x97,0xb6,
0xd,0xda,0xb7,0x10,0x9c,0x41,0x2e,0x70,0x67,0xc9,0x70,0xf8,
0x49,0xea,0xdb,0xad,0xad,0xe6,0xf3,0xaf,0x1c,0xf7,0x4,0x0,
0x0,0xff,0xff,0x1,0x0,0x0,0xff,0xff,0x9b,0xe3,0x63,0x1,
0xf4,0x0,0x0,0x0,

	}
	assetsTmp["style.css"] = []byte{
		0x1f,0x8b,0x8,0x0,0x0,0x9,0x6e,0x88,0x0,0xff,0x4a,0xca,
0x4f,0xa9,0x54,0xa8,0xe6,0xe2,0x4c,0x4a,0x4c,0xce,0x4e,0x2f,
0xca,0x2f,0xcd,0x4b,0xd1,0x4d,0xce,0xcf,0xc9,0x2f,0xb2,0x52,
0x48,0x49,0x2c,0xca,0x2e,0x4a,0x4d,0xb1,0xe6,0xe2,0x84,0xa,
0x94,0x67,0x64,0x96,0xa4,0x2,0xb9,0x69,0xf9,0x79,0x25,0xba,
0x69,0x89,0xb9,0x99,0x39,0x95,0x56,0xa,0x89,0x45,0x99,0x89,
0x39,0xd6,0x5c,0xb5,0x5c,0x0,0x0,0x0,0x0,0xff,0xff,0x1,
0x0,0x0,0xff,0xff,0xf7,0x31,0x38,0x22,0x49,0x0,0x0,0x0,

	}
	
	assets = assetsTmp
}
