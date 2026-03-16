package monster_siren

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	msrClient "github.com/AyakuraYuki/monster-siren-api-go/client"

	"github.com/go-resty/resty/v2"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/panjf2000/ants/v2"
)

var version string

const (
	saveTo = "./monster-siren"
)

type MonsterSiren struct {
	client     *msrClient.Client
	downloader *resty.Client
	pool       *ants.Pool
	progress   progress.Writer
}

func New(versions ...string) *MonsterSiren {
	if len(versions) > 0 && versions[0] != "" {
		version = versions[0]
	} else {
		version = "build-" + strings.ReplaceAll(time.Now().Format("20060102150405.000000"), ".", "_")
	}

	downloader := resty.New().
		SetTimeout(20 * time.Minute).
		SetRetryCount(3).
		SetHeaders(map[string]string{
			"Accept":          "*/*",
			"Accept-Language": "zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,en-GB;q=0.6,en-US;q=0.5",
			"Referer":         "https://monster-siren.hypergryph.com/",
			"User-Agent":      fmt.Sprintf("Go/%s monster-siren-downloader/%s", runtime.Version(), version),
		})

	instance := &MonsterSiren{
		client:     msrClient.NewClient(),
		downloader: downloader,
		progress:   customizedProgress(),
	}

	pool, err := ants.NewPool(5, ants.WithPanicHandler(instance.antsPanicHandler))
	if err != nil {
		panic(err)
	}
	instance.pool = pool

	return instance
}

func (m *MonsterSiren) antsPanicHandler(_ any) {
	buf := make([]byte, 4<<10) // 4K
	buf = buf[:runtime.Stack(buf, false)]
	m.progress.Log("panic: %s", string(buf))
}
