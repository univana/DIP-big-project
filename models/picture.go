package models

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"myApp/utils"
	"os"
	"path"

	uuid "github.com/satori/go.uuid"
)

//Point :点模型
type Point struct {
	X int //横坐标
	Y int //纵坐标
	R uint8
	G uint8
	B uint8 //RGB 值

}

//P :包含灰度信息的点模型
type P struct {
	X    int //横坐标
	Y    int //纵坐标
	Gray int //灰度
}

//Picture :图片模型
type Picture struct {
	Path   string          //图片的存放路径
	Bounds image.Rectangle //范围
	PNum   int             //像素数目
	Ext    string          //图片格式
	Points []Point         //点集
}

//check :错误检测
func check(err error) {
	if err != nil {
		panic(err)
	}
}

//NewPicture :新建图片实例
func NewPicture() *Picture {
	return &Picture{}
}

//Analyse :解析图片
func (m *Picture) Analyse(filePath string) (result *Picture) {

	result = NewPicture()

	result.Path = fmt.Sprintf("/%s", filePath)

	//确定图片格式
	result.Ext = path.Ext(filePath)

	switch result.Ext {
	case ".jpg":
		result.Bounds, result.Points, result.PNum = analyseJPG(filePath)
	case ".png":
		result.Bounds, result.Points, result.PNum = analysePNG(filePath)
	}

	return
}

func analyseJPG(filePath string) (bounds image.Rectangle, points []Point, pNum int) {

	//读取图片
	reader, err := os.Open(filePath)
	check(err)
	defer reader.Close()
	img, _ := jpeg.Decode(reader)

	bounds = img.Bounds()
	pNum = 0

	//读取每个点的RGB值
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			points = append(points, Point{X: x, Y: y, R: originalColor.R, G: originalColor.G, B: originalColor.B})
			pNum++
		}
	}

	return
}

func analysePNG(filePath string) (bounds image.Rectangle, points []Point, pNum int) {

	//读取图片
	reader, err := os.Open(filePath)
	check(err)
	defer reader.Close()

	img, _ := png.Decode(reader)

	bounds = img.Bounds()
	pNum = 0

	//读取每个点的RGB值
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			points = append(points, Point{X: x, Y: y, R: originalColor.R, G: originalColor.G, B: originalColor.B})
			pNum++
		}
	}

	return

}

//Equalize :均衡化
func Equalize(picture *Picture) ([256]float64, [256]float64, *Picture, [256]int) {

	var (
		minGray int //原始图像的最小灰度值
		maxGray int //原始图像的最大灰度值

		cdf     [256]float64 //原始图像灰度分布函数
		mapping [256]int     //原始图像灰度值与结果灰度值的映射表

		oriGrayNums      [256]int     //原始图像灰度值统计数目
		oriHistogramData [256]float64 //原始图像灰度直方图数据
		resHistogramData [256]float64 //结果图像灰度直方图数据

		resPicture *Picture //结果图像对象，用于生成结果图像
	)

	/* 遍历原始图像的每一个点 */
	for _, p := range picture.Points {

		/* 将RGB值转化为HSI值 */
		_, _, i := utils.RGB2HSI(p.R, p.G, p.B)

		gray := int(i)

		/* 确定灰度值范围 */
		if gray > maxGray {
			maxGray = gray
		}
		if gray < minGray {
			minGray = gray
		}

		/* 统计灰度值数目 */
		oriGrayNums[gray]++
	}

	for i := 0; i < 256; i++ {

		/* 计算原始图像灰度直方图数据 */
		oriHistogramData[i] = float64(oriGrayNums[i]) / float64(picture.PNum)

		/* 计算累积分布函数 */
		if i > 0 {
			cdf[i] = cdf[i-1] + oriHistogramData[i]
		} else {
			cdf[0] = oriHistogramData[0]
		}

		/* 计算映射函数 */
		mapping[i] = int(math.Floor(float64(maxGray-minGray)*cdf[i] + float64(minGray) + 0.5))

		/* 计算均衡化后灰度直方图数据 */
		resHistogramData[mapping[i]] = oriHistogramData[i]

	}

	/* 生成结果图像 */
	resPicture = NewPicture()
	resPicture.Bounds = picture.Bounds
	resPicture.PNum = picture.PNum
	resPicture.Points = make([]Point, 0)
	resPicture.Ext = picture.Ext

	for _, p := range picture.Points {

		/* 将RGB值转化为HSI值 */
		h, s, i := utils.RGB2HSI(p.R, p.G, p.B)

		/* 将 I 映射为均衡化结果 */
		i = float64(mapping[int(i)])

		/* 将HSI值转回RGB值 */
		r, g, b := utils.HSI2RGB(h, s, i)

		resPicture.Points = append(resPicture.Points, Point{X: p.X, Y: p.Y, R: r, G: g, B: b})

	}

	return oriHistogramData, resHistogramData, resPicture, mapping

}

//MakePicture :制作新的图像
func MakePicture(picture *Picture) string {

	if picture.Ext == ".jpg" {

		newColor := image.NewRGBA(picture.Bounds)

		for _, p := range picture.Points {

			c := color.RGBA{R: p.R, G: p.G, B: p.B}
			newColor.Set(p.X, p.Y, c)
		}
		//生成新的文件名
		uuid := uuid.NewV4()
		fileName := fmt.Sprintf("%s%s", uuid, picture.Ext)
		filePath := fmt.Sprintf("static/result/%s", fileName)

		fg, err := os.Create(filePath)
		check(err)
		defer fg.Close()
		err = jpeg.Encode(fg, newColor, nil)
		return fmt.Sprintf("/%s", filePath)
	} else {

		newColor := image.NewNRGBA(picture.Bounds)

		for _, p := range picture.Points {

			c := color.NRGBA{R: p.R, G: p.G, B: p.B, A: 255}
			newColor.Set(p.X, p.Y, c)
		}

		//生成新的文件名
		uuid := uuid.NewV4()
		fileName := fmt.Sprintf("%s%s", uuid, picture.Ext)
		filePath := fmt.Sprintf("static/result/%s", fileName)

		fg, err := os.Create(filePath)
		check(err)
		defer fg.Close()

		err = png.Encode(fg, newColor)
		return fmt.Sprintf("/%s", filePath)
	}

}

//Specificate :规定化
func Specificate(oriPicture *Picture, oriHistogramData [256]float64, oriMapping [256]int, matchMapping [256]int) ([256]float64, *Picture) {

	var (
		resMapping       [256]int     //匹配后的映射函数
		resHistogramData [256]float64 //结果图像灰度直方图数据
		resPicture       *Picture     //结果图像对象
	)

	/* 遍历每一个灰度级 */
	for i := 0; i < 256; i++ {

		oriMappingGray := oriMapping[i] //原始图像均衡化后的灰度值

		var offset float64 = 1000 //两个图像映射结果的差距
		var flag int              //目标图像最匹配灰度值

		/* 遍历每个目标灰度值 */
		for j := 0; j < 256; j++ {

			matchMappingGray := matchMapping[j] //目标图像均衡化后的灰度值

			/* 计算偏移量 */
			temp := math.Abs(float64(oriMappingGray - matchMappingGray))

			if temp < offset {
				offset = temp
				flag = j
			}

			resMapping[i] = flag
		}

	}

	/* 计算结果直方图数据 */
	for i := 0; i < 256; i++ {
		resHistogramData[resMapping[i]] = oriHistogramData[i]
	}

	/* 生成结果图像对象 */
	resPicture = NewPicture()
	resPicture.Bounds = oriPicture.Bounds
	resPicture.Ext = oriPicture.Ext
	resPicture.PNum = oriPicture.PNum
	resPicture.Points = make([]Point, 0)

	for _, p := range oriPicture.Points {

		/* 将RGB值转化为HSI值 */
		h, s, i := utils.RGB2HSI(p.R, p.G, p.B)

		/* 将 I 映射为均衡化结果 */
		i = float64(resMapping[int(i)])

		/* 将HSI值转回RGB值 */
		r, g, b := utils.HSI2RGB(h, s, i)

		resPicture.Points = append(resPicture.Points, Point{X: p.X, Y: p.Y, R: r, G: g, B: b})
	}

	return resHistogramData, resPicture

}
