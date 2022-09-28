#pragma once

#include <string>
#include <functional>

void sigsetup();
void register_crash_exit_fn(const std::function<void(void)> &fn);
