package hunter

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sunface/tools"
)

func (h *Hunter) daemon() {
	go func() {
		Logger.Info("start daemon")
		for {
			hour := time.Now().Hour()

			var pt int
			if Conf.Hunter.PullTime == 0 {
				pt = 12
			} else {
				pt = Conf.Hunter.PullTime
			}

			if pt == hour { // 在pull的时间范围内
				//判断需要pull的文件是否已经存在
				year, mon, day := time.Now().Date()
				h.Date = fmt.Sprintf("%04d", year) + "-" + fmt.Sprintf("%02d", mon) + "-" + fmt.Sprintf("%02d", day)

				fn := Conf.Hunter.OutPath + "/" + h.Date + ".out"
				if tools.FileExist(fn) { //若存在，则继续循环
					Logger.Info("in pull time,but file exist")
					goto INTV
				} else {
					Logger.Info("start to pull")
					h.hunter()
				}
			}

		INTV:
			time.Sleep(15 * time.Second)
		}
	}()

	// 等待服务器停止信号
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	<-chSig
	Logger.Info("stop daemon")
}
