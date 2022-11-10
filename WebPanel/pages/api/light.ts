import type { NextApiRequest, NextApiResponse } from 'next';
import { getSession } from 'next-auth/react';
import { ClientAPIResponse, createAuth, GetLightAPIResponse, getServerConnection, LightAPIRequest, UpdateLightAPIResponse } from '../../common';

type UpdateLightResponse = ClientAPIResponse<GetLightAPIResponse | UpdateLightAPIResponse>;

export default async function power(req: NextApiRequest, resp: NextApiResponse<UpdateLightResponse>) {
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

    const body: LightAPIRequest = req.body;


    try {
        if (req.method === "GET") {
            const light = await client.getLightState({ auth: AuthObject });
            if (light.state)
                resp.status(200).send({ message: "success", success: true, data: { state: light.state } });
            else
                resp.status(500).send({ message: "failed to get light state", success: false, data: undefined });
        }
        else if (req.method === "POST") {
            await client.setLightState({ auth: AuthObject, state: body.state });
            resp.status(200).json({ success: true, message: "ok", data: { state: body.state } });
        }
    } catch (error) {
        console.log(error)
        resp.status(503).send({ message: "server error", success: false, data: undefined });
    }
}
