package main

type HTTPListenerSpec struct {
	Port string `json:"port" yaml:"port"`
	SSL  bool   `json:"ssl" yaml:"ssl"`
}
