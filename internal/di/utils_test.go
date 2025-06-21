package di

func ResetYadi() {
	ClearDeferredUpdates()
	err := CloseContextSoft()
	if err != nil {
		panic(err)
	}
}
