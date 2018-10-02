# Overview

This project takes images in the format PNG, GIF, BMP, JPEG and displays them within your terminal as colored character art.
The goal is to have the characters displayed in color not just in black & white.

## ASCII/ANSI Art

When looking into ASCII art there is a lot of information that overwhelms you and its hard to find a starting point.
For example here is just a few phrases mentioned around the Internet:

* Fixed-width Fonts
* Grayscale
* Block Elements
* Character Ramps
* Extended ASCII
* ANSI Escape Codes
* Image Scaling
* Luminosity & Brightness
* Pixel Density / Dithering

## Approach

To simplify the steps in rendering an image as characters here are the steps I took:

- [X] Determine terminal dimensions
- [X] Load the image
- [X] Scale the image to fit terminal
- [X] Calculate brightness from Pixel Density
- [X] Map brightness to ASCII character
- [X] Draw characters as series of ANSI escape codes
- [ ] Refactor code into a Go Package

## Example Screens

[rainbow256]: /assets/rainbow-256.png "ANSI 256 Colors"
[rainbowBG256]: /assets/rainbow-background-256.png "ANSI 256 Colors"
[yakfChars]: /assets/yie-ar-kung-fu-characters.png "Yie-Ar-Kung-Fu Characters"
[yakfAnsiChars1]: /assets/yie-ar-kung-fu-ansichars1.png "Yie-Ar-Kung-Fu ANSI Characters"
[yakfAnsiChars2]: /assets/yie-ar-kung-fu-ansichars2.png "Yie-Ar-Kung-Fu ANSI Characters"


Screen | Captures
------------ | -------------
![ANSI 256 Colors][rainbow256] | ![ANSI 256 Colors][rainbowBG256]
![Yie-Ar-Kung-Fu][yakfChars] | ![Yie-Ar-Kung-Fu][yakfAnsiChars1]

![Yie-Ar-Kung-Fu][yakfAnsiChars2]

# References

* https://en.wikipedia.org/wiki/ASCII_art
* https://en.wikipedia.org/wiki/Block_Elements
* https://en.wikipedia.org/wiki/ANSI_escape_code
* https://en.wikipedia.org/wiki/Grayscale
* https://en.wikipedia.org/wiki/Dither
* https://en.wikipedia.org/wiki/Lanczos_resampling
* https://en.wikipedia.org/wiki/Vector_quantization
* https://en.wikipedia.org/wiki/Quantization_(signal_processing)
* https://www.johndcook.com/blog/2009/08/24/algorithms-convert-color-grayscale/
* http://www.roysac.com/tutorial/asciiarttutorial.html
* http://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html
* https://robertheaton.com/2018/06/12/programming-projects-for-advanced-beginners-ascii-art/
* https://www.fightersgeneration.com/games/yie-ar-kung-fu.html