package client

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var buildVersion string
var buildVersionOnce sync.Once

func init() {
	if buildVersion == "" {
		buildVersionOnce.Do(func() {
			timeTag := time.Now().Format("20060102150405.000000")
			timeTag = strings.ReplaceAll(timeTag, ".", "_")
			buildVersion = fmt.Sprintf("dev-%s", timeTag)
		})
	}
}

func Version() string {
	return buildVersion
}
