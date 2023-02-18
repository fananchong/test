#include <gtest/gtest.h>
#include "../add.hpp"

// Demonstrate some basic assertions.
TEST(MyTest, Add2)
{
    auto v = add(71, 6);
    EXPECT_EQ(v, 77);
}