package main

// cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
// compilazione GOOS=js GOARCH=wasm go build -o wasm.wasm
//https://www.kirsle.net/wiki/Go-WebAssembly
// attenzione : chrome mette in cache i file wasm, prima di eseguire il codice
// premere CTRL+SHIFT+R

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"syscall/js"
)

func funzione() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		//////////////////////////////////////////////////////////////////////////////
		var SCREEN_WIDTH int = 400
		var SCREEN_HEIGHT int = 400

		var re_min float32 = -2.0
		var im_min float32 = -1.2
		var re_max float32 = 1.0
		var im_max float32 = 1.2

		var iterazioni int = 1024

		var a, b float32
		var x, y, x_new, y_new, somma float32
		var k, i, j int

		var re_factor float32 = (re_max - re_min)
		var im_factor float32 = (im_max - im_min)

		m := image.NewRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))

		for i = 0; i < SCREEN_HEIGHT; i++ {
			for j = 0; j < SCREEN_WIDTH; j++ {
				a = re_min + (float32(j) * re_factor / float32(SCREEN_WIDTH))
				b = im_min + (float32(i) * im_factor / float32(SCREEN_HEIGHT))
				x = 0
				y = 0

				for k = 0; k < iterazioni; k++ {
					x_new = (float32(x) * float32(x)) - (float32(y) * float32(y)) + float32(a)
					y_new = (float32(2) * float32(x) * float32(y)) + float32(b)
					somma = (x_new * x_new) + (y_new * y_new)
					if somma > 4 {
						if k%2 == 0 {
							m.Set(j, i, color.RGBA{0, 0, 0, 255})
						} else {
							m.Set(j, i, color.RGBA{255, 255, 255, 255})
						}
						break
					}
					x = x_new
					y = y_new
				}
			}
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, m); err != nil {
			return ("unable to encode png")
		}
		readBuf, _ := ioutil.ReadAll(buf)

		enc := base64.StdEncoding.EncodeToString([]byte(readBuf))

		return "data:image/png;base64," + enc
	})
	return jsonFunc
}

func main() {
	fmt.Println("Inizio")
	js.Global().Set("go_function", funzione())
	fmt.Println(funzione())
	<-make(chan bool)
}
