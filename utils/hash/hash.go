package hash

type HashSum string

func (h HashSum) String() string {
	return string(h)
}

type HashCode interface {
	String() string
	Hash() HashSum
}
