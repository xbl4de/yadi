# yadi
YADI — Yet Another DI([Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection)) library for go

![image](https://github.com/user-attachments/assets/d3730f68-780d-49f0-b28b-fdad2797879a)


# Disclaimer

>Library in beta stage and documentation in development, not all aspects are described. Please wait for updates.

## Use context

Before access to any bean or value, you should enable DI context in your application:

```go
package main

import (
	"github.com/xbl4de/yadi"
)

func main() {
    yadi.UseLazyContext()
	// access to bean
}
```

You can provide bean before the call `yadi.UseLazyContext()` - all defined providers will be passed to the context. But you cannot to access to beans or values without this call.

## Provide a bean

YADI allows providing a way to build your structures:

```go
package main

import (
	"fmt"
	"github.com/xbl4de/yadi"
)

type ServiceE struct {
	Description string
}

var _ = yadi.SetBeanProvider[*ServiceE](func(ctx yadi.Context) (*ServiceE, error) {
	return &ServiceE{
		Description: "service E",
	}, nil
})

func main() {
    yadi.UseLazyContext()
	UseServiceE()
}

func UseServiceE() {
	serviceE, err := yadi.GetBean[*ServiceE]()
	if err != nil {
		// process error
	} else {
		fmt.Println(serviceE.Value.Description)
	}
}
```

If you have constructor functions to build your bean, you can pass them to YADI:

```go
package main

import (
	"github.com/xbl4de/yadi"
)

type ServiceE struct {
	Description string
}

func NewServiceE() *ServiceE {
	return &ServiceE{
		Description: "some description",
    }
}
func NewServiceEV2() (*ServiceE, error) {
	return &ServiceE{
		Description: "some description",
    }, nil
}

func main() {
	yadi.UseLazyContext()
	UseServiceE()
}

var _ = yadi.SetBeanProviderFunc[*ServiceE](NewServiceE)
// or
//var _ = yadi.SetBeanProviderFunc[*ServiceE](NewServiceEV2)

 func UseServiceE() {
	 ///
 }
```

> Important: YADI supports provider functions only with one or two return values. In the first case, return value should be exact bean. In the second one, the first value still should be a bean, the second value should be always `error`.

## Require a bean

There are available function `yadi.RequireBean[T]()`:

```go
func UseServiceE() {
	serviceE := yadi.RequireBean[*ServiceE]() // panics if cannot provide the bean
}
```

# Values

YADI supports key-value storage for any type. You can store values by string path and then access them at any part of your application:

```go
package main

import (
	"github.com/xbl4de/yadi"
)

var _ = yadi.SetValue[int]("some.int", 33)

func main() {
	yadi.UseLazyContext()
	SomeFunc()
}

func SomeFunc() {
	someInt, err := yadi.GetValue[int]("some.int")
	if err != nil {
		// process error
	} else {
		println(someInt.Value)
    }
}
```

## Provide values to builder functions

You can tell YDI where to find func args:

```go
package main

import (
	"github.com/xbl4de/yadi"
)

type ServiceE struct {
	Description string
}

func NewServiceE(desc string) *ServiceE {
	return &ServiceE{
		Description: desc,
    }
}

var _ = yadi.SetBeanProviderFunc[*ServiceE](NewServiceE, yadi.WithDefaultValueAt(0, "abcd"))
// or
var _ = yadi.SetBeanProviderFunc[*ServiceE](NewServiceE, 
	yadi.WithValuePathAt(0, "serviceE.description"))
```

## YADI tag

YADI provides a Go tag with folloving structure:

```
`yadi:"ignore;path=path.to.value"`
```

### Ignore

Marks field as ignored. YADI will NOT try to inject it. Incompatible with other parameters

### Path

Value a path. YADI will look for this path when does injection. You should provide the value by this path, otherwise yadi raises error.

## Guess the bean

If you don't provide a way to build the bean, YADI will try to create the provider by itself. YADI will inject all structure and interface fields, and will set all values if their paths were provided.

```go
package main

import (
	"github.com/xbl4de/yadi"
)

type ServiceE struct {
	Description string `yadi:"path=serviceE.description"`
}

var _ = yadi.SetValue[string]("serviceE.description", "abcd") // important: YADI cannot guess the values


func main() {
	yadi.UseLazyContext()
	UseServiceE()
}

func UseServiceE() {
	serviceE, err := yadi.GetBean[*ServiceE]() // YADI guess the provider
	if err != nil {
		// process error
	} else {
		// ...
	}
}

```

## Inject to value

YADI can do injection to provided value:

```go
package main

import (
	"github.com/xbl4de/yadi"
)

type ServiceE struct {
	Description string `yadi:"path=serviceE.description"`
}

type ServiceT struct {
	e *ServiceE
}
func main() {
	yadi.UseLazyContext()
	SomeFunc()
}


func SomeFunc() {
	t := &ServiceT{}
	err := yadi.Inject(&t)
	if err != nil {
		// process the error
    }
}

```

## Lazy beans

YADI supports lazy bean injection, but with some limitations:

```go
package main

import (
	"fmt"
	"github.com/xbl4de/yadi"
)

type ServiceE struct {
	Description string `yadi:"path=serviceE.description"`
}

func main() {
	yadi.UseLazyContext()
	SomeFunc()
}

func SomeFunc() {
	var e yadi.LazyBean[*ServiceE]
	e = yadi.InjectLazyBean[*ServiceE]()

	fmt.Println(e().Description)
}
```

The main limitation is injection: YADI cannot inject LazyBean by itself: you should manually to call ` yadi.InjectLazyBean` to get a function object.
