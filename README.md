# img2rgb565

A command-line tool that converts image files to RGB565 format C arrays for embedded systems and displays.

## Features

- Converts images to 16-bit RGB565 format
- Generates paired .c/.h files with uint16_t arrays
- Supports BMP, JPEG, PNG, and GIF formats
- Optimized for embedded graphics applications

## Build

```bash
# Using Make
make build

# Using Go directly
go build
```

## Usage

```bash
./img2rgb565 image.png
```

This generates:
- `image.c` - Array definition with RGB565 pixel data
- `image.h` - Header file with extern declaration

### Example

```bash
$ ./img2rgb565 logo.png
# Creates logo.c and logo.h

# In your C code:
#include "logo.h"
// Use img_logo[] array (16-bit RGB565 values)
```

## Adding Format Support

To support additional image formats, add the corresponding decoder to the import section in `main.go`.
