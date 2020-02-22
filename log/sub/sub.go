// sub
package sub

import (
	"../mylog"
)

func LogTest() {
	logger := mylog.Newlog()
	logger.Info("sub LogTest")
	logger.Warn("sub LogTest")
}
