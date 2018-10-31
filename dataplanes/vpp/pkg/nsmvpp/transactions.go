package nsmvpp

import (
	govppapi "git.fd.io/govpp.git/api"
	govpp "git.fd.io/govpp.git/core"
	"github.com/sirupsen/logrus"
)

type operation interface {
	apply(apiCh govppapi.Channel) error
	rollback() operation
}

func rollback(tx []operation, pos int, apiCh govppapi.Channel) error {
	logrus.Infof("Rolling back operations...")
	var err error
	for i := pos - 1; i >= 0; i-- {
		err = tx[i].rollback().apply(apiCh)
		if err != nil {
			logrus.Errorf("error while rolling back, (I will continue rollback operations): %v", err)
		}
	}
	logrus.Info("Done. I did my best to roll things back")
	return err
}

func perform(tx []operation, conn *govpp.Connection) (int, error) {
	logrus.Infof("Programming dataplane...")
	for i := range tx {
		apiCh, err := conn.NewAPIChannel()
		if err != nil {
			logrus.Errorf("Problems creating API channel: %v", err)
		}
		err = tx[i].apply(apiCh)
		if err != nil {
			logrus.Errorf("error performing operation %v", err)
			return i, err
		}
		apiCh.Close()
	}
	logrus.Infof("Transaction completed!")
	return len(tx), nil
}
