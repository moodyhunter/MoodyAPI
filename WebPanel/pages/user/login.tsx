import { AccountCircleOutlined, KeyOutlined, Visibility, VisibilityOff } from "@mui/icons-material";
import { Alert, Box, Button, Container, FormControl, IconButton, InputAdornment, TextField, Typography } from "@mui/material";
import { GetServerSideProps } from "next";
import { getCsrfToken, useSession } from 'next-auth/react';
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

type SigninProperty = {
    csrfToken: string;
    errorReason: string;
};

export const getServerSideProps: GetServerSideProps = async (context) => {
    return {
        props: {
            csrfToken: await getCsrfToken(context) ?? "",
            errorReason: context.query["error"] ?? "",
        }
    };
};


export default function SignIn({ csrfToken, errorReason }: SigninProperty) {
    const [showPassword, setShowPassword] = useState(false);
    const session = useSession();
    const router = useRouter();

    const handleToggleShowPassword = () => { setShowPassword(!showPassword); };

    useEffect(() => {
        if (session.status === "authenticated") {
            setTimeout(() => router.push('/'), 1000);
        }
    });

    if (session.status === "authenticated") {
        return (
            <Container maxWidth="sm" sx={{ mt: '4vh' }}>
                <Alert severity="info">You have successfully logged in, redirecting...</Alert>
            </Container>
        );
    }

    return (
        <Container maxWidth="sm" sx={{ mt: "4vh" }}>
            <Typography variant="h5">Login With Your Username and Password</Typography>
            <br />
            {errorReason == "CredentialsSignin" && <><Alert severity="warning">Incorrect Username or Password.</Alert><br /></>}

            <form method="post" action="/api/auth/callback/credentials">
                <Box sx={{ display: 'grid', rowGap: 2 }}>
                    <input name="csrfToken" type="hidden" defaultValue={csrfToken} />

                    <FormControl >
                        <TextField label="username" variant="standard" id="login_username" name="username" InputProps={{ startAdornment: <AccountCircleOutlined sx={{ marginRight: 1 }} /> }} />
                    </FormControl>

                    <FormControl>
                        <TextField label="Password" variant="standard" name="password" type={showPassword ? "text" : "password"}
                            InputProps={{
                                startAdornment: <KeyOutlined sx={{ marginRight: 1 }} />,
                                endAdornment:
                                    <InputAdornment position="end">
                                        <IconButton onClick={handleToggleShowPassword} onMouseDown={(event) => event.preventDefault()}>
                                            {showPassword ? <VisibilityOff /> : <Visibility />}
                                        </IconButton>
                                    </InputAdornment>
                            }}
                        />
                    </FormControl>

                    <Button type="submit" variant="contained">Continue</Button>
                </Box>
            </form>
        </Container>
    );
}
