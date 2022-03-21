import { Container } from '@mui/material';
import { GetServerSideProps } from 'next';

export default function Home() {
    return (
        <Container>
        </Container>
    );
}

export const getServerSideProps: GetServerSideProps = async () => {
    return {
        props: {
            title: "Home"
        }
    };
};
