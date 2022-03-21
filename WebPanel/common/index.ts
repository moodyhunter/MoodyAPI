import { createChannel, ChannelCredentials, createClient } from "nice-grpc";
import { MoodyAPIServiceDefinition } from "./protos/MoodyAPI";

export { APIClient } from "./protos/MoodyAPI";

export function getServerConnection() {
    const API_Server = process.env['API_SERVER'] ?? "localhost:1920";
    const API_TLS = Boolean(process.env['API_TLS']) ?? false;

    const channel = createChannel(API_Server, API_TLS ? ChannelCredentials.createSsl() : ChannelCredentials.createInsecure());
    return createClient(MoodyAPIServiceDefinition, channel);
}
