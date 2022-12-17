#!/bin/bash

mkdir -p /tmp1/
pushd /tmp1/
if [ ! -d /tmp1/gperftools ]; then
    git clone https://github.com/gperftools/gperftools.git
fi
cd ./gperftools
mkdir -p my_build_out
rm -rf my_build_out/*
my_build_out=$PWD/my_build_out
./autogen.sh
./configure CXXFLAGS='-fPIC -std=c++14' --prefix=${my_build_out} --enable-libunwind -enable-frame-pointers
make -j 8
make install
popd
mkdir -p ./dep/gperftools
cp -rf /tmp1/gperftools/my_build_out/* ../dep/gperftools
