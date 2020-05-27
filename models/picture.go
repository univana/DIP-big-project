package models

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path"
)

//Point :点模型
type Point struct {
	X int //横坐标
	Y int //纵坐标
	R uint8
	G uint8
	B uint8 //RGB 值

}

//Picture :图片模型
type Picture struct {
	Path   string          //图片的存放路径
	Bounds image.Rectangle //范围
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
		result.Bounds, result.Points = analyseJPG(filePath)
	case ".png":
		fmt.Println("是PNG")
	}

	return
}

func analyseJPG(filePath string) (bounds image.Rectangle, points []Point) {

	//读取图片
	reader, err := os.Open(filePath)
	check(err)
	defer reader.Close()
	img, _ := jpeg.Decode(reader)

	bounds = img.Bounds()

	//读取每个点的RGB值
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			points = append(points, Point{X: x, Y: y, R: originalColor.R, G: originalColor.G, B: originalColor.B})
		}
	}

	return
}
