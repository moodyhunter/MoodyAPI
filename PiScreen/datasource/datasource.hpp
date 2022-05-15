#pragma once

#include "IDataSource.hpp"
#include "datasource/cpu/source.hpp"
#include "datasource/mem/source.hpp"
#include "datasource/service/source.hpp"
#include "ip/source.hpp"

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
    RegisterDataSource(MemoryDataSource);
    RegisterDataSource(CPUDataSource);
    RegisterDataSource(SystemdServiceDataSource);
} // namespace PiScreen::datasource