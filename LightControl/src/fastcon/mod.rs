// This is a reverse-engineered implementation of the BLE FastCon protocol
// used by some light bulbs.

// parse the broadcast packets
pub mod broadcast_parser;

// wrapper around some of the commands
pub mod command_wrapper;

// device type enum
pub mod device_type;
