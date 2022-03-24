import { Container } from '@mui/material';
import { GetServerSideProps } from 'next';

export const getServerSideProps: GetServerSideProps = async () => {
    return {
        props: {
            title: "Home"
        }
    };
};

export default function Home() {
    return (
        <Container>
            Home
        </Container>
    );
}
