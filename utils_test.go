package yadi

func ResetYadi() {
	clearDeferredUpdates()
	err := closeContextSoft()
	if err != nil {
		panic(err)
	}
}
