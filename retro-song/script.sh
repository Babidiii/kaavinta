#!/usr/bin/env bash

#songs=$(tail -n +2 ./songs.csv)
output_dir="/tmp/retro-songs"
tmp_dir="${output_dir}/tmp"
res_dir="${output_dir}/res"
file="./songs.csv"


download(){
  mkdir -p "${tmp_dir}"
  echo "* $(cat ${file} | wc -l) songs to download"

  while read -r s; do
   title="$(echo $s | cut -d ',' -f 1)"
   link="$(echo $s | cut -d ',' -f 3)"
   echo "* Downloading ${title} ${link}"
   if ! youtube-dl -q -i --extract-audio -f 140  -o "${tmp_dir}/%(title)s.%(ext)s" "${link}" ; then
     echo "[ytb-dl](fail) ${title} --> ${link}" >> errors.log
   fi
  done <<< $( tail -n +2 ${file} )v
}

convert_files(){
  echo "* Convert each files into .ogg with max 256kb/s max bit rate"
  mkdir -p "${res_dir}"
  for i in "${tmp_dir}"/*; do 
    echo "* Converting ${i} file"
    name=${i##*/}
    echo "   to ${res_dir}/${name%.*}.ogg"
    if ! ffmpeg -hide_banner -loglevel error -i "${i}" -b:a 256000 "${res_dir}/${name%.*}.ogg"; then
      echo "[ffmpeg](fail) ${i}" >> errors.log
    fi
  done
}

#download
convert_files
