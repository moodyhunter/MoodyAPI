import { Home as HomeIcon, Laptop as LaptopIcon, NetworkCheck as NetworkIcon, NotificationsNone as NotificationIcon, Settings as SettingsIcon } from '@mui/icons-material';
import { AppBar, Box, Container, CssBaseline, Divider, Drawer, Link, List, ListItemButton, ListItemIcon, ListItemText, Toolbar, Typography } from '@mui/material';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { ReactElement } from 'react';

const DrawerWidth = 220;

type AppListButtonProps = {
    name: string,
    link: string,
    icon: ReactElement
};

export default function MyApp({ Component, pageProps }: AppProps) {
    const router = useRouter();

    function AppListButton(props: AppListButtonProps) {
        return (
            <Link href={props.link} underline='none'>
                <ListItemButton selected={router.route === props.link} key={props.name} sx={{ minHeight: 48, justifyContent: 'initial', px: 2.5 }}>
                    <ListItemIcon sx={{ minWidth: 0, mr: 3, justifyContent: 'center' }}>
                        {props.icon}
                    </ListItemIcon>
                    <ListItemText primary={props.name} sx={{ opacity: 1 }} />
                </ListItemButton>
            </Link>
        );
    }

    return (
        <>
            <Head>
                {/* https://stackoverflow.com/a/19903063/16018952 */}
                <title>{[pageProps.title, "MoodyAPI Dashboard"].filter(Boolean).join(" — ")}</title>
                <link rel="icon" href="/favicon.ico" />
            </Head>
            <CssBaseline />
            <Box sx={{ display: 'flex' }}>
                <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
                    <Toolbar>
                        <Typography variant="h6" noWrap>MoodyAPI Dashboard</Typography>
                    </Toolbar>
                </AppBar>

                <Drawer variant="permanent" sx={{ width: DrawerWidth }}>
                    <Toolbar />
                    <List>
                        <AppListButton link='/' name='Home' icon={(<HomeIcon />)} />
                        <AppListButton link='/clients' name='API Clients' icon={(<LaptopIcon />)} />
                        <AppListButton link='/wg' name='Wireguard Clients' icon={(<NetworkIcon />)} />
                        <AppListButton link='/notifications' name='Notifications' icon={(<NotificationIcon />)} />
                    </List>
                    <Divider />
                    <List>
                        <AppListButton link='/settings' name='Settings' icon={(<SettingsIcon />)} />
                    </List>
                </Drawer>

                <Container>
                    <Toolbar />
                    <Component {...pageProps}></Component>
                </Container>
            </Box>
        </>
    );
}
