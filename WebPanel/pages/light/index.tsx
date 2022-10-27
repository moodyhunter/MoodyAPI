import { Container, Switch } from '@mui/material';
import { GetServerSideProps } from 'next';
import { useState } from 'react';
import { LightAPIRequest, UpdateLightAPIResponse } from '../../common';

export const getServerSideProps: GetServerSideProps = async () => {
    return {
        props: {
            title: "Light Control",
        }
    };
};

export default function Content() {
    const [power, setPower] = useState(false);
    const handlePowerChange = () => {
        const req: LightAPIRequest = {
            state: {
                on: !power,
                brightness: 255,
                colored: undefined,
                warmwhite: true,
            }
        };
        // post to api '/api/light'
        fetch('/api/light', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(req)
        }).then(res => {
            res.json().then((data: UpdateLightAPIResponse) => {
                setPower(data.state.on);
            }).catch(err => {
                console.error(err);
            });
        }).catch(err => {
            console.log(err);
        });



    };
    return (
        <Container>
            Light Control
            <Switch checked={power} onChange={handlePowerChange} />
        </Container>
    );
}
