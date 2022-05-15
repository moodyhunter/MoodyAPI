#pragma once

#include "IDataSource.hpp"
#include "ip/source.hpp"

#include <map>

// clang-format off
#define RegisterDataSource(name)                                                                                                                                         \
    constexpr auto name##_ID = __COUNTER__;                                                                                                                              \
    namespace __details                                                                                                                                                  \
    {                                                                                                                                                                    \
        struct name##_dsreg { name##_dsreg() { registrations.insert({ name##_ID, []() { return (IDataSource *) new name(); } }); } };                                    \
        inline volatile name##_dsreg name##_reg;                                                                                                                         \
    }
// clang-format on

namespace PiScreen::datasource
{
    RegisterDataSource(IPAddressDataSource);
} // namespace PiScreen::datasource