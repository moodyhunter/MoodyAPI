import { ChannelCredentials, createChannel, createClient } from "nice-grpc";
import { APIClient, MoodyAPIServiceDefinition } from "./protos/MoodyAPI";

export { APIClient } from "./protos/MoodyAPI";

export function getServerConnection() {
    const API_Server = process.env['API_SERVER'] ?? "localhost:1920";
    const API_TLS = Boolean(process.env['API_TLS']) ?? false;

    const channel = createChannel(API_Server, API_TLS ? ChannelCredentials.createSsl() : ChannelCredentials.createInsecure());
    return createClient(MoodyAPIServiceDefinition, channel);
}

export type CreateClientAPIResponse = { client: APIClient }
export type UpdateClientAPIResponse = { client: APIClient }
export type DeleteClientAPIResponse = { deleted: boolean };
export type ListClientsAPIResponse = { clients: APIClient[] }

export declare type ClientAPIResponse<T> = {
    success: boolean,
    message: string,
    data: T | undefined
};
