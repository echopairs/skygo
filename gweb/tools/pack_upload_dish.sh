#!/usr/bin/env bash

# ubuntu 16.4
set -e
script=$(readlink -f "$0")
route=$(dirname "$script")

### 0. get ready
git_branch=$(git branch | grep "*" | cut -d" " -f2)
pkg_name=zweb
tgt_install_prefix=/opt
version=$(git describe --abbrev=0 --tags)
working_dir=zweb

if [ "$1" != "-t" ] && [ "$git_branch" != "master" ]; then
    echo "release must be on master branch"
    exit
fi

### 1. make the working dir

### 2. make various file under DEBIAN dir

### 3. cp binary and config file ...