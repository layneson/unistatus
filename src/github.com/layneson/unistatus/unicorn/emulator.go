package unicorn

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"time"
)

const packetSize = 50

//Emulator represents a Provider which controls an external Unicorn emulator.
type Emulator struct {
	//Port is the local machine port on which the emulator is listening.
	Port int
}

//dial creates a new connection to the Emulator.
func (e Emulator) dial() (net.Conn, error) {
	return net.Dial("udp", ":"+strconv.Itoa(e.Port))
}

//Init implements the method from the Provider interface.
func (e Emulator) Init() error {
	conn, err := e.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	bs := make([]byte, packetSize)
	bs[0] = 0

	_, err = conn.Write(bs)
	if err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.Read(bs)
	if err != nil {
		return err
	}

	return nil
}

//Deinit implements the method from the Provider interface.
func (e Emulator) Deinit() error {
	return nil
}

//SetBrightness implements the method from the Provider interface.
func (e Emulator) SetBrightness(br float32) error {
	bint := int32(br * 255)

	conn, err := e.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	buff := new(bytes.Buffer)

	err = buff.WriteByte(1)
	if err != nil {
		return err
	}

	err = binary.Write(buff, binary.BigEndian, bint)
	if err != nil {
		return err
	}

	_, err = conn.Write(buff.Bytes())
	if err != nil {
		return err
	}

	return nil
}

//SetPixel implements the method from the Provider interface.
func (e Emulator) SetPixel(x, y, r, g, b int) error {
	ints := []int{x, y, r, g, b}

	conn, err := e.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	buff := new(bytes.Buffer)

	err = buff.WriteByte(2)
	if err != nil {
		return err
	}

	for _, i := range ints {
		err = binary.Write(buff, binary.BigEndian, int32(i))
		if err != nil {
			return err
		}
	}

	_, err = conn.Write(buff.Bytes())
	if err != nil {
		return err
	}

	return nil
}

//Show implements the method from the Provider interface.
func (e Emulator) Show() error {
	conn, err := e.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	bs := make([]byte, packetSize)
	bs[0] = 3

	_, err = conn.Write(bs)
	if err != nil {
		return err
	}

	return nil
}
