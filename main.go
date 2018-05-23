package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"io/ioutil"
	"bytes"
	"os"
	"image/jpeg"
)

func main() {
	img := kingpin.Flag("image", "image file").Short('i').Default("").String()
	html := kingpin.Flag("html", "html file").Short('h').Default("./index.html").String()
	font := kingpin.Flag("font", "the font you want").Short('f').Default("y").String()

	kingpin.Parse()
	if *img == "" {
		help()
		return
	}

	err := createImg(*img, *html, *font)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success")
	}

}

func createImg(img, html, font string) (err error) {

	file, err := ioutil.ReadFile(img)
	if err != nil {
		return
	}
	imageByte := bytes.NewBuffer(file)
	image, err := jpeg.Decode(imageByte)
	if err != nil {
		return
	}
	rectangle := image.Bounds()

	htmlFile, err := os.Create(html)
	defer htmlFile.Close()

	content := `<html>
<head>
    <meta charset="utf-8">
    <title> you </title>
    <style type="text/css">
        body {
            margin: 0px; 
			padding: 0px; 
			line-height:100%; 
			letter-spacing:0px; 
			text-align: center;
            width: auto !important;
            background-color: #ffffff;
        }
		font{
			font-size:1px;
		}
    </style>
</head>
<body>
<div>`

	fontRune := []rune(font)
	fontLen := len(fontRune)

	for i := rectangle.Min.Y; i < rectangle.Max.Y; i++ {
		fmt.Println(i*100/rectangle.Max.Y, "%")
		for j := rectangle.Min.X; j < rectangle.Max.X; j++ {
			color := image.At(j, i)
			r, g, b, _ := color.RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			rString := fmt.Sprint(r)
			gString := fmt.Sprint(g)
			bString := fmt.Sprint(b)
			num := j % fontLen
			fontNow := string(fontRune[num])
			content += `<font style='color:rgb(` + rString + `,` + gString + `,` + bString + `)'>` + fontNow + `</font>`
		}
		content += "<br>"
	}

	content += `</div>
			</body>
		</html>`

	htmlFile.Write([]byte(content))
	return
}

func help() {
	fmt.Println("输入错误，请按照下面的格式输入：")
	fmt.Println("使用: img2html [OPTION] [VALUE]")
	fmt.Println("  Options is flow:")
	fmt.Println("    -i         图片路径")
	fmt.Println("    -h         生成的html文件路径")
	fmt.Println("    -f         要取的字")
}
