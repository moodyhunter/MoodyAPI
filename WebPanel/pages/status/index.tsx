import { ExpandMore as ExpandMoreIcon } from '@mui/icons-material';
import { Accordion, AccordionDetails, AccordionSummary, Container, Typography } from '@mui/material';
import { GetServerSideProps } from 'next';
import Image from 'next/image';
import { SyntheticEvent, useState } from 'react';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const getServerSideProps: GetServerSideProps = async (context) => {
    return {
        props: {
            title: "Server Status"
        }
    };
};

export default function Content() {

    const [expanded, setExpanded] = useState<string | false>(false);

    const handleChange =
        (panel: string) => (event: SyntheticEvent, isExpanded: boolean) => {
            setExpanded(isExpanded ? panel : false);
        };

    return (
        <Container sx={{ height: '100vh' }}>
            <Accordion expanded={expanded === 'server-status'} onChange={handleChange('server-status')}>
                <AccordionSummary expandIcon={<ExpandMoreIcon />} >
                    <Typography>Server Status</Typography>
                </AccordionSummary>
                <AccordionDetails>
                    <Image alt='Server Load' src='https://dash.mooody.me/.image/server_load' width={800} height={500}></Image>
                </AccordionDetails>
            </Accordion>
            <Accordion expanded={expanded === 'rpi-status'} onChange={handleChange('rpi-status')}>
                <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                    <Typography>Raspberry Pi Status</Typography>
                </AccordionSummary>
                <AccordionDetails>
                    <Image alt='RPi System Load' src='https://dash.mooody.me/.image/rpi_load' width={800} height={500}></Image>
                </AccordionDetails>
            </Accordion>
        </Container>
    );
}

