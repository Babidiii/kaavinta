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
  done <<< $( tail -n +2 ${file} )
}

convert_files(){
  echo "* Convert each files into .ogg with max 256kb/s max bit rate"
  mkdir -p "${res_dir}"
  for i in "${tmp_dir}"/*; do 
    name=${i##*/}
    echo "* Converting ${name} file to .ogg"
    if ! ffmpeg -hide_banner -loglevel error -i "${i}" -b:a 256000 "${res_dir}/${name%.*}.ogg"; then
      echo "[ffmpeg](fail) ${i}" >> errors.log
    fi
  done

  # clean all zero bytes files
  find ${res_dir} -name "*" -size 0 -print0 | xargs -0 rm

}

try_mp3(){
  while read -r f; do
    name=${f##*/}
#    echo "* Converting ${name} file to .mp3"
    echo "* $f"
    echo "* ${res_dir}/${name%.*}.mp3"
    if ! ffmpeg -hide_banner -loglevel error -i "${f}" -b:a 256000 "${res_dir}/${name%.*}.mp3"; then
      echo "[ffmpeg](fail) ${f}" >> errors_try.log
    fi
  done <<< $(grep "ffmpeg" ./errors.log | cut -d ' ' -f 2-)
}

#download
#convert_files
try_mp3
