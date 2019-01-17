package giterator

import "reflect"

type Iterator []interface{}

func Over(v interface{}) *Iterator {
	e := reflect.New(reflect.TypeOf(v)).Elem()
	e.Set(reflect.ValueOf(v))

	o := reflect.Indirect(e)

	switch o.Kind() {
	case reflect.Slice:
		i := make(Iterator, o.Len(), o.Cap())

		for index := 0; index < o.Len(); index++ {
			i[index] = o.Index(index).Interface()
		}

		return &i
	default:
		return &Iterator{}
	}
}

func (i *Iterator) Map(v interface{}) *Iterator {
	f := reflect.ValueOf(v)

	n := make(Iterator, 0, len(*i))
	for _, x := range *i {
		e := reflect.New(f.Type().In(0)).Elem()
		e.Set(reflect.ValueOf(x))

		n = append(n, (f.Call([]reflect.Value{e})[0]).Interface())
	}

	return &n
}

func (i *Iterator) Reduce(v1, v2 interface{}) interface{} {
	f := reflect.ValueOf(v1)

	a := reflect.New(f.Type().In(1)).Elem()
	a.Set(reflect.ValueOf(v2))

	for _, x := range *i {
		e := reflect.New(f.Type().In(0)).Elem()
		e.Set(reflect.ValueOf(x))

		a.Set(f.Call([]reflect.Value{e, a})[0])
	}

	return a.Interface()
}

func (i *Iterator) ForEach(v interface{}) *Iterator {
	f := reflect.ValueOf(v)

	for _, x := range *i {
		e := reflect.New(f.Type().In(0)).Elem()
		e.Set(reflect.ValueOf(x))

		f.Call([]reflect.Value{e})
	}

	return i
}

func (i *Iterator) FilterOn(v interface{}) *Iterator {
	f := reflect.ValueOf(v)

	n := make(Iterator, 0, len(*i))
	for _, x := range *i {
		e := reflect.New(f.Type().In(0)).Elem()
		e.Set(reflect.ValueOf(x))

		if (f.Call([]reflect.Value{e})[0]).Bool() {
			n = append(n, x)
		}
	}

	return &n
}
