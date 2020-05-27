package models

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
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

//Equalize :均衡化
func Equalize(picture *Picture) (source []P, result []P) {

	source = make([]P, 0) //原始图像的灰度信息
	result = make([]P, 0) //均衡化后图像的灰度信息

	//遍历图像的每一个点
	for _, p := range picture.Points {

		//将RGB值转化为HSI值
		_, _, i := utils.RGB2HSI(p.R, p.G, p.B)
		source = append(source, P{X: p.X, Y: p.Y, Gray: int(i)})
	}

	return

}

//MakePicture :制作新的图像
func MakePicture(ps []P, bounds image.Rectangle, ext string) {

	newColor := image.NewRGBA(bounds)

	for _, p := range ps {
		gray := uint8(p.Gray)
		c := color.RGBA{R: gray, G: gray, B: gray}
		newColor.Set(p.X, p.Y, c)
	}
	//生成新的文件名
	uuid := uuid.NewV4()
	fileName := fmt.Sprintf("%s%s", uuid, ext)
	filePath := fmt.Sprintf("static/result/%s", fileName)

	fg, err := os.Create(filePath)
	check(err)
	defer fg.Close()
	err = jpeg.Encode(fg, newColor, nil)

}
