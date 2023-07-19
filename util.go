/*
Copyright 2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nftables

import (
	"fmt"
	"strconv"
	"strings"
)

// Optional can be used to fill in optional field values in objects
func Optional[T any](val T) *T {
	return &val
}

var numericPriorities = map[string]int{
	"raw":      -300,
	"mangle":   -150,
	"dstnat":   -100,
	"filter":   0,
	"security": 50,
	"srcnat":   100,
}

var bridgeNumericPriorities = map[string]int{
	"dstnat": -300,
	"filter": -200,
	"out":    100,
	"srcnat": 300,
}

// ParsePriority tries to convert the string form of a chain priority into a number
func ParsePriority(family Family, priority string) (int, error) {
	val, err := strconv.Atoi(priority)
	if err == nil {
		return val, nil
	}

	modVal := 0
	if i := strings.IndexAny(priority, "+-"); i != -1 {
		mod := priority[i:]
		modVal, err = strconv.Atoi(mod)
		if err != nil {
			return 0, fmt.Errorf("could not parse modifier %q: %v", mod, err)
		}
		priority = priority[:i]
	}

	var found bool
	if family == BridgeFamily {
		val, found = bridgeNumericPriorities[priority]
	} else {
		val, found = numericPriorities[priority]
	}
	if !found {
		return 0, fmt.Errorf("unknown priority %q", priority)
	}

	return val + modVal, nil
}

// Join joins multiple string values together into a multi-valued set/map key/value.
func Join(values ...string) string {
	return strings.Join(values, " . ")
}

// Split splits an Element.Key or Element.Value into its component parts
func Split(values string) []string {
	return strings.Split(values, " . ")
}