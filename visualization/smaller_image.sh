# the script converts 512x512 PNG to 128x128 JPG for loading speed,
# input dir: archetype_icon_512, output dir: archetype_icon

cd archetype_icon_512
mogrify -path ../archetype_icon/ -format jpg -quality 90 -resize 128x128 *.png
