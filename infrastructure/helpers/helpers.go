package helpers

import (
	"fmt"
	"go-api/infrastructure/myLog"
	"io"
)

func CloseSafely(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		message := fmt.Sprintf("Failed to close: \n%v", err)
		myLog.Info.SkipLog(1, message)
	}
}
