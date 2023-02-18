#include <gtest/gtest.h>
#include "math/add.hpp"

// Demonstrate some basic assertions.
TEST(MyTest, Add1)
{
    auto v = add(7, 6);
    EXPECT_EQ(v, 13);
}