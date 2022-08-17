#define ANKERL_NANOBENCH_IMPLEMENT
#include <nanobench.h>
#include <chrono>

inline int64_t get_current_time_nanos()
{
    return std::chrono::duration_cast<std::chrono::nanoseconds>(
               std::chrono::system_clock::now() -
               std::chrono::system_clock::from_time_t(0))
        .count();
}

int main()
{
    int64_t d = 0;
    auto fn = [&]
    {
        d = get_current_time_nanos();
        ankerl::nanobench::doNotOptimizeAway(d);
    };
    ankerl::nanobench::Bench().run("time test", fn);
}