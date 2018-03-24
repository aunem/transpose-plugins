package main

type HTTPListenerSpec struct {
	Port int  `json:"port" yaml:"port"`
	SSL  bool `json:"ssl" yaml:"ssl"`
}
