package main

import (
	"fmt"
)

type osImageError struct {
	os     string
	images string
}
type notEnoughResources struct {
	cpu string
	ram string
	storage string
	cpuLeft string
	ramLeft string
	storageLeft string
}

func (r *osImageError) Error() string {
	return fmt.Sprintf("%s is not found. Please use one from %s", r.os, r.images)
}
func (r *notEnoughResources) Error() string {
	return fmt.Sprintf("You need %s MB ram, %s cpu, %s MB storage for this server. Only %s MB ram, %s cpu, %s MB storage left. You don't have enough resources to create this instance. Please delete other instances or buy new resources", r.ram, r.cpu, r.storage, r.ramLeft, r.cpuLeft, r.storageLeft)
}