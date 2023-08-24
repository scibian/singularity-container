// Copyright (c) 2021-2022, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package siftool

import (
	"os"

	"github.com/sylabs/sif/v2/pkg/sif"
)

// withFileImage calls fn with a FileImage loaded from path.
func withFileImage(path string, writable bool, fn func(*sif.FileImage) error) error {
	flag := os.O_RDONLY
	if writable {
		flag = os.O_RDWR
	}

	f, err := sif.LoadContainerFromPath(path, sif.OptLoadWithFlag(flag))
	if err != nil {
		return err
	}

	err = fn(f)

	if uerr := f.UnloadContainer(); err == nil {
		err = uerr
	}

	return err
}
