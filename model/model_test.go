package model

import (
	"testing"
	"fmt"
)

func TestCalculateValues(t *testing.T) {
	fmt.Println(calculateValues(0, 355, 10))
	fmt.Println(calculateValues(8.609999656677246, 11.9399995803833, 6))
	fmt.Println(calculateValues(103, 355, 10))
	fmt.Println(calculateValues(-0.042, 0.03, 10))
	fmt.Println(calculateValues(-0.072, -0.01, 10))
}
