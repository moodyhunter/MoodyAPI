import { ServiceError } from '@grpc/grpc-js';
import type { NextApiRequest, NextApiResponse } from 'next';
import { getSession } from 'next-auth/react';
import { APIClient, ClientAPIResponse, createAuth, CreateClientAPIResponse, DeleteClientAPIResponse, getServerConnection, ListClientsAPIResponse, UpdateClientAPIResponse } from '../../common';

type ClientAPIServerResponse = ClientAPIResponse<CreateClientAPIResponse | ListClientsAPIResponse | UpdateClientAPIResponse | DeleteClientAPIResponse>

export default async function clients(req: NextApiRequest, resp: NextApiResponse<ClientAPIServerResponse>) {
    const session = await getSession({ req });
    if (!session) {
        resp.status(401);
        resp.end();
    }

    const client = getServerConnection();
    const AuthObject = createAuth();

    if (!AuthObject) {
        resp.status(500);
        resp.end();
    }

    const requestedClient: APIClient = req.body;

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
