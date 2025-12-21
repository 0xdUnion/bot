package content

import (
	"fmt"
	"time"
)

func TimeStamp() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("当前时间戳为: `%d`", timestamp)

}
