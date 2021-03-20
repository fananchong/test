all:
	g++ -I. -I./lib/include -c gperftools.cpp
	ar cr libgperftools.a gperftools.o
	go build -o main -gcflags=all="-N -l" gperftools.go main.go
clean:
	rm -f main gperftools.o libgperftools.a *.heap

