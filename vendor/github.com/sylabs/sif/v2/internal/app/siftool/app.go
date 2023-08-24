// Copyright (c) 2021-2022, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package siftool

import (
	"io"
	"os"
)

// appOpts contains configured options.
type appOpts struct {
	out io.Writer
	err io.Writer
}

// AppOpt are used to configure optional behavior.
type AppOpt func(*appOpts) error

// App holds state and configured options.
type App struct {
	opts appOpts
}

// OptAppOutput specifies that output should be written to w.
func OptAppOutput(w io.Writer) AppOpt {
	return func(o *appOpts) error {
		o.out = w
		return nil
	}
}

// OptAppError specifies that errors should be written to w.
func OptAppError(w io.Writer) AppOpt {
	return func(o *appOpts) error {
		o.err = w
		return nil
	}
}

// New creates a new App configured with opts.
//
// By default, application output and errors are written to os.Stdout and os.Stderr respectively.
// To modify this behavior, consider using OptAppOutput and/or OptAppError.
func New(opts ...AppOpt) (*App, error) {
	a := App{
		opts: appOpts{
			out: os.Stdout,
			err: os.Stderr,
		},
	}

	for _, opt := range opts {
		if err := opt(&a.opts); err != nil {
			return nil, err
		}
	}

	return &a, nil
}
