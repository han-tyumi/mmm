#!/usr/bin/env bash

set -e

readonly build_dir="build"
readonly out_name="mmm"
readonly file_prefix="${out_name}_"

declare -A targets=(
    [windows]="amd64 386"
    [darwin]="amd64"
    [linux]="amd64 386"
)

mkdir -p $build_dir
rm $build_dir/*

for os in "${!targets[@]}"
do
    IFS=" " read -r -a arches <<< "${targets[$os]}"

    for arch in "${arches[@]}"
    do
        env GOOS="$os" GOARCH="$arch" go build -x

        file_suffix="$os-$arch"
        base_path="$build_dir/$file_prefix$file_suffix"

        if [[ $os = "windows" ]]
        then
            out_file="$out_name.exe"
            zip -v9 "$base_path" $out_file
        else
            out_file="$out_name"
            tar -czvf "$base_path.tar.gz" $out_file
        fi

        rm $out_file
    done
done
