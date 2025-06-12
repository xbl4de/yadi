package yadi

type ServiceIgnore struct {
	ServiceA *ServiceA `yadi:"ignore"`
}

type ServiceBagTag struct {
	ServiceA *ServiceA `yadi:"tag"`
}

type ServiceClose struct {
	Closed bool
}

func NewServiceClose() *ServiceClose {
	return &ServiceClose{
		Closed: false,
	}
}

func (s *ServiceClose) Close() error {
	s.Closed = true
	return nil
}

type ServiceE struct {
	Description string `yadi:"path=serviceE.description"`
}

type CountInterface interface {
	GetCount() int
}

type ServiceF struct {
	Count int `yadi:"path=serviceF.count"`
}

func (f *ServiceF) GetCount() int {
	return f.Count
}

type ServiceG struct {
	Enabled bool `yadi:"path=serviceG.enabled"`
}

type ServiceH struct {
	Timeout int `yadi:"path=serviceH.timeout"`
}

type ServiceA struct {
	Name     string `yadi:"path=serviceA.name"`
	ServiceE *ServiceE
	ServiceF *ServiceF
}

type ServiceB struct {
	Age      int `yadi:"path=serviceB.age"`
	ServiceF *ServiceF
	ServiceH *ServiceH
}

type ServiceC struct {
	Location string `yadi:"path=serviceC.location"`
	ServiceG *ServiceG
}

const (
	ServiceAName        = "A-Name"
	ServiceBAge         = 10
	ServiceCLocation    = "C-Location"
	ServiceEDescription = "E-Description"
	ServiceFCount       = 20
	ServiceGEnabled     = true
	ServiceHTimeout     = 30
)

func ProvideDefaultValues() {
	SetValue("serviceA.name", ServiceAName)
	SetValue("serviceB.age", ServiceBAge)
	SetValue("serviceC.location", ServiceCLocation)
	SetValue("serviceE.description", ServiceEDescription)
	SetValue("serviceF.count", ServiceFCount)
	SetValue("serviceH.timeout", ServiceHTimeout)
	SetValue("serviceG.enabled", ServiceGEnabled)
}

func NewServiceE(description string) *ServiceE {
	return &ServiceE{Description: description}
}

func NewServiceF(count int) *ServiceF {
	return &ServiceF{Count: count}
}

func NewServiceG(enabled bool) *ServiceG {
	return &ServiceG{Enabled: enabled}
}

func NewServiceH(timeout int) *ServiceH {
	return &ServiceH{Timeout: timeout}
}

func NewServiceA(name string, e *ServiceE) *ServiceA {
	return &ServiceA{Name: name, ServiceE: e}
}

func NewServiceB(age int, f *ServiceF, h *ServiceH) *ServiceB {
	return &ServiceB{Age: age, ServiceF: f, ServiceH: h}
}

func NewServiceC(location string, g *ServiceG) *ServiceC {
	return &ServiceC{Location: location, ServiceG: g}
}

type MainService struct {
	ServiceA *ServiceA
	ServiceB *ServiceB
	ServiceC *ServiceC
}

func NewMainService(a *ServiceA, b *ServiceB, c *ServiceC) *MainService {
	return &MainService{
		ServiceA: a,
		ServiceB: b,
		ServiceC: c,
	}
}
