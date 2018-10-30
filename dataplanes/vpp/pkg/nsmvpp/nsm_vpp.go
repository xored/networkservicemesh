// Copyright (c) 2018 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nsmvpp

import (
	"os"

	"git.fd.io/govpp.git/adapter/vppapiclient"
	govppapi "git.fd.io/govpp.git/api"
	govpp "git.fd.io/govpp.git/core"
	"github.com/sirupsen/logrus"
)

// Interface lists methods available to manipulate VPPDataplane controller information
type Interface interface {
	GetAPIChannel() govppapi.Channel
	GetDataplaneSocket() string
	IsConnected() bool
}

// VPPDataplane defines fields of NSM VPP dataplane controller
type VPPDataplane struct {
	dataplaneSocket string
}

func (v *VPPDataplane) GetDataplaneSocket() string {
	return v.dataplaneSocket
}

func (v *VPPDataplane) IsConnected() bool {
	return true
}

// GetAPIChannel returns VPP Dataplane API channel. API channel is used by dataplane programming
// functions.
func (v *VPPDataplane) GetAPIChannel() govppapi.Channel {
	vppConn, vppConnCh, err := govpp.AsyncConnect(vppapiclient.NewVppAdapter(""))
	if err != nil {
		logrus.Errorf("Failed to reconnect VPP with error: %+v retrying in %s", err, "0")
		os.Exit(1)
	}
	status := <-vppConnCh
	if status.State != govpp.Connected {
		logrus.Errorf("Timed out to reconnect to VPP, retrying in %s", "0")
		os.Exit(1)
	}
	// Locking VPPDataplane for updating some fields
	apiCh, err := vppConn.NewAPIChannel()
	if err != nil {
		logrus.Errorf("Failed to get API channel, retrying in %s", "0")
		os.Exit(1)
	}

	return apiCh
}

// NEWVPPDataplane starts VPP binary, waits until it is ready and populate
// VPPDataplane controller structure.
func NEWVPPDataplane(dataplaneSocket string) (Interface, error) {
	VPPDataplaneController := &VPPDataplane{
		dataplaneSocket: dataplaneSocket,
	}

	return VPPDataplaneController, nil
}
