package atexit

var functions []func()

// AtExit , call in main func ...
func AtExit() {
	for _, f := range functions {
		f()
	}
}

// Add functions to call on program exit
func Add(y ...func()) {
	functions = append(functions, y...)
}
