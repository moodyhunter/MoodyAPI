import type { NextApiRequest, NextApiResponse } from 'next';
import { getServerConnection } from '../../common';
import { APIClient, Auth, UpdateClientInfoRequest } from '../../common/protos/MoodyAPI';


export default async function clients(req: NextApiRequest, resp: NextApiResponse<APIClient | APIClient[]>) {
    const client = getServerConnection();
    const requestedClient: APIClient = req.body;

    const API_CLIENTID = process.env["API_CLIENTID"];
    if (!API_CLIENTID) {
        console.error("API_CLIENTID is not set.");
        resp.status(503).send(requestedClient);
        return;
    }

    const AuthObject: Auth = { clientId: API_CLIENTID };

    if (req.method == "PATCH") {
        const request: UpdateClientInfoRequest = {
            auth: AuthObject,
            clientInfo: requestedClient
        };
        const response = await client.updateClientInfo(request);
        resp.status(response.success ? 200 : 403).json(requestedClient);
    } else if (req.method == "GET") {
        const result = await client.listClients({ auth: AuthObject });
        resp.status(result.success ? 200 : 403).json(result.clients);
    }

    return;
}
