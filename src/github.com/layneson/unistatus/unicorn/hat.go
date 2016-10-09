package unicorn

/*
#include "hat.go.h""
*/
import "C"

//The HATProvider provides for a legitimate Unicorn HAT.
type DefaultProvider struct{}

//Init implements the method of the Provider interface.
func (d DefaultProvider) Init() error {
	C.init()

	return nil
}

//Deinit implements the method of the Provider interface.
func (d DefaultProvider) Deinit() error {
	C.stop()

	return nil
}

//SetBrightness implements the method of the Provider interface.
func (d DefaultProvider) SetBrightness(b float32) error {
	C.setBrightness(C.int(int(b / 255.0)))

	return nil
}

//SetPixel implements the method of the Provider interface.
func (d DefaultProvider) SetPixel(x, y, r, g, b int) error {
	rgb := (g << 16) | (r << 8) | b // please, don't ask...
	C.setPixel(C.int(x), C.int(y), C.uint32_t(rgb))

	return nil
}

//Show implements the method of the Provider interface.
func (d DefaultProvider) Show() error {
	C.show()

	return nil
}
