package datasets

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/dfcf"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
)

// DataQuarterlyReport 季报
type DataQuarterlyReport struct {
	DataCache
	cache map[string]dfcf.QuarterlyReport
}

func init() {
	_ = cache.Register(&DataQuarterlyReport{})
}

func (r *DataQuarterlyReport) Usage() string {
	//TODO implement me
	panic("implement me")
}

func (r *DataQuarterlyReport) Print(code string, date ...string) {
	//TODO implement me
	panic("implement me")
}

func (r *DataQuarterlyReport) Kind() DataKind {
	return BaseQuarterlyReports
}

func (r *DataQuarterlyReport) Name() string {
	return mapDataSets[r.Kind()].Name
}

func (r *DataQuarterlyReport) Key() string {
	return mapDataSets[r.Kind()].Key
}

func (r *DataQuarterlyReport) Filename(date, code string) string {
	//TODO implement me
	panic("implement me")
}

func (r *DataQuarterlyReport) Init(barIndex *int, date string) error {
	*barIndex++
	r.cache = IntegrateQuarterlyReports(barIndex, date)
	return nil
}

func (r *DataQuarterlyReport) Update(cacheDate, featureDate string) {
	_ = cacheDate
	_ = featureDate
}

func (r *DataQuarterlyReport) Repair(cacheDate, featureDate string) {
	_ = cacheDate
	_ = featureDate
}

func (r *DataQuarterlyReport) Increase(snapshot quotes.Snapshot) {
	_ = snapshot
}

func (r *DataQuarterlyReport) Clone(date string, code string) DataSet {
	var dest = DataQuarterlyReport{DataCache: DataCache{Date: date, Code: code}}
	return &dest
}

// IntegrateQuarterlyReports 更新季报数据
func IntegrateQuarterlyReports(barIndex *int, date string) map[string]dfcf.QuarterlyReport {
	modName := "季报概要信息"
	logger.Info(modName + ", 任务开始启动...")

	allReports := []dfcf.QuarterlyReport{}
	reports, pages, _ := dfcf.QuarterlyReports(date)
	if pages < 1 || len(reports) == 0 {
		return nil
	}
	allReports = append(allReports, reports...)
	bar := progressbar.NewBar(*barIndex, "执行["+modName+"]", pages-1)
	for pageNo := 2; pageNo < pages+1; pageNo++ {
		bar.Add(1)
		list, pages, err := dfcf.QuarterlyReports(date, pageNo)
		if err != nil || pages < 1 {
			logger.Error(err)
			break
		}
		count := len(list)
		if count == 0 {
			break
		}
		allReports = append(allReports, list...)
		if count < dfcf.EastmoneyQuarterlyReportAllPageSize {
			break
		}
	}
	mapReports := map[string]dfcf.QuarterlyReport{}
	if len(allReports) > 0 {
		for _, v := range allReports {
			mapReports[v.SecurityCode] = v
		}
		filename := cache.ReportsFilename(date)
		err := api.SlicesToCsv(filename, allReports)
		if err != nil {
			logger.Errorf("cache %s failed, error: %+v", filename, err)
		}
	}
	//// 确定更新日期
	//currentDate := cachel5.DefaultCanUpdateDate()
	//c1d := flash.CacheF10()
	//merge := func(securityCode string, f10 *flash.F10) (ok bool) {
	//	f10.Date = currentDate
	//	cover, ok := mapReports[securityCode]
	//	if ok {
	//		f10.BPS = cover.BPS
	//		f10.BasicEPS = cover.BasicEPS
	//	}
	//	return ok
	//}
	//c1d.Checkout(currentDate)
	//c1d.Apply(merge)
	//logger.Info(modName+", 任务执行完毕.", time.Now())
	return mapReports
}
