package repo

func MetaEqual(m, equals map[string]any) bool {
	for key, val := range equals {
		if v, ok := m[key]; !ok {
			return false
		} else {
			if v != val {
				return false
			}
		}
	}
	return true
}
