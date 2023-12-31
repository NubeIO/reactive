package reactive

import "github.com/sirupsen/logrus"

func (n *BaseNode) AddLogger(logger *logrus.Logger) {
	n.logger = logger
}
func (n *BaseNode) GetLogger() *logrus.Logger {
	return n.logger
}
