package progress

type Writer interface {
	Start(l int64, name string) error
	Inc(n int64) error
	Stop() error
}
