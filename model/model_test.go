package model

import (
	"testing"
	"fmt"
)

func TestCalculateValues(t *testing.T) {
	fmt.Println(calculateValues(0, 355, 10))
	fmt.Println(calculateValues(-0.002, 0.03, 10))
	fmt.Println(calculateValues(-0.072, -0.01, 10))
}
