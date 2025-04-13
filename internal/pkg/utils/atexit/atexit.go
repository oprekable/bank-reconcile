package atexit

var functions []func()

// AtExit , call in main func ...
func AtExit() {
	for i := len(functions) - 1; i >= 0; i-- {
		functions[i]()
	}
}

// Add functions to call on program exit
func Add(y ...func()) {
	functions = append(functions, y...)
}
