package detect

// DetectorList is a sortable slice of Detectors, and implements
// sort.Interface.
type DetectorList []*Detector

func (l DetectorList) Len() int {
	return len(l)
}

func (l DetectorList) Less(i, j int) bool {
	// Even though this fucntion is "less", we sort by highest priority first.
	return l[i].Priority > l[j].Priority
}

func (l DetectorList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
