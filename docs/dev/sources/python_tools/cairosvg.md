# CairoSVG Documentation

> Source: https://cairosvg.org/
> Fetched: 2026-02-01T11:55:34.106120+00:00
> Content-Hash: c1186671992f4a52
> Type: html

---

CairoSVG

- [Home](/)

- [News](/news/)
- [Documentation](/documentation/)

# CairoSVG

Convert your SVG files to PDF and PNG.

## About CairoSVG

CairoSVG is:

- a SVG 1.1 to PNG, PDF, PS and SVG converter;
- a command-line interface;
- a Python 3.6+ library;
- known to work at least on Linux, OS X, and Windows;
- based on the Cairo 2D graphics library;
- tested using the W3C test suite;
- LGPLv3-licensed free software.

## Let's go

### Install

CairoSVG is available on [PyPI](https://pypi.python.org/pypi/CairoSVG/), you can install it with pip:

    pip3 install cairosvg

You can have more information and help in the ["Installation" part of the documentation](/documentation/#installation).

### Convert

You can use CairoSVG as a standalone command-line program:

    cairosvg image.svg -o image.png

### Embed

You can also use CairoSVG as a Python 3 library:

    $ python3
    >>> import cairosvg
    >>> cairosvg.svg2pdf(url='image.svg', write_to='image.pdf')

### Want more?

Please read the [documentation](/documentation/) to learn more about how to use CairoSVG.

## What's new?

Latest version of CairoSVG is 2.7.1, released on August 5, 2023 ([changelog](https://github.com/Kozea/CairoSVG/blob/master/NEWS.rst)).

### Version 2.5.2

March 6, 2021

CairoSVG 2.5.2 has been released!

- Fix marker path scale

[Read more news…](/news/)

- [A CourtBouillon Project](https://www.courtbouillon.org)
- [Fork me on GitHub](https://github.com/Kozea/CairoSVG)

  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
