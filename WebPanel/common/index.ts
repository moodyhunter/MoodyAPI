import { Channel, ChannelCredentials, Client, createChannel, createClient } from "nice-grpc";
import { APIClient, Auth } from "./protos/common/common";
import { LightState } from "./protos/light/light";
import { MoodyAPIServiceDefinition } from "./protos/MoodyAPI";

export { APIClient } from "./protos/common/common";

let channel: Channel | undefined;
let client: Client<MoodyAPIServiceDefinition> | undefined;

export function getServerConnection() {
    const API_Server = process.env['API_SERVER'] ?? "localhost:1920";
    const API_TLS = Boolean(process.env['API_TLS']) ?? false;

    if (channel === undefined || client === undefined) {
        if (channel === undefined)
            channel = createChannel(API_Server, API_TLS ? ChannelCredentials.createSsl() : ChannelCredentials.createInsecure());
        if (client === undefined)
            client = createClient(MoodyAPIServiceDefinition, channel);
    }

    return client;
}

export function createAuth(): Auth | undefined {
    const API_CLIENTID = process.env["API_CLIENTID"];
    if (!API_CLIENTID) {
        console.error("API_CLIENTID is not set.");
        return undefined;
    }

    const AuthObject: Auth = { clientUuid: API_CLIENTID };
    return AuthObject;
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
