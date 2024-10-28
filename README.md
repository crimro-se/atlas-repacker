# Atlas Repacker

![Example image](example.png)

This is a tool made to identify pixel islands in image(s) and repack them onto a new canvas of a specified size, with optional padding between such islands.

A pixel is considered real for island detection purposes if it has alpha > 0, so you'll probably be using this tool on PNGs

actual packing is deferred to https://github.com/nothings/stb/blob/master/stb_rect_pack.h

I needed this for some specific AI training, it might not fit your needs.

## Building/Installing

(Go)[https://go.dev] and (cgo)[https://github.com/go101/go101/wiki/CGO-Environment-Setup] are required, then simply:

```bash
go install github.com/crimro-se/atlas-repacker@latest
```

I'll add binaries when the project matures.

## Usage

```
atlas-repacker [flags] [input.png] [input2.png ...]
Flags:
  -align int
        how to align a box within its margin?
        0 = top left, 1 = center, 2 = bottom right (default 1)
  -diagonal
        when set, diagonally adjacent pixels are considered connected during island detection.
  -findmaxmargin
        when set, will find the largest margin value for which all islands still fit in the output.
  -findminsquare int
        If set > 0, finds the smallest output image size for which w and h is a multiple of this value.
  -h int
        height of output image (default 512)
  -margin int
        margin to use for each box (default 1)
  -o string
        filename of output (default "output.png")
  -w int
        width of output image (default 512)
```

## TODO

- add chroma mask support
- .atlas file read/write support (maybe)