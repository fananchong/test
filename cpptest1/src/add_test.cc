#include <gtest/gtest.h>
#include "add.hpp"

// Demonstrate some basic assertions.
TEST(MyTest, Add)
{
    auto v = add(7, 6);
    EXPECT_EQ(v, 13);
}