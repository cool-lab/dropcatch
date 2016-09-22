package hunter

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/tealeg/xlsx"
	"github.com/uber-go/zap"
)

var mux = &sync.Mutex{}

func (h *Hunter) Filter() {
	for _, sheet := range h.Excel.Sheets {
		Logger.Info("the total number of domains to be filter", zap.Int("domains", len(sheet.Rows)))
		cpuNum := runtime.NumCPU()
		preLen := len(sheet.Rows)/cpuNum - 1
		lastLen := len(sheet.Rows) % cpuNum
		for i := 0; i < cpuNum; i++ {
			if i == cpuNum-1 {
				go sheetRange(sheet, h.Out, preLen*i, lastLen)
			} else {
				go sheetRange(sheet, h.Out, preLen*i, preLen)
			}
		}
	}
}

func sheetRange(sheet *xlsx.Sheet, out *os.File, start int, l int) {
	wg.Add(1)
	for _, row := range sheet.Rows[start : start+l] {
		for _, cell := range row.Cells {
			s, _ := cell.String()
			s = strings.ToLower(s)

			pass, suffixLen := BaseFilter(s)
			if !pass {
				continue
			}

			pass, n := AdvanceFilter(s, suffixLen)
			if !pass {
				continue
			}

			// 输出到outs中
			s1 := strconv.Itoa(n)
			mux.Lock()
			outs[s1] = append(outs[s1], s)
			mux.Unlock()
		}
	}
	wg.Done()
}
