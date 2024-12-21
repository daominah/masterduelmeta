# the script converts 512x512 PNG to 128x128 for loading speed,
# input dir: archetype_icon_512, output dir: archetype_icon.
# Requires ImageMagick-6.9 installed.

cd archetype_icon_512
mogrify -path ../archetype_icon/ -format png -resize 128x128\! *.png
