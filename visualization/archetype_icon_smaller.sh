# the script converts 512x512 PNG to 128x128 for loading speed,
# input dir: archetype_icon_512, output dir: archetype_icon.
# Requires ImageMagick-6.9 installed.

cd archetype_icon_512
nTargetExisted=0
nTarnetNew=0
for file in *.png; do
    output_file="../archetype_icon/${file}"
    if [ ! -f "${output_file}" ]; then
        echo "converting ${file} to 128x128"

        convert "${file}" -resize 128x128\! "${output_file}"

        nTarnetNew=$((nTarnetNew+1))
    else
        nTargetExisted=$((nTargetExisted+1))
    fi
done

echo "done, ${nTarnetNew} new files, ${nTargetExisted} files already existed."
