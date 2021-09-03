#!/usr/bin/env sh


mkdir ./output

cpt=1
while read link; do
  curl -s -S $link -o "./output/image-${cpt}"
  cpt=$(expr $cpt + 1)
  echo $cpt
done < <(cat ./russian-images.csv | cut -d ',' -f 1)
