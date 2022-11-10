import { Container, Slider, Switch, Typography } from '@mui/material';
import { GetServerSideProps } from 'next';
import { useState } from 'react';
import { ClientAPIResponse, createAuth, getServerConnection, LightAPIRequest, UpdateLightAPIResponse } from '../../common';
import { LightState } from '../../common/protos/light/light';

export const getServerSideProps: GetServerSideProps = async () => {
    const client = getServerConnection();
    const AuthObject = createAuth();
    const resp = await client.getLightState({ auth: AuthObject });

    if (resp.state === undefined)
        resp.state = {} as LightState;

    if (resp.state.brightness === undefined)
        resp.state.brightness = 0;

    if (resp.state.on === undefined)
        resp.state.on = false;

    if (resp.state.colored === undefined)
        resp.state.colored = { red: 0, green: 0, blue: 0 };

    if (resp.state.warmwhite === undefined)
        resp.state.warmwhite = true;

    return {
        props: {
            title: "Light Control",
            lightState: resp.state,
        }
    };
};

export default function Content(props: { title: string, lightState: LightState }) {
    const state = props.lightState;

    const [power, setPower] = useState(state.on);
    const [brightness, setBrightness] = useState(state.brightness);

    function doUpdate(req: LightAPIRequest) {
        fetch('/api/light', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req)
        }).then(res => {
            return res.json();
        }).then((data: ClientAPIResponse<UpdateLightAPIResponse>) => {
            console.log(data);
            if (data.success && data.data?.state) {
                const state = data.data.state;
                setPower(state.on);
                setBrightness(state.brightness);
            }
        }).catch(err => {
            console.error(err);
        });
    }

    const handlePowerChange = () => {
        const req: LightAPIRequest = {
            state: {
                on: !power,
                brightness: brightness,
                colored: undefined,
                warmwhite: true,
            }
        };

        doUpdate(req);
    };

    const handleBrightnessChange = (_event: Event, newValue: number | number[]) => {
        const req: LightAPIRequest = {
            state: {
                on: power,
                brightness: newValue as number,
                colored: undefined,
                warmwhite: true,
            }
        };

        doUpdate(req);
    };


    return (
        <Container>
            <Typography variant='h4'>Light</Typography>
            Power
            <Switch checked={power} onChange={handlePowerChange} />
            <br />
            Brightness
            <Slider defaultValue={0} value={brightness} onChange={handleBrightnessChange} min={0} max={255} aria-label="Disabled slider" />
        </Container>
    );
}
