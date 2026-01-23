package util

type Set struct {
	m map[string]bool
}

func NewSet() *Set {
	return &Set{m: make(map[string]bool)}
}

func (s *Set) Add(value string) {
	s.m[value] = true
}

func (s *Set) Remove(value string) {
	delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	_, exists := s.m[value]
	return exists
}

func (s *Set) Size() int {
	return len(s.m)
}

// because vals are bools, we only have function to get keys
func (s *Set) GetKeys() []string {
	keys := make([]string, 0, len(s.m))
	for k, v := range s.m {
		if v {
			keys = append(keys, k)
		}
	}
	return keys
}

func (s *Set) AddVals(value ...string) {
	for _, v := range value {
		s.m[v] = true
	}
}

func (s *Set) RemoveVals(values ...string) {
	for _, v := range values {
		delete(s.m, v)
	}
}
