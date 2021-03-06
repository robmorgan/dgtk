package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	builtAt        = time.Now()
	compiledAssets = assetIntFS{}
	DevPath        string
	assets         assetFileSystemI
)

type assetProxy struct {
}

func (a *assetProxy) Open(name string) (http.File, error) {
	return a.fileSystem().Open(name)
}

func (a *assetProxy) AssetNames() []string {
	return a.fileSystem().AssetNames()
}

func (a *assetProxy) fileSystem() assetFileSystemI {
	if DevPath != "" {
		stat, e := os.Stat(DevPath)
		if e == nil && stat.IsDir() {
			assets = &assetOsFS{root: DevPath}
			return assets
		}
	}
	return compiledAssets
}

func FileSystem() assetFileSystemI {
	return &assetProxy{}
}

type assetFileSystemI interface {
	Open(name string) (http.File, error)
	AssetNames() []string
}

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

func (aFS *assetOsFS) AssetNames() []string {
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
