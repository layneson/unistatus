package unicorn

/*
#include "hat.go.h"
*/
import "C"

//The HATProvider provides for a legitimate Unicorn HAT.
type HATProvider struct{}

//Init implements the method of the Provider interface.
func (d HATProvider) Init() error {
	C.init()

	return nil
}

//Deinit implements the method of the Provider interface.
func (d HATProvider) Deinit() error {
	C.clear()
	C.show()

	C.stop()

	return nil
}

//SetBrightness implements the method of the Provider interface.
func (d HATProvider) SetBrightness(b float32) error {
	C.setBrightness(C.int(int(b / 255.0)))

	return nil
}

//SetPixel implements the method of the Provider interface.
func (d HATProvider) SetPixel(x, y, r, g, b int) error {
	//Flip x-axis
	x = 7 - x

	rgb := (g << 16) | (r << 8) | b // please, don't ask...
	C.setPixel(C.int(x), C.int(y), C.uint32_t(rgb))

	return nil
}

//Show implements the method of the Provider interface.
func (d HATProvider) Show() error {
	C.show()

	return nil
}
