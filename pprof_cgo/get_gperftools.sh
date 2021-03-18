#!/bin/bash

git clone https://github.com/gperftools/gperftools.git
pushd ./gperftools 
./autogen.sh
./configure CXXFLAGS='-fPIC' --prefix=$PWD/../lib --enable-libunwind -enable-frame-pointers
make -j 8
make install
popd
