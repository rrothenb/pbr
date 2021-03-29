package main

import (
	"fmt"
	"os"

	"github.com/rrothenb/pbr/pkg/camera"
	"github.com/rrothenb/pbr/pkg/geom"
)

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
}

func printInfo(b *geom.Bounds, surfaces int, c *camera.SLR) {
	fmt.Println("Min:", b.Min)
	fmt.Println("Max:", b.Max)
	fmt.Println("Center:", b.Center)
	fmt.Println("Camera:", c)
}
