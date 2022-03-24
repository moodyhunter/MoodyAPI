import { Container } from '@mui/material';
import { GetServerSideProps } from 'next';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const getServerSideProps: GetServerSideProps = async (context) => {
    return {
        props: {
            title: "Home"
        }
    };
};

export default function Content() {
    return (
        <Container>
            Settings Page
        </Container>
    );
}

