package nsmvpp

import (
	"fmt"
	"strconv"

	govppapi "git.fd.io/govpp.git/api"
	"github.com/ligato/networkservicemesh/dataplanes/vpp/bin_api/l2"
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
	// Testing API calls, just to explore weird API call stuff

	logrus.Infof("connecting inexisting interfaces...")
	xconnectReq := l2.SwInterfaceSetL2Xconnect{
		RxSwIfIndex: 100,
		TxSwIfIndex: 101,
		Enable:      1,
	}
	xconnectRpl := l2.SwInterfaceSetL2XconnectReply{}
	if err := apiCh.SendRequest(&xconnectReq).ReceiveReply(&xconnectRpl); err != nil {
		logrus.Errorf("error: %v", err)
	} else {
		logrus.Errorf("No error, but should be error")
	}

	logrus.Infof("connecting inexisting interfaces again...")
	xconnectReq = l2.SwInterfaceSetL2Xconnect{
		RxSwIfIndex: 100,
		TxSwIfIndex: 101,
		Enable:      1,
	}
	xconnectRpl = l2.SwInterfaceSetL2XconnectReply{}
	if err := apiCh.SendRequest(&xconnectReq).ReceiveReply(&xconnectRpl); err != nil {
		logrus.Errorf("error: %v", err)
	} else {
		logrus.Errorf("No error, but should be error")
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

	pos, err := perform(tx, apiCh)
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
