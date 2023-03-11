#include <coroutine>
#include <iostream>
#include <stdexcept>
#include <thread>

auto switch_to_new_thread(std::jthread &out)
{
    std::cout << "switch_to_new_thread()" << std::endl;
    struct awaitable
    {
        awaitable(std::jthread *out)
            : p_out(out)
        {
            std::cout << "awaitable()" << std::endl;
        }

        ~awaitable()
        {
            std::cout << "~awaitable()" << std::endl;
        }
        std::jthread *p_out;
        bool await_ready()
        {
            std::cout << "awaitable::await_ready()" << std::endl;
            return false;
        }
        void await_suspend(std::coroutine_handle<> h)
        {
            std::cout << "awaitable::await_suspend()" << std::endl;
            std::jthread &out = *p_out;
            if (out.joinable())
                throw std::runtime_error("jthread 输出参数非空");
            out = std::jthread([h]
                               { h.resume(); });
            // 潜在的未定义行为：访问潜在被销毁的 *this
            // std::cout << "新线程 ID：" << p_out->get_id() << '\n';
            std::cout << "新线程 ID：" << out.get_id() << '\n'; // 这样没问题
        }
        void await_resume()
        {
            std::cout << "awaitable::await_resume()" << std::endl;
        }
    };
    return awaitable{&out};
}

struct task
{
    task()
    {
        std::cout << "task()" << std::endl;
    }

    ~task()
    {
        std::cout << "~task()" << std::endl;
    }

    struct promise_type
    {
        promise_type()
        {
            std::cout << "promise_type()" << std::endl;
        }

        ~promise_type()
        {
            std::cout << "~promise_type()" << std::endl;
        }

        task get_return_object()
        {
            std::cout << "promise_type::get_return_object()" << std::endl;
            return {};
        }
        std::suspend_never initial_suspend()
        {
            std::cout << "promise_type::initial_suspend()" << std::endl;
            return {};
        }
        std::suspend_never final_suspend() noexcept
        {
            std::cout << "promise_type::final_suspend()" << std::endl;
            return {};
        }
        void return_void()
        {
            std::cout << "promise_type::return_void()" << std::endl;
        }
        void unhandled_exception()
        {
            std::cout << "promise_type::unhandled_exception()" << std::endl;
        }
    };
};

task resuming_on_new_thread(std::jthread &out)
{
    std::cout << "resuming_on_new_thread()" << std::endl;
    std::cout << "协程开始，线程 ID：" << std::this_thread::get_id() << '\n';
    co_await switch_to_new_thread(out);
    // 等待器在此销毁
    std::cout << "协程恢复，线程 ID：" << std::this_thread::get_id() << '\n';
}

int main()
{
    std::jthread out;
    resuming_on_new_thread(out);
}