#!/bin/bash

set -e
[[ -z $DEBUG ]] || set -x

if [ $EUID != 0 ]; then
    sudo "$0" "$@"
    exit $?
fi

machine=""
case $(uname -m) in
x86_64) machine="amd64";;
i386|i686) machine="386";;
esac

os=""
case $(uname -s) in
Linux)  os="linux";;
Darwin) os="darwin";;
*)      echo "OS not supported" && exit 1;;
esac

latest=$(curl -sL 'https://api.github.com/repos/lighttiger2505/lab/releases/latest' | grep tag_name | grep --only 'v[0-9\.]\+' | cut -c2-)
curl -sL "https://github.com/lighttiger2505/lab/releases/download/v${latest}/lab-v${latest}-${os}-${machine}.tar.gz" | tar -C /tmp/ -xzf -
cp /tmp/${os}-${machine}/lab /usr/local/bin/lab
echo "Successfully installed lab into /usr/local/bin/"
