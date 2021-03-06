# unistatus
Unistatus is a program designed to run on an Raspberry PI. It uses the Pimoroni Unicorn HAT to cycle between different animated status icons which can display information about the current weather, bus delays, etc.

Since the Unicorn HAT only supports an 8x8 pixel grid, it is not possible to display a large amount of information. However, each pixel is an RGB LED, which means that many colors are able to be displayed. Thus, with clever use of animation and color, more can be conveyed than one might initially think.

## Building
1. Clone the repository onto a Raspberry Pi.
2. Install go (https://golang.org) and gb (https://getgb.io).
3. Run `make` from inside the `unistatus` directory.
4. Run `sudo make install` from inside the `unistatus` directory.
5. Create the appropriate `credentials.json` file within the `/opt/unistatus` folder.

The final binary will be located in the `/opt/unistatus` folder. To uninstall unistatus, run `sudo make uninstall` from within the `unistatus` directory.