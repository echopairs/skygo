#!/usr/bin/env bash

# ubuntu 16.4
#set -e
script=$(readlink -f "$0")
route=$(dirname "$script")

### 0. get ready
git_branch=$(git branch | grep "*" | cut -d" " -f2)
pkg_name=zweb
tgt_install_prefix=/opt
version=$(git describe --abbrev=0 --tags 2>/dev/null)
revcnt=$(git rev-list --count HEAD 2>/dev/null)
working_dir=gweb

if [[ "$1" != "-t" ]] && [[ "$git_branch" != "master" ]]; then
    echo "release must be on master branch"
    exit
fi

if test -z "$version"
then
    version=$revcnt.$(git rev-parse --short HEAD)
fi


### 1. make the working dir
mkdir -p ${route}/../dist/${working_dir}
mkdir -p ${route}/../dist/${working_dir}/DEBIAN
mkdir -p ${route}/../dist/${working_dir}/${tgt_install_prefix}/${pkg_name}

### 2. make various file under DEBIAN dir
cd ${route}/../dist/${working_dir}/DEBIAN
touch control
(cat << EOF
Package: gweb
Version: ${version}
Section: x11
Priority: optional
Depends:
Suggests:
Architecture: amd64
Maintainer: zsy
CopyRight: commercial
Provider: zsy.
Description: golang web.
EOF
) > control

touch postinst
chmod a+x postinst
# todo
touch postrm
chmod a+x postrm
# todo
touch prerm
chmod a+x prerm
# todo

### 3. cp binary and config file ...

${route}/build.sh
cp -rf ${route}/../cmd/bin/cloudweb ${route}/../dist/${working_dir}/${tgt_install_prefix}/${pkg_name}

### 4. make xxx.deb package
cd ${route}/../

dpkg -b dist/${working_dir} dist/gweb_${version}_amd64.deb