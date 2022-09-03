#include <memory>
#include <string>
#include <vector>

class IMsg
{
public:
    IMsg()
    {
        printf("call IMsg()\n");
    }
    virtual ~IMsg()
    {
        printf("call ~IMsg()\n");
    }
};

class A : public IMsg, public std::enable_shared_from_this<A>
{
public:
    A()
    {
        printf("call A()\n");
    }
    virtual ~A()
    {
        printf("call ~A()\n");
    }
};

class B : public IMsg, public std::enable_shared_from_this<B>
{
public:
    B()
    {
        printf("call B()\n");
    }
    virtual ~B()
    {
        printf("call ~B()\n");
    }
};

int main(int argc, char **argv)
{
    std::vector<std::shared_ptr<IMsg>> c;
    {
        auto a = std::make_shared<A>();
        auto b = std::make_shared<B>();
        c.push_back(a->shared_from_this());
        c.push_back(b->shared_from_this());
    }
    for (auto &v : c)
    {
        printf("use_count:%d\n", v.use_count());
    }
    c.clear();
    return 0;
}