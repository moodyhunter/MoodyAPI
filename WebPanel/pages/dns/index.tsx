import { Container } from '@mui/material';
import { GetServerSideProps } from 'next';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const getServerSideProps: GetServerSideProps = async (context) => {
    return {
        props: {
            title: "DNS Records"
        }
    };
};

export default function Content() {
    return (
        <Container>
            DNS Record Entries
        </Container>
    );
}

