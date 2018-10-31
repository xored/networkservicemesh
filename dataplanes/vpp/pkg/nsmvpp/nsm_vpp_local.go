package nsmvpp

import (
	"fmt"
	"os"
	"strconv"

	"git.fd.io/govpp.git/adapter/vppapiclient"
	govppapi "git.fd.io/govpp.git/api"
	govpp "git.fd.io/govpp.git/core"
	"github.com/ligato/networkservicemesh/pkg/nsm/apis/common"
	"github.com/sirupsen/logrus"
)

type vppInterface struct {
	mechanism *common.LocalMechanism // we want to save parameters here in order to recreate interface
	id        uint32
}

var (
	connections map[int][]operation = make(map[int][]operation)
	lastId      int                 = 0
)

// CreateLocalConnect sanity checks parameters passed in the LocalMechanisms and call nsmvpp.CreateLocalConnect
func CreateLocalConnect(apiCh govppapi.Channel, src, dst *common.LocalMechanism) (string, error) {
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
	srcIntf := &vppInterface{}
	dstIntf := &vppInterface{}

	tx := []operation{
		&createLocalInterface{
			localMechanism: src,
			intf:           srcIntf,
		},
		&createLocalInterface{
			localMechanism: dst,
			intf:           dstIntf,
		},
		&interfaceXconnect{
			rx:     srcIntf,
			tx:     dstIntf,
			enable: 1,
		},
		&interfaceXconnect{
			rx:     dstIntf,
			tx:     srcIntf,
			enable: 1,
		},
		&interfaceUpDown{
			intf:   srcIntf,
			upDown: 1,
		},
		&interfaceUpDown{
			intf:   dstIntf,
			upDown: 1,
		},
	}

	if (src.Type == common.LocalMechanismType_MEM_INTERFACE) && (dst.Type == common.LocalMechanismType_MEM_INTERFACE) {
		var err error
		if tx, err = memifDirectConnect(src.Parameters, dst.Parameters); err != nil {
			return "", err
		}
	}

	pos, err := perform(tx, vppConn)
	if err != nil {
		rollback(tx, pos, apiCh)
		return "", err
	}

	lastId++
	connections[lastId] = tx // save transaction log to perform rollback on delete connection
	return fmt.Sprintf("%d", lastId), nil
}

// DeleteLocalConnect
func DeleteLocalConnect(apiCh govppapi.Channel, connID string) error {
	id, _ := strconv.Atoi(connID)
	tx := connections[id]
	return rollback(tx, len(tx), apiCh)
}
