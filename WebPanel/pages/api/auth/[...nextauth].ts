import { compareSync } from "bcrypt";
import NextAuth from "next-auth/next";
import CredentialsProvider from "next-auth/providers/credentials";

const CredentialDescriptor = {
    username: { label: "Username", type: "text", placeholder: "Your Username" },
    password: { label: "Password", type: "password", placeholder: "Your Password" },
};

// Need env:NEXTAUTH_SECRET

export default NextAuth({
    pages: {
        signIn: '/user/login',
        error: `/user/login`,
        signOut: '/user/logout'
    },
    secret: process.env["NEXTAUTH_SECRET"],
    jwt: { secret: process.env["NEXTAUTH_SECRET"] },
    providers: [
        CredentialsProvider({
            name: "credentials",
            type: "credentials",
            credentials: CredentialDescriptor,
            authorize: async (cred, req) => {
                // unused
                req;

                // admin:admin
                const panelUser = process.env["PANEL_USERNAME"] ?? "admin";
                const panelPasswordHash = process.env["PANEL_PASSWORD_HASH"] ?? "$2a$10$ovuekWbiVDbnhpPWHLwQiexXz2sshxdPEZ1inh3S017RibzvtARN.";

                if (!cred)
                    return null;

                console.log(`User '${cred.username}' is trying to login.`);

                if ((cred.username ?? "_invalid_") != panelUser)
                    return null;

                if (!compareSync(cred.password, panelPasswordHash)) {
                    return null;
                }

                console.log(`User ${cred.username} logged in successfully.`);
                return { id: 1, name: cred.username, email: "" };
            }
        })
    ]
});
