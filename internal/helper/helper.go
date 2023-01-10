package helper

import "github.com/sirupsen/logrus"

// WrapCloser call close and log the error
func WrapCloser(close func() error) {
	if close == nil {
		return
	}
	if err := close(); err != nil {
		logrus.Error(err)
	}
}
