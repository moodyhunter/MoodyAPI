import type { NextApiRequest, NextApiResponse } from 'next';

type Data = {
    name: string
}

export default function handler(req: NextApiRequest, res: NextApiResponse<Data>) {
    if (req.method == "PATCH") {
        console.log("Updating a client.");
    }
    res.status(200).json(req.body)
}
