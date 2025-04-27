package shp

import "fmt"

type PointMapping struct {
	X, Y     float64 // 本地坐标
	Lon, Lat float64 // 对应的经纬度
}

// 提供多组对应点
var mappings = []PointMapping{
	// 2261 POLYGON ((
	// 23549 48945, 左上
	// 28381 48945, 右上
	// 28381 46597, 右下
	// 23549 46597 左下))
	{23549, 46597, 11584406.852, 3574681.706}, // 2261 左下角
	{23549, 48945, 11584406.852, 3574684.106}, // 2261 左上角
	{28381, 48945, 11584412.352, 3574684.106}, // 2261 右上角

	// 2242 POLYGON ((
	// 77407 46622, 右下
	// 72110 46622, 左下
	// 72110 48964, 左上
	// 77407 48964 右上))
	{72110, 46622, 11584456.252, 3574681.706}, // 2242 左下角
	{72110, 48964, 11584456.252, 3574684.106}, // 2242 左上角
	{77407, 48964, 11584461.752, 3574684.106}, // 2242 右上角

	// 2312 POLYGON ((
	// 10548 93, 左下
	// 10548 5546, 左上
	// 12923 5546,  右上
	// 12923 93 右下))
	{10548, 93, 11584394.352, 3574634.806},   // 2312 左下角
	{10548, 5546, 11584394.352, 3574640.306}, // 2312 左上角
	{12923, 5546, 11584396.752, 3574640.306}, // 2312 右上角

	// 2332 POLYGON ((
	// 69075 105, 左下
	// 69075 5633, 左上
	// 71434 5633,  右上
	// 71434 105 右下))
	{69075, 5633, 11584453.302, 3574640.306}, // 2332 左上角
	{71434, 5633, 11584455.702, 3574640.306}, // 2332 右上角
	{71434, 105, 11584455.702, 3574634.806},  // 2332 右下角

	// 2279 POLYGON ((
	// 77388 11559, 右下
	// 72342 11559, 左下
	// 72342 13911, 左上
	// 77388 13911 右上))
	{72342, 11559, 11584456.252, 3574646.406}, // 2279 左下角
	{77388, 13911, 11584461.752, 3574648.806}, // 2279 右上角
	{77388, 11559, 11584461.752, 3574646.406}, // 2279 右下角
}

// 仿射变换参数
type AffineTransform struct {
	A, B, C float64
	D, E, F float64
}

// 使用最小二乘法从多组点拟合仿射变换
func computeAffineLeastSquares(mappings []PointMapping) (*AffineTransform, error) {
	n := len(mappings)
	if n < 3 {
		return nil, fmt.Errorf("至少需要3个点进行仿射拟合")
	}

	// 构建矩阵求解
	var sumX, sumY, sum1 float64
	var sumXX, sumXY, sumYY float64
	var sum1Lon, sum1Lat float64
	var sumXLon, sumYLon float64
	var sumXLat, sumYLat float64

	for _, m := range mappings {
		x, y := m.X, m.Y
		lon, lat := m.Lon, m.Lat

		sumX += x
		sumY += y
		sum1 += 1
		sumXX += x * x
		sumXY += x * y
		sumYY += y * y

		sumXLon += x * lon
		sumYLon += y * lon
		sum1Lon += lon

		sumXLat += x * lat
		sumYLat += y * lat
		sum1Lat += lat
	}

	// 构造正规方程组并求解
	// 系数矩阵
	M := [3][3]float64{
		{sumXX, sumXY, sumX},
		{sumXY, sumYY, sumY},
		{sumX, sumY, sum1},
	}

	// 右边向量
	bLon := [3]float64{sumXLon, sumYLon, sum1Lon}
	bLat := [3]float64{sumXLat, sumYLat, sum1Lat}

	// 解线性方程组
	A, B, C := solveLinearSystem(M, bLon)
	D, E, F := solveLinearSystem(M, bLat)

	return &AffineTransform{
		A: A, B: B, C: C,
		D: D, E: E, F: F,
	}, nil
}

// 用高斯消元法解3x3线性方程组
func solveLinearSystem(M [3][3]float64, b [3]float64) (x, y, z float64) {
	// 拷贝矩阵
	a := M
	v := b

	// 消元
	for i := range 3 {
		// 主元归一化
		pivot := a[i][i]
		for j := range 3 {
			a[i][j] /= pivot
		}
		v[i] /= pivot

		// 消去其他行
		for k := range 3 {
			if k != i {
				factor := a[k][i]
				for j := range 3 {
					a[k][j] -= factor * a[i][j]
				}
				v[k] -= factor * v[i]
			}
		}
	}

	x = v[0]
	y = v[1]
	z = v[2]
	return
}

// 进行仿射变换
func (t *AffineTransform) Transform(x, y float64) (lon, lat float64) {
	lon = t.A*x + t.B*y + t.C
	lat = t.D*x + t.E*y + t.F
	return
}
