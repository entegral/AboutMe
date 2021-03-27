package errors

import "github.com/sirupsen/logrus"

func Warn(id string, err error) bool {
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Warnf("ERROR: %s", err.Error())
		return true
	}
	return false
}