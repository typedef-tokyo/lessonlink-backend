package vo

func SetVOConstructor[T, V any](t *T, f func(V) (T, error), _s V) error {

	v, err := f(_s)
	if err != nil {
		return err
	}

	*t = v
	return nil
}
