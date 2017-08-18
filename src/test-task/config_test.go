package main

import "testing"

func TestShouldLoadConfig(t *testing.T) {
	var config Config;
	e := config.Load()
	if e != nil { t.Fatal(e.Error()) }
}
