package biz

var _ Facade = (*NoteFacade)(nil)

type NoteFacade struct {
	TData
}

func (c *NoteFacade) New() Facade {
	return &UserFacade{}
}

func (c *NoteFacade) GetTData() TData {
	return c.TData
}

func (c *NoteFacade) SetTData(tData TData) {
	c.TData = tData
}
