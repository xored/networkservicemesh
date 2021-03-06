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

package nsmdataplane

import (
	"context"
	"fmt"
	"net"

	"github.com/ligato/networkservicemesh/dataplanes/vpp/pkg/nsmvpp"
	"github.com/ligato/networkservicemesh/pkg/nsm/apis/common"
	dataplaneapi "github.com/ligato/networkservicemesh/pkg/nsm/apis/dataplane"
	"github.com/ligato/networkservicemesh/pkg/tools"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes" // "context"
	"google.golang.org/grpc/status"
)

// DataplaneServer keeps k8s client and gRPC server
type DataplaneServer struct {
	vppDataplane     nsmvpp.Interface
	server           *grpc.Server
	remoteMechanisms []*common.RemoteMechanism
	localMechanisms  []*common.LocalMechanism
	updateCh         chan Update
}

// Update is a message used to communicate any changes in operational parameters and constraints
type Update struct {
	remoteMechanisms []*common.RemoteMechanism
	localMechanisms  []*common.LocalMechanism
}

// ConnectRequest is called when NSM sends the request to interconnect two containers' namespaces.
func (d DataplaneServer) ConnectRequest(ctx context.Context, req *dataplaneapi.Connection) (*dataplaneapi.Reply, error) {
	logrus.Infof("ConnectRequest was called for namespace %+v", req)

	// First check, is VPP is operational? If not return grpc error dataplane is not available
	if !d.vppDataplane.IsConnected() {
		// VPP is not currently connected, failing this request.
		return &dataplaneapi.Reply{
			Success: false,
		}, status.Error(codes.Unavailable, "VPP dataplane is currently unavailable.")
	}

	// TODO (sbezverk) need to add a check for requested local and remote mechanisms. The dataplane controller
	// configuration should be global as it is needed by multiple packages.

	// There are two known types of destinations:
	// Local - when both PODs are running on the same host,
	// Remote - when NSM Client is local, but the requested Network Service Endpoint is runing on a
	// different from the client host.
	switch req.Destination.(type) {
	case *dataplaneapi.Connection_Local:
		logrus.Infof("Destination is local: %+v", req)
		destination := req.Destination.(*dataplaneapi.Connection_Local)
		logrus.Infof("Destination struct: %+v", destination.Local)
		connID, err := nsmvpp.CreateLocalConnect(d.vppDataplane.GetAPIChannel(), req.LocalSource, destination.Local)
		if err != nil {
			errStr := fmt.Sprintf("fail to build the cross connect with error: %+v", err)
			return &dataplaneapi.Reply{
				Success: false,
			}, status.Error(codes.Unavailable, errStr)
		}
		req.ConnectionId = connID
		return &dataplaneapi.Reply{
			Success:    true,
			Connection: req,
		}, nil

	case *dataplaneapi.Connection_Remote:
		logrus.Infof("Destination is remote: %+v", req)
		destination := req.Destination.(*dataplaneapi.Connection_Remote)
		logrus.Infof("Destination struct: %+v", destination.Remote)
		// Remote destination support is not yet implemented, failing this request
		return &dataplaneapi.Reply{
			Success: false,
		}, status.Error(codes.Unavailable, "Remote Destination currently is not supported")
	default:
		// Destination type does not match to any known/supported types, failing this request.
		return &dataplaneapi.Reply{
			Success: false,
		}, status.Error(codes.Unknown, "Unknown destination type")
	}
}

// DisconnectRequest is called when NSM sends the request to disconnect two containers' namespaces.
func (d DataplaneServer) DisconnectRequest(ctx context.Context, req *dataplaneapi.Connection) (*dataplaneapi.Reply, error) {
	logrus.Infof("DisconnectRequest was called for namespace %+v", req)
	// First check, is VPP is operational? If not return grpc error dataplane is not available
	if !d.vppDataplane.IsConnected() {
		// VPP is not currently connected, failing this request.
		return &dataplaneapi.Reply{
			Success: false,
		}, status.Error(codes.Unavailable, "VPP dataplane is currently unavailable.")
	}

	switch req.Destination.(type) {
	case *dataplaneapi.Connection_Local:
		logrus.Infof("Destination is local: %+v", req)
		destination := req.Destination.(*dataplaneapi.Connection_Local)
		logrus.Infof("Destination struct: %+v", destination.Local)

		if err := nsmvpp.DeleteLocalConnect(d.vppDataplane.GetAPIChannel(), req.ConnectionId); err != nil {
			errStr := fmt.Sprintf("fail to delete the cross connect with error: %+v", err)
			return &dataplaneapi.Reply{
				Success: false,
			}, status.Error(codes.Unavailable, errStr)
		}
		return &dataplaneapi.Reply{
			Success: true,
		}, nil
	case *dataplaneapi.Connection_Remote:
		logrus.Infof("Destination is remote: %+v", req)
		destination := req.Destination.(*dataplaneapi.Connection_Remote)
		logrus.Infof("Destination struct: %+v", destination.Remote)
		// Remote destination support is not yet implemented, failing this request
		return &dataplaneapi.Reply{
			Success: false,
		}, status.Error(codes.Unavailable, "Remote Destination currently is not supported")
	default:
		// Destination type does not match to any known/supported types, failing this request.
		return &dataplaneapi.Reply{
			Success: false,
		}, status.Error(codes.Unknown, "Unknown destination type")
	}

	return &dataplaneapi.Reply{Success: true}, nil
}

// UpdateDataplane implements method of dataplane interface, which is informing NSM of any changes
// to operational prameters or constraints
func (d DataplaneServer) MonitorMechanisms(empty *common.Empty, updateSrv dataplaneapi.DataplaneOperations_MonitorMechanismsServer) error {
	logrus.Infof("Update dataplane was called")
	if err := updateSrv.Send(&dataplaneapi.MechanismUpdate{
		RemoteMechanisms: d.remoteMechanisms,
		LocalMechanisms:  d.localMechanisms,
	}); err != nil {
		logrus.Errorf("vpp dataplane server: Deteced error %s, grpc code: %+v on grpc channel", err.Error(), status.Convert(err).Code())
		return nil
	}
	for {
		select {
		// Waiting for any updates which might occur during a life of dataplane module and communicating
		// them back to NSM.
		case update := <-d.updateCh:
			d.remoteMechanisms = update.remoteMechanisms
			d.localMechanisms = update.localMechanisms
			if err := updateSrv.Send(&dataplaneapi.MechanismUpdate{
				RemoteMechanisms: update.remoteMechanisms,
				LocalMechanisms:  update.localMechanisms,
			}); err != nil {
				logrus.Errorf("vpp dataplane server: Deteced error %s, grpc code: %+v on grpc channel", err.Error(), status.Convert(err).Code())
				return nil
			}
		}
	}
}

// StartDataplaneServer instantiates gRPC server to serve NSM dataplane programming requests
func StartDataplaneServer(vpp nsmvpp.Interface) error {
	//  Start server on our dataplane socket.

	dataplaneServer := DataplaneServer{
		updateCh:     make(chan Update),
		vppDataplane: vpp,
		localMechanisms: []*common.LocalMechanism{
			&common.LocalMechanism{
				Type: common.LocalMechanismType_KERNEL_INTERFACE,
			},
			&common.LocalMechanism{
				Type: common.LocalMechanismType_MEM_INTERFACE,
			},
		},
		remoteMechanisms: []*common.RemoteMechanism{
			&common.RemoteMechanism{
				Type: common.RemoteMechanismType_VXLAN,
			},
		},
	}
	dataplaneSocket := vpp.GetDataplaneSocket()

	if err := tools.SocketCleanup(dataplaneSocket); err != nil {
		return fmt.Errorf("vpp dataplane server: failure to cleanup stale socket %s with error: %+v", dataplaneSocket, err)
	}
	dataplaneConn, err := net.Listen("unix", dataplaneSocket)
	if err != nil {
		return fmt.Errorf("vpp dataplane server: fail to open socket %s with error: %+v", dataplaneSocket, err)
	}
	dataplaneServer.server = grpc.NewServer()
	// Binding dataplane Interface API to gRPC server
	dataplaneapi.RegisterDataplaneOperationsServer(dataplaneServer.server, dataplaneServer)

	// Starting gRPC server, if there is something wrong with starting it, it will be caught by following gRPC test
	go dataplaneServer.server.Serve(dataplaneConn)

	// Check if the socket of device plugin server is operation
	testSocket, err := tools.SocketOperationCheck(dataplaneSocket)
	if err != nil {
		return fmt.Errorf("vpp dataplane server: failure to communicate with the socket %s with error: %+v", dataplaneSocket, err)
	}
	defer testSocket.Close()
	logrus.Infof("vpp dataplane server: Test Dataplane controller is ready to serve...")

	return nil
}
