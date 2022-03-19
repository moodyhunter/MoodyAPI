import { Container, Toolbar, Typography } from '@mui/material';
import { GetServerSideProps } from 'next';
import { ClientsTable } from '../../components';

export default function Home() {
    return (
        <Container>
            <br />
            <Toolbar>
                <Typography variant='h4'>API Clients</Typography>
            </Toolbar>
            <ClientsTable></ClientsTable>
        </Container>
    )
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const getServerSideProps: GetServerSideProps = async (context) => {
    return {
        props: {
            title: "API Clients"
        }
    }
}