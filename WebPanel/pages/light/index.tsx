import { Container, Switch } from '@mui/material';
import { GetServerSideProps } from 'next';
import { useState } from 'react';
import { ClientAPIResponse, UpdateLightAPIResponse } from '../../common';

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
        // post to api '/api/light'
        fetch('/api/light', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ power: power })
        }).then(res => {
            const status = res as unknown as ClientAPIResponse<UpdateLightAPIResponse>;
            if (status.success) {
                setPower(status.data?.state.on ?? false);
            }
            else {
                console.log("Error: " + status.message);
            }
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
