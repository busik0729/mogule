package helpers

import (
	"github.com/tsaikd/inject"
	"reflect"
)

var DI DependencyInjection

type DependencyInjection struct {
	DIInjector inject.Injector
}

func GetDI() DependencyInjection {
	return DI
}

func (di DependencyInjection) ProvideDI(injectObject interface{}) DependencyInjection {
	DI.DIInjector.Provide(injectObject)
	return DI
}

func (di DependencyInjection) Map(injectObject interface{}) DependencyInjection {
	DI.DIInjector.Map(injectObject)
	return DI
}

func (di DependencyInjection) Invoke(injectObject interface{}) []reflect.Value {

	r, _ := DI.DIInjector.Invoke(injectObject)

	return r
}

func Initialize() DependencyInjection {
	di := DependencyInjection{DIInjector: inject.New()}
	DI = di
	return DI
}
