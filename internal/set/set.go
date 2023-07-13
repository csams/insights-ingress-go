package set

type Set[K comparable] map[K]bool

func New[K comparable](ks []K) Set[K] {
    var res Set[K]
    for _, k := range ks {
        res[k] = true
    }
    return res
}

func (s Set[K]) Contains(k K) bool {
    return s[k]
}
