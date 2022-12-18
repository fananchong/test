all:
	gcc src/main.cpp src/internal/memhook.cpp src/internal/backtrace.cpp -o out -O3 -rdynamic -g -fno-builtin-malloc -fno-builtin-calloc -fno-builtin-realloc -fno-builtin-free -std=c++11 -I. -I./dep/gperftools/include -L./dep/gperftools/lib -ltcmalloc -lpthread -lstdc++ -lc -lunwind -lm
	
nouse:
	gcc src/main.cpp src/internal/memhook.cpp -o out -std=c++11 -I. -Wl,--wrap=malloc -Wl,--wrap=free -lpthread -lstdc++ -lc -lunwind -lm

clean:
	rm -rf *.o out
