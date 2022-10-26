import type { NextApiRequest, NextApiResponse } from 'next';
import { getSession } from 'next-auth/react';
import { ClientAPIResponse, GetLightAPIResponse, getServerConnection, LightAPIRequest, UpdateLightAPIResponse } from '../../common';
import { Auth } from '../../common/protos/common/common';
import { SetLightResponse } from '../../common/protos/light/light';

type UpdateLightResponse = ClientAPIResponse<GetLightAPIResponse | UpdateLightAPIResponse>;

export default async function power(req: NextApiRequest, resp: NextApiResponse<UpdateLightResponse>) {
    const session = await getSession({ req });
    if (!session) {
        resp.status(401);
        resp.end();
    }

    const client = getServerConnection();
    const body: LightAPIRequest = req.body;

    const API_CLIENTID = process.env["API_CLIENTID"];
    if (!API_CLIENTID) {
        console.error("API_CLIENTID is not set.");
        resp.status(503).send({ message: "invalid server configuration", success: false, data: undefined });
        return;
    }

    const AuthObject: Auth = { clientUuid: API_CLIENTID };

    try {
        const result: SetLightResponse = await client.setLight({ auth: AuthObject, state: body.state });
        console.log(result);
        resp.status(200).json({ success: true, message: "ok", data: { state: body.state } });
    } catch (error) {
        resp.status(503).send({ message: "server error", success: false, data: undefined });
    }
}
