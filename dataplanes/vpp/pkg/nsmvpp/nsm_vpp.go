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
	"sync"
	"time"

	"git.fd.io/govpp.git/adapter/vppapiclient"
	govppapi "git.fd.io/govpp.git/api"
	govpp "git.fd.io/govpp.git/core"
	"github.com/sirupsen/logrus"
)

var (
	vppReconnectInterval = time.Second * 30
)

// Interface lists methods available to manipulate VPPDataplane controller information
type Interface interface {
	GetDataplaneSocket() string
	IsConnected() bool
	SetRegistered()
	SetUnRegistered()
	Test() error
	Shutdown()
	BreakConnection()
	SetUnRegisterCallback(f func(vpp Interface))
	GetAPIChannel() govppapi.Channel
}

// VPPDataplane defines fields of NSM VPP dataplane controller
type VPPDataplane struct {
	// conn   *govpp.Connection
	// connCh chan govpp.ConnectionEvent
	// status *govpp.ConnectionEvent
	// apiCh  govppapi.Channel
	sync.RWMutex
	nsmRegistered       bool
	dataplaneSocket     string
	dataplaneUnregister func(vpp Interface)
}

// GetAPIChannel returns VPP Dataplane API channel. API channel is used by dataplane programming
// functions.
func (v *VPPDataplane) GetAPIChannel() govppapi.Channel {
	logrus.Infof(" ***** Connecting to VPP ***** ")
	vppConn, vppConnCh, err := govpp.AsyncConnect(vppapiclient.NewVppAdapter(""))
	if err != nil {
		logrus.Errorf("Failed to reconnect VPP with error: %+v retrying in %s", err, vppReconnectInterval.String())
		os.Exit(1)
	}
	status := <-vppConnCh
	if status.State != govpp.Connected {
		logrus.Errorf("Timed out to reconnect to VPP, retrying in %s", vppReconnectInterval.String())
		os.Exit(1)
	}
	ch, err := vppConn.NewAPIChannel()
	if err != nil {
		logrus.Errorf("failed open new channel: %v", err)
		return nil
	}
	return ch
}

// SetUnRegisterCallback sets a callback function which will be called upon detection
// VPP Disconnected event. It is used to prevent NSM to call the dataplane controller
// while VPP is not functional/connected.
func (v *VPPDataplane) SetUnRegisterCallback(f func(vpp Interface)) {
	v.Lock()
	defer v.Unlock()
	v.dataplaneUnregister = f
}

// GetDataplaneSocket returns dataplane socket location, this dataplane controller
// will be service requests.
func (v *VPPDataplane) GetDataplaneSocket() string {
	return v.dataplaneSocket
}

// IsConnected returns true if VPP state is connected
func (v *VPPDataplane) IsConnected() bool {
	return true
}

// SetRegistered marks VPP Dataplane controller as registered with NSM
func (v *VPPDataplane) SetRegistered() {
	v.Lock()
	defer v.Unlock()
	v.nsmRegistered = true
}

// SetUnRegistered marks VPP Dataplane controller as registered with NSM
func (v *VPPDataplane) SetUnRegistered() {
	v.Lock()
	defer v.Unlock()
	v.nsmRegistered = false

}

// Shutdown shuts down api channel and closes connection with VPP.
func (v *VPPDataplane) Shutdown() {
	//v.apiCh.Close()
	// v.conn.Disconnect()
}

// eventMonitor listens for Disconnected event, upon receiving it,
// 1. Changes status of VPPDataplane controller,
// 2. Calls for NSM dataplane unregister function, so while vpp is not connected
//    nsm does not try to communicate with VPPDataplane controller. TODO (sbezverk)
// 3. Starts vpp reconnector function which will attempt to re-connect to VPP
// 4. Exits, new monitor will be started once the connection gets re-established.
func (v *VPPDataplane) eventMonitor() {
	logrus.Info("Starting event monitor")
}

// reConnector is called once Disconnect message is recieved by the event monitor. It will infinetly
// attempts to re connect to VPP, once it is succeeded, it will mark VPPDataplane controller as Connected
// and start the dataplane registration with NSM function TODO (sbezverk).
func (v *VPPDataplane) reConnector() {
	// ticker := time.NewTicker(vppReconnectInterval)
	// startTime := time.Now()
}

// NEWVPPDataplane starts VPP binary, waits until it is ready and populate
// VPPDataplane controller structure.
func NEWVPPDataplane(dataplaneSocket string) (Interface, error) {
	// startTime := time.Now()
	// vppConn, vppConnCh, err := govpp.AsyncConnect(vppapiclient.NewVppAdapter(""))
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to start NSM VPP Dataplaneagent with error:%+v", err)
	// }

	// status := <-vppConnCh
	// if status.State != govpp.Connected {
	// 	return nil, fmt.Errorf("Failed to start NSM VPP Dataplaneagent with error:%+v", err)
	// }
	// vppConnectTime := time.Since(startTime)
	// logrus.Info("Connecting to VPP took ", vppConnectTime)

	VPPDataplaneController := &VPPDataplane{
		// apiCh:           apiCh,
		nsmRegistered:   false,
		dataplaneSocket: dataplaneSocket,
	}
	// Starting VPP event monitor routine
	go VPPDataplaneController.eventMonitor()

	return VPPDataplaneController, nil
}

// Test is used only in Debug mode, it runs some common api
// to confirm VPP is fully functional
func (v *VPPDataplane) Test() error {
	// Playground
	// End of playground
	return nil
}

// BreakConnection is used only for debugging mode to simulate Disconnected
// message from VPP, to see how NSM VPP dataplane controller behaves
func (v *VPPDataplane) BreakConnection() {
	// v.conn.Disconnect()
	// v.connCh <- govpp.ConnectionEvent{
	// 	Timestamp: time.Now(),
	// 	State:     govpp.Disconnected,
	// 	Error:     fmt.Errorf("Simulating VPP disconnect"),
	// }
}
