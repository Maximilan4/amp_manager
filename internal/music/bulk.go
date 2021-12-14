package music

type AmidBulk struct {
    MaxItems int
    Values   []string
}

func (ab *AmidBulk) Add(id string) {
    ab.Values = append(ab.Values, id)
}

func (ab *AmidBulk) Flush() []string {
    values := ab.Values
    ab.Values = make([]string, 0, ab.MaxItems)
    return values
}

func (ab *AmidBulk) CanAdd() bool {
    return ab.MaxItems > len(ab.Values)
}
