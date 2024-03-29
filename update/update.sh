#!/bin/bash

# run it from the top directory

set -eu
set -o pipefail

function _file_to_url {
    local fname=$1
    if echo "$fname" | grep -q "html$"; then
        echo "https://htmlpreview.github.io/?https://github.com/jakub-m/bulletin/blob/mainline/$fname"
    else
        echo $fname
    fi
}

function _date_from_fname {
    local fname=$1
    echo $fname | perl -ne 'print "$1" if /bulletin-([\d-]+).(html|md)/'
}

nl="
"

dir=$(cd $(dirname $0) && pwd)

cd bulletins
../bin/bulletin
cd .. 

template_file="$dir/README.template.md"


recent_file=$(find bulletins -name bulletin\*.html -or -name bulletin\*.md | sort | tail -n1)
recent_url=$(_file_to_url $recent_file)

sed "s|__CURRENT__|$recent_url|" $template_file

for fname in $(find bulletins -name bulletin\*.html -or -name bulletin\*.md | sort -r); do
    echo "- [$(_date_from_fname $fname)]($(_file_to_url $fname))"
done
