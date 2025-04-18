package types

import (
	"fmt"
)

func GenItemKey(taskId uint32) string {
	return fmt.Sprintf("%d", taskId)
}
