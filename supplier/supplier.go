package supplier

type Func func() string
type ArrayFunc func() []interface{}

func GetConstSupplier(value string) Func {
	return func() string {
		return value
	}
}

func GetArraySupplier(value []interface{}) ArrayFunc {
	return func() []interface{} {
		return value
	}
}
