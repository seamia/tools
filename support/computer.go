// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package support

import (
	"net"
	"sort"
	"strings"
)

func GetComputerID() string {

	is, err := net.Interfaces()
	if err != nil {
		// let's use the description of the error as 'id'
		return err.Error()
	}

	collection := make([]string, 0, 12)
	for _, ifi := range is {

		if ifi.Flags&net.FlagUp == 0 || ifi.Flags&net.FlagLoopback != 0 || ifi.Flags&net.FlagPointToPoint != 0 {
			// ignoring 'downed', loopback and p2p interfaces
			continue
		}

		if ifi.Flags&net.FlagBroadcast != 0 && ifi.Flags&net.FlagMulticast != 0 {
			// fmt.Printf("%s. addr: %v, hash: %s\n", ifi.Name, ifi.HardwareAddr, hash(ifi.HardwareAddr))
			collection = append(collection, Hash(ifi.HardwareAddr))
		}
	}

	// arrange potential multitude of 'strings' in a deterministic manner
	if len(collection) > 0 {
		sort.Strings(collection)
	} else {
		collection = append(collection, "Homeless device.")
	}

	// prevent 'flooding' attack
	if len(collection) > 3 {
		collection = collection[:3]
	}

	return strings.Join(collection, "|")
}
