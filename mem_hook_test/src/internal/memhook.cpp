#include "memhook.hpp"
#include <gperftools/malloc_hook.h>
#include <gperftools/tcmalloc.h>
#include <gperftools/nallocx.h>

// extern "C"
// {
//     void *__wrap_malloc(size_t size)
//     {
//         void *ptr = __real_malloc(sizeof(int64_t) + size);
//         printf("__wrap_malloc, ptr=%p\n", &ptr);
//         *(int64_t *)ptr = size;
//         return (int64_t *)ptr + 1;
//     }

//     void __wrap_free(void *ptr)
//     {
//         printf("__wrap_free, ptr=%p\n", &ptr);
//         __real_free((int64_t *)ptr - 1);
//     }
// }

void MyNewHook(const void *ptr, size_t size)
{
    if (ptr != nullptr)
    {
        auto real_size = tc_nallocx(size, 0);
        printf("new size: %d, real_size: %d, ptr=%p\n", size, real_size, &ptr);

        // // 例子打印堆栈
        // if (size > 1024)
        // {
        //     MallocHook::GetCallerStackTrace()
        // }
    }
}

void MyDeleteHook(const void *ptr)
{
    if (ptr != nullptr)
    {
        int64_t size = tc_malloc_size(const_cast<void *>(ptr));
        printf("delete size: %d, ptr=%p\n", size, &ptr);
    }
}

void init_mem_hook()
{
    MallocHook::AddNewHook(&MyNewHook);
    MallocHook::AddDeleteHook(&MyDeleteHook);
}
