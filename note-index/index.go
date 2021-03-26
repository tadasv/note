package main

type SetIndex struct {
	data map[string][]string
}

func (p *SetIndex) Add(key, value string) error {
	if items, ok := p.data[key]; ok {
		for _, val := range items {
			if val == value {
				return nil
			}
		}

		items = append(items, value)
		p.data[key] = items
	} else {
		p.data[key] = []string{value}
	}
	return nil
}
