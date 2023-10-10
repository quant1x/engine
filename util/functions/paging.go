package functions

// GetPages 计算页数
func GetPages(pageSize, count int) (pages int) {
	//pages = int(math.Ceil(float64(raw.Data.TotalHits) / float64(EastmoneyNoticesPageSize)))
	pages = count / pageSize
	n := count % pageSize
	if n > 0 {
		pages++
	}
	return pages
}
