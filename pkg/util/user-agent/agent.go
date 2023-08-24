// Copyright (c) 2018-2022, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package useragent

import (
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var value string

// Value contains the Singularity user agent.
//
// For example, "Singularity/3.0.0 (linux amd64) Go/1.10.3".
func Value() string {
	if value == "" {
		panic("useragent.InitValue() must be called before calling useragent.Value()")
	}

	return value
}

// InitValue sets value that will be returned when
// user queries singularity version.
func InitValue(name, version string) {
	value = fmt.Sprintf("%v (%v %v) %v",
		singularityVersion(name, version),
		cases.Title(language.English).String(runtime.GOOS),
		runtime.GOARCH,
		goVersion())
}

func singularityVersion(name, version string) string {
	product := cases.Title(language.English).String(name)
	ver := strings.Split(version, "-")[0]
	return fmt.Sprintf("%v/%v", product, ver)
}

func goVersion() string {
	version := strings.TrimPrefix(runtime.Version(), "go")
	return fmt.Sprintf("Go/%v", version)
}
