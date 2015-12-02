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

// Modified from original version.
// Iso 8601 is the standard for date/time exchange

package util

import (
	"time"
)

// Iso8601DateTime is a type for decoding and encoding json
// date times that follow Iso 8601 format. The type currently
// decodes and encodes with exactly precision to seconds. If more
// formats of Iso8601 need to be supported additional work
// will be needed.
type Iso8601DateTime struct {
	time.Time
}

// NewDateTime creates a new Iso8601DateTime taking a string as input.
// It must follow the "2006-01-02T15:04:05" pattern.
func NewIso8601DateTime(input string) (val *Iso8601DateTime, err error) {
	val = &Iso8601DateTime{}
	err = val.Parse(input)
	if err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

// Parse converts string give to a Iso8601DateTime object.
// Errors will occur if the string when converted to a string
// doesn't match the format "2006-01-02T15:04:05".
func (r *Iso8601DateTime) Parse(data string) error {
	//Check last character is Z or not
	
	timeVal, err := time.Parse(format, data)
	if err != nil {
		return err
	}
	r.Time = timeVal
	return nil
}

// MarshalJSON converts a Iso8601DateTime to a string.
func (r *Iso8601DateTime) String() string {
	return r.Time.Format(format)
}

func (r *Iso8601DateTime) UnmarshalJSON(data []byte) error {
	return r.Parse(string(data))
}

func (r *Iso8601DateTime) MarshalJSON() ([]byte, error) {
	val := r.Time.Format(format)
	return []byte(val), nil
}

const format = `"2006-01-02T15:04:05.999999Z"`
