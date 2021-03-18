all:
	g++ -I. -I./lib/include -g3 -O0 -c gperftools.cpp
	ar cr libgperftools.a gperftools.o
	go build -o main -gcflags=all="-N -l" gperftools.go main.go
clean:
	rm -f main libgperftools.a gperftools.o *.heap

