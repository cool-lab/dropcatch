package hunter

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sunface/tools"
	"github.com/tealeg/xlsx"
	"github.com/uber-go/zap"
)

type Hunter struct {
	Date      string
	ExcelName string
	Excel     *xlsx.File
	Out       *os.File
}

func New() *Hunter {
	return &Hunter{}
}

var wg = &sync.WaitGroup{}

func (h *Hunter) Start() {
	// 初始化zap logger
	InitLogger(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)

	if Conf.Hunter.IsDaemon {
		h.daemon()
	} else {
		if Conf.Hunter.Date == "" {
			year, mon, day := time.Now().Date()
			h.Date = fmt.Sprintf("%04d", year) + "-" + fmt.Sprintf("%02d", mon) + "-" + fmt.Sprintf("%02d", day)
		} else {
			h.Date = Conf.Hunter.Date
		}

		h.hunter()
	}

}

func (h *Hunter) DownLoad() {
	url := "https://www.dropcatch.com/DownloadCenter/ExpiringDomainsXLS?date=" + h.Date
	r, err := http.Get(url)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	f, err := os.Create(h.Date + ".xlsx")
	defer f.Close()
	if err != nil {
		panic(err)
	}

	io.Copy(f, r.Body)
}

func (h *Hunter) SetExcel() {
	en := "./" + h.Date + ".xlsx"
	var err error
	h.Excel, err = xlsx.OpenFile(en)
	if err != nil {
		Logger.Fatal("set excel error", zap.Error(err))
	}

	h.ExcelName = en
}

func (h *Hunter) SetOut() {
	// 若输出目录不存在，则创建目录
	if !tools.FileExist(Conf.Hunter.OutPath) {
		os.Mkdir(Conf.Hunter.OutPath, os.ModePerm)
	}

	path := Conf.Hunter.OutPath + "/" + h.Date + ".out"
	f, err := os.Create(path)
	if err != nil {
		Logger.Fatal("set out error", zap.Error(err))
	}

	h.Out = f
}

func (h *Hunter) hunter() {

	t := time.Now()
	// 下载指定日期的域名源文件
	h.DownLoad()
	Logger.Info("download time used", zap.Float64("eclapsed(s)", float64(time.Now().Sub(t).Nanoseconds())/1000000000))

	// 打开并初始化excel handler
	h.SetExcel()

	// 打开并初始化out文件
	h.SetOut()

	// 启动字符分类
	writerStart()

	ts := time.Now()
	// 开始根据指定规则过滤所有域名
	h.Filter()

	// 等待子线程完成
	time.Sleep(time.Millisecond)
	wg.Wait()

	// 数据按照一定的标准写入到文件中
	write(h.Out)

	Logger.Info("filter time used", zap.Float64("eclapsed(ms)", float64(time.Now().Sub(ts).Nanoseconds())/1000000))

	h.Out.Close()

	// 移除下载的excel文件
	os.Remove(h.ExcelName)
}
