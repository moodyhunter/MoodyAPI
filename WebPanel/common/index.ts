import { ChannelCredentials, createChannel, createClient } from "nice-grpc";
import { APIClient } from "./protos/common/common";
import { LightState } from "./protos/light/light";
import { MoodyAPIServiceDefinition } from "./protos/MoodyAPI";

export { APIClient } from "./protos/common/common";

export function getServerConnection() {
    const API_Server = process.env['API_SERVER'] ?? "localhost:1920";
    const API_TLS = Boolean(process.env['API_TLS']) ?? false;

    const channel = createChannel(API_Server, API_TLS ? ChannelCredentials.createSsl() : ChannelCredentials.createInsecure());
    return createClient(MoodyAPIServiceDefinition, channel);
}

export type CreateClientAPIResponse = { client: APIClient }
export type ListClientsAPIResponse = { clients: APIClient[] }
export type UpdateClientAPIResponse = { client: APIClient }
export type DeleteClientAPIResponse = { deleted: boolean };

export type LightAPIRequest = { state: LightState };
export type UpdateLightAPIResponse = { state: LightState };
export type GetLightAPIResponse = { state: LightState };

export type ClientAPIResponse<T> = {
    success: boolean,
    message: string,
    data: T | undefined
};
