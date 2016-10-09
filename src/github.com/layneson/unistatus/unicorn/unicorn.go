package unicorn

//A Provider represents a Unicorn display, whether it is an emulation or the real thing.
type Provider interface {
	//Init initializes the Provider.
	Init() error

	//Deinit deinitializes (destroys) the Provider.
	Deinit() error

	//SetBrightness sets the display brightness, between 0 and 1.
	SetBrightness(b float32) error

	//SetPixel sets the color of the specified pixel.
	SetPixel(x, y, r, g, b int) error

	//Show displays the changes made to the buffer.
	Show() error
}

//provider is the package Provider.
var provider Provider

//InitProvider initializes and sets the package provider as long as no errors were thrown during initialization.
func InitProvider(p Provider) error {
	err := p.Init()
	if err != nil {
		return err
	}

	provider = p
	return nil
}

//Deinit deinitializes the package provider.
func Deinit() error {
	return provider.Deinit()
}

//SetBrightness sets the brightness for the package provider.
func SetBrightness(b float32) error {
	return provider.SetBrightness(b)
}

//SetPixel sets the color of the specified pixel for the package provider.
func SetPixel(x, y, r, g, b int) error {
	return provider.SetPixel(x, y, r, g, b)
}

//Show displays the changes made to the buffer for the package provider.
func Show() error {
	return provider.Show()
}
