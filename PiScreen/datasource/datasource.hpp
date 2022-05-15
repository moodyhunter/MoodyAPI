#pragma once

#include "IDataSource.hpp"
#include "ip/source.hpp"

#include <functional>
#include <map>

// clang-format off
#define RegisterDataSource(name)                                                                                                                                         \
    constexpr auto name##_ID = __COUNTER__;                                                                                                                              \
    namespace __details                                                                                                                                                  \
    {                                                                                                                                                                    \
        struct name##auto_registration { name##auto_registration() { registrations.insert({ name##_ID, std::function([]() { return new name(); }) }); } };               \
        inline volatile name##auto_registration name##auto_reg;                                                                                                          \
    }
// clang-format on

namespace PiScreen::datasource
{
    RegisterDataSource(IPAddressDataSource);
} // namespace PiScreen::datasource