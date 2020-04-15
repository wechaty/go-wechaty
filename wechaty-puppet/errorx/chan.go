package errorx

type ChanErr struct {
  errChan chan error
}

func (c *ChanErr) Put(err error) {
  c.errChan <- err
}

func (c *ChanErr) WaitErr() <-chan error {
  return c.errChan
}

func NewChanErr(threadNum int) *ChanErr {
  return &ChanErr{errChan: make(chan error, threadNum)}
}
