#pragma once

#include <iostream>

namespace PiScreen::common
{
    [[noreturn]] inline void Unreachable()
    {
        std::cout << "Unreachable code reached." << std::endl;
        std::abort();
    }
} // namespace PiScreen::common
