package web

const RequestScope string = "request"

type AppRequestScope struct {
}

func NewAppRequestScope() AppRequestScope {
	return AppRequestScope{}
}

func (scope AppRequestScope) GetPeaObject(peaName string) interface{} {
	return nil
}

func (scope AppRequestScope) RemovePeaObject(peaName string) interface{} {
	return nil
}
