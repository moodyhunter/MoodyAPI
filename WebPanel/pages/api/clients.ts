import { ServiceError } from '@grpc/grpc-js';
import type { NextApiRequest, NextApiResponse } from 'next';
import { ClientAPIResponse, CreateClientAPIResponse, DeleteClientAPIResponse, getServerConnection, ListClientsAPIResponse, UpdateClientAPIResponse } from '../../common';
import { APIClient, Auth } from '../../common/protos/MoodyAPI';

type ClientAPIServerResponse = ClientAPIResponse<CreateClientAPIResponse | ListClientsAPIResponse | UpdateClientAPIResponse | DeleteClientAPIResponse>

export default async function clients(req: NextApiRequest, resp: NextApiResponse<ClientAPIServerResponse>) {
    const client = getServerConnection();
    const requestedClient: APIClient = req.body;

    const API_CLIENTID = process.env["API_CLIENTID"];
    if (!API_CLIENTID) {
        console.error("API_CLIENTID is not set.");
        resp.status(503).send({ message: "invalid server configuration", success: false, data: undefined });
        return;
    }

    await new Promise(f => setTimeout(f, 2000));

    const AuthObject: Auth = { clientUuid: API_CLIENTID };
    try {
        if (req.method == "GET") {
            const result = await client.listClients({ auth: AuthObject });
            resp.status(result.success ? 200 : 400).json({ success: result.success, message: "ok", data: { clients: result.clients } as ListClientsAPIResponse });
            return;
        } else if (req.method == "POST") {
            const result = await client.createClient({ auth: AuthObject, client: requestedClient });
            resp.status(result.success ? 201 : 400).json({ success: result.success, message: "ok", data: { client: result.client } as CreateClientAPIResponse });
            return;
        } else if (req.method == "PATCH") {
            const result = await client.updateClient({ auth: AuthObject, client: requestedClient });
            resp.status(result.success ? 200 : 400).json({ success: result.success, message: "ok", data: { client: requestedClient } as UpdateClientAPIResponse });
            return;
        } else if (req.method == "DELETE") {
            const result = await client.deleteClient({ auth: AuthObject, client: requestedClient });
            resp.status(result.success ? 200 : 400).json({ success: result.success, message: "ok", data: { deleted: result.success } as DeleteClientAPIResponse });
            return;
        }
    } catch (error) {
        resp.status(400).json({ success: false, message: (error as ServiceError).details, data: undefined });
        return;
    }

    resp.status(405).json({ success: false, message: "Method not allowed", data: undefined });
}
