// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

// Package testUtil has helpers to be used with unit tests
package testutil

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Equals fails the test if exp is not equal to act.
// Code was copied from https://github.com/benbjohnson/testing MIT license
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// Assert fails the test if the condition is false.
// Code was copied from https://github.com/benbjohnson/testing MIT license
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// IsNil ensures that the act interface is nil
// otherwise an error is raised.
func IsNil(tb testing.TB, act interface{}) {
	if act != nil {
		tb.Error("expected nil", act)
		tb.FailNow()
	}
}