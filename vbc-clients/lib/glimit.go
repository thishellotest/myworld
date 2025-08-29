package lib

type GLimit struct {
	n int
	c chan struct{}
}

// NewGLimit n:最大并行处理数量
func NewGLimit(n int) *GLimit {
	return &GLimit{
		n: n,
		c: make(chan struct{}, n),
	}
}

// Run f in a new goroutine but with limit.
func (g *GLimit) Run(f func()) {
	g.c <- struct{}{}
	go func() {
		f()
		<-g.c
	}()
}
