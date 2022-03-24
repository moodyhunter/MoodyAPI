import { Container } from '@mui/material';
import { GetServerSideProps } from 'next';

export const getServerSideProps: GetServerSideProps = async () => {
    return {
        props: {
            title: "Notifications"
        }
    };
};

export default function Content() {
    return (
        <Container>
            Notifications
        </Container>
    );
}
