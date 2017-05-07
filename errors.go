package main

import (
		"fmt"
)
type osImageError struct{
	os string
	images string
}
func (r *osImageError) Error() string {
	return fmt.Sprintf("%s is not found. Please use one from %s",r.os, r.images)
}