package utils

func PanicOn(err error) {
	if err != nil {
		panic(err)
	}
}
