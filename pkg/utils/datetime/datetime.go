package datetime

import (
	"fmt"
	"time"
)

func GetDateTimeStr() string {
	dt := time.Now().UTC()

	return fmt.Sprintf("%s", dt.Format("2006-01-02 15:04:05"))
}
