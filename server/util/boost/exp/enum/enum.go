package enum

type FlagEnum uint

func (r FlagEnum) Add(flags ...FlagEnum) FlagEnum {
	f := r
	for _, v := range flags {
		f |= v
	}
	return f
}

func (r FlagEnum) Remove(flags ...FlagEnum) FlagEnum {
	f := r
	for _, v := range flags {
		f &^= v
	}
	return f
}

func (r FlagEnum) Exist(flags ...FlagEnum) bool {
	for _, v := range flags {
		if r&v != v {
			return false
		}
	}
	return true
}

func (r FlagEnum) ExistAny(flags ...FlagEnum) bool {
	for _, v := range flags {
		if r&v == v {
			return true
		}
	}
	return false
}
