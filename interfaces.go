package ipfilter

type Filter interface {
	Contains(string) bool
}

type WhiteList interface {
	Contains(string) bool
	Remove(string) bool
}
