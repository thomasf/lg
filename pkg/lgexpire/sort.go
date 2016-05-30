package lgexpire

type logfilesByTime []logFile

func (p logfilesByTime) Len() int {
	return len(p)
}

func (p logfilesByTime) Less(i, j int) bool {
	return p[i].Time.Before(p[j].Time)
}

func (p logfilesByTime) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
