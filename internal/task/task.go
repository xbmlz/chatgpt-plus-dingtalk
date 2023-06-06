package task

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/config"
)

func Initialize() {
	c := cron.New()
	err := c.AddFunc(config.Instance.CleanAllSessionCron, CleanAllSessionInterval)
	if err != nil {
		fmt.Println("添加定时任务失败：", err)
		return
	}
	c.Start()
}
