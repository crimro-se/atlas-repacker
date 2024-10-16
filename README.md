# Atlas Repacker

This is a tool made to identify pixel islands in an image and repack them onto a new canvas of a specified size, with optional padding between such islands.

A pixel is considered real for island detection purposes if it has alpha > 0, so you'll probably be using this tool on PNGs

actual packing is deferred to https://github.com/nothings/stb/blob/master/stb_rect_pack.h

I needed this for some specific AI training, it may or may not fit your needs.

## Usage

```
atlas-repacker [flags] [input.png] [input2.png ...]
Flags:
  -diagonal
        when set, diagonally adjacent pixels are considered connected during island detection.
  -h int
        height of output image (default 512)
  -margin int
        margin to use for each box (default 1)
  -o string
        filename of output (default "output.png")
  -offset int
        ammount to offset each box. Useful values are 0, margin/2, =margin
  -w int
        width of output image (default 512)
```

## TODO

- unit tests
- add chroma mask support
- .atlas file read/write support (maybe)