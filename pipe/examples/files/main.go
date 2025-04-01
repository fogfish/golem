//
// Copyright (C) 2022 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fogfish/golem/pipe/v2/fork"
)

var (
	fsys    = os.DirFS("../../../")
	threads = runtime.NumCPU()
)

func main() {
	ctx, close := context.WithCancel(context.Background())
	defer close()

	// Parallel SHA1 digest
	seq, errf := walk(ctx)
	sha, errh := fork.Map(ctx, threads, seq, fork.Lift(digest))
	str, errs := fork.Map(ctx, threads, sha, fork.Lift(stringify))
	<-fork.ForEach(ctx, threads, str,
		fork.Pure(
			func(x string) string {
				fmt.Printf("==> %s\n", x)
				return x
			},
		),
	)

	if err := <-fork.Join(ctx, errf, errh, errs); err != nil {
		fmt.Printf("failed: %s\n", err)
	}
}

// Walk pipeline visit each file in the mounted file system.
func walk(ctx context.Context) (<-chan string, <-chan error) {
	out := make(chan string)
	exx := make(chan error, 1)

	go func() {
		defer close(out)
		exx <- fs.WalkDir(fsys, ".",
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() || !strings.HasSuffix(path, ".go") {
					return nil
				}

				select {
				case out <- path:
				case <-ctx.Done():
					return fs.SkipAll
				}
				return nil
			},
		)
	}()

	return out, exx
}

// Digest file content into SHA1 hash
func digest(path string) (*hash, error) {
	fd, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}

	h := sha1.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return nil, err
	}

	return &hash{
		path: path,
		sha1: h.Sum(nil),
	}, nil
}

type hash struct {
	path string
	sha1 []byte
}

// Convert digest to string
func stringify(h *hash) (string, error) {
	return fmt.Sprintf("%x %s", h.sha1, filepath.Base(h.path)), nil
}
