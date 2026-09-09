#!/usr/bin/env python3
"""Regenerate the derived raster assets committed under static/images/.

Three of the site's images are derived from sources the Go binary cannot read:

  * the portrait is authored as WebP, which PDF cannot embed at all;
  * the logo is an Illustrator SVG with 275 paths and clip paths;
  * the Open Graph cards are compositions that need a text renderer.

Rather than pull an image toolchain into the server, they are rendered once
here and committed. This script exists so the outputs are reproducible instead
of being undocumented binaries in the tree.

Usage:  python scripts/generate_assets.py
Needs:  Pillow, ImageMagick (`magick` on PATH), and a system sans-serif.
"""

from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path

from PIL import Image, ImageDraw, ImageFilter, ImageFont

ROOT = Path(__file__).resolve().parent.parent
IMAGES = ROOT / "static" / "images"
PORTRAIT_SRC = IMAGES / "photos" / "ewerton.webp"
LOGO_SRC = IMAGES / "logo.svg"

# The palette the PDF header uses, so every surface reads as one brand.
BRAND_DEEP = (0x00, 0x33, 0x66)
BRAND = (0x00, 0x45, 0x87)
BRAND_LIGHT = (0x18, 0x7C, 0xD4)
BAND_TEXT = (0xCB, 0xDF, 0xF5)
GREEN = (0x00, 0x97, 0x39)
GOLD = (0xFE, 0xDD, 0x00)
WHITE = (0xFF, 0xFF, 0xFF)

# Open Graph cards. 1200x630 is the size every major renderer crops to.
OG_SIZE = (1200, 630)
OG_CARDS = {
    "en": {
        "headline": "Senior Backend & Data Engineer",
        "location": "Belo Horizonte, Brazil",
    },
    "pt-br": {
        "headline": "Engenheiro Backend & de Dados Sênior",
        "location": "Belo Horizonte, Brasil",
    },
}
NAME = "Ewerton de Souza"
SITE = "ewerton.com.br"
LINKEDIN = "linkedin.com/in/ewertones"

FONT_CANDIDATES = {
    "bold": ["C:/Windows/Fonts/segoeuib.ttf", "C:/Windows/Fonts/arialbd.ttf",
             "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"],
    "regular": ["C:/Windows/Fonts/segoeui.ttf", "C:/Windows/Fonts/arial.ttf",
                "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"],
}


def font(weight: str, size: int) -> ImageFont.FreeTypeFont:
    for path in FONT_CANDIDATES[weight]:
        if os.path.exists(path):
            return ImageFont.truetype(path, size)
    raise SystemExit(f"no {weight} font found; edit FONT_CANDIDATES")


def lerp(a, b, t):
    return tuple(int(a[i] + (b[i] - a[i]) * t) for i in range(3))


def portrait_jpeg() -> None:
    """Portrait for the PDF header. PDF's DCTDecode filter is JPEG, so the
    bytes are embedded verbatim; 360px is ~320 DPI at its printed size."""
    out = IMAGES / "photos" / "ewerton.jpg"
    im = Image.open(PORTRAIT_SRC).convert("RGB").resize((360, 360), Image.LANCZOS)
    im.save(out, "JPEG", quality=85, optimize=True, progressive=False, subsampling=2)
    report(out)


def logo_jpeg() -> None:
    """Logo mark for the PDF footer. Flattened onto white because the footer
    sits on white paper and DCTDecode carries no alpha channel."""
    out = IMAGES / "logo-mark.jpg"
    subprocess.run(
        ["magick", "-background", "white", "-density", "400", str(LOGO_SRC),
         "-resize", "160x160", "-flatten", "-strip", "-quality", "88",
         "-sampling-factor", "1x1", "-interlace", "none", str(out)],
        check=True,
    )
    report(out)


def og_card(locale: str, headline: str, location: str) -> None:
    """Social card, composed to mirror the PDF header band."""
    width, height = OG_SIZE

    img = Image.new("RGB", OG_SIZE)
    px = img.load()
    for y in range(height):
        for x in range(width):
            t = (x / width + y / height) / 2
            px[x, y] = lerp(BRAND_DEEP, BRAND, t / 0.5) if t < 0.5 \
                else lerp(BRAND, BRAND_LIGHT, (t - 0.5) / 0.5)

    orbs = Image.new("RGBA", OG_SIZE, (0, 0, 0, 0))
    draw = ImageDraw.Draw(orbs)
    draw.ellipse([width - 330, -230, width + 130, 230], fill=(255, 255, 255, 20))
    draw.ellipse([width - 620, height - 160, width - 300, height + 160], fill=(255, 255, 255, 16))
    img = Image.alpha_composite(img.convert("RGBA"), orbs.filter(ImageFilter.GaussianBlur(2))).convert("RGB")

    photo = Image.open(PORTRAIT_SRC).convert("RGB").resize((300, 300), Image.LANCZOS)
    mask = Image.new("L", (300, 300), 0)
    ImageDraw.Draw(mask).ellipse([0, 0, 299, 299], fill=255)
    ox, oy = 96, 165

    halo = Image.new("RGBA", OG_SIZE, (0, 0, 0, 0))
    ImageDraw.Draw(halo).ellipse([ox - 22, oy - 22, ox + 322, oy + 322], fill=(255, 255, 255, 34))
    img = Image.alpha_composite(img.convert("RGBA"), halo.filter(ImageFilter.GaussianBlur(6))).convert("RGB")
    img.paste(photo, (ox, oy), mask)

    draw = ImageDraw.Draw(img)
    draw.ellipse([ox, oy, ox + 299, oy + 299], outline=WHITE, width=5)

    tx = 452
    draw.text((tx, 196), NAME, font=font("bold", 76), fill=WHITE)
    draw.text((tx, 292), headline, font=font("regular", 38), fill=BAND_TEXT)
    draw.rectangle([tx, 352, tx + 108, 358], fill=GOLD)
    draw.text((tx, 392), SITE, font=font("bold", 30), fill=WHITE)
    draw.text((tx, 436), f"{location}  ·  {LINKEDIN}", font=font("regular", 26), fill=BAND_TEXT)

    # Flag ribbon along the bottom edge, as on the PDF header.
    ribbon = 12
    for x in range(width):
        t = x / width
        if t < 0.62:
            colour = BRAND_LIGHT
        elif t < 0.85:
            colour = lerp(BRAND_LIGHT, GREEN, (t - 0.62) / 0.23)
        else:
            colour = lerp(GREEN, GOLD, (t - 0.85) / 0.15)
        draw.rectangle([x, height - ribbon, x + 1, height], fill=colour)

    out = IMAGES / f"og-card-{locale}.jpg"
    img.save(out, "JPEG", quality=90, optimize=True, progressive=False)
    report(out)


def report(path: Path) -> None:
    print(f"  {path.relative_to(ROOT)}  ({path.stat().st_size:,} bytes)")


def main() -> int:
    print("Regenerating derived assets:")
    portrait_jpeg()
    logo_jpeg()
    for locale, card in OG_CARDS.items():
        og_card(locale, **card)
    return 0


if __name__ == "__main__":
    sys.exit(main())
