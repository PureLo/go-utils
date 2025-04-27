package shp

import (
	"fmt"
	"math"
	"testing"
)

func TestTransform(t *testing.T) {
	transformer, err := computeAffineLeastSquares(mappings)
	if err != nil {
		t.Fatalf("计算仿射矩阵参数失败: %v", err)
	}

	for _, mapping := range mappings {
		lon, lat := transformer.Transform(mapping.X, mapping.Y)
		fmt.Printf("Input XY: (%.2f, %.2f) → LonLat: (%.6f, %.6f)\n", mapping.X, mapping.Y, lon, lat)
		fmt.Printf("Real LonLat: (%.2f, %.2f)\n", mapping.Lon, mapping.Lat)
		fmt.Printf("Distance:%.2f m\n", CalcDistance(Point{mapping.Lon, mapping.Lat}, Point{lon, lat}))
		fmt.Println("----------------------------------------")
	}
}

type Point struct {
	X float64
	Y float64
}

func CalcDistance(p1, p2 Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}
