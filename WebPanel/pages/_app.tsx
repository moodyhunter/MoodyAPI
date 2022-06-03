import { Dashboard as DashboardIcon, Home as HomeIcon, Laptop as LaptopIcon, Logout as LogoutIcon, Menu as MenuIcon, NetworkCheck as NetworkIcon, NotificationsNone as NotificationIcon, Settings as SettingsIcon } from '@mui/icons-material';
import { AppBar, Box, Button, CssBaseline, Divider, Drawer, IconButton, Link, List, ListItemButton, ListItemIcon, ListItemText, Toolbar, Typography } from '@mui/material';
import { SessionProvider, signIn, signOut, useSession } from 'next-auth/react';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { memo, ReactElement, useEffect, useMemo, useState } from 'react';
import { LoadingScreen } from '../components';
import NextNProgress from "nextjs-progressbar";
import NextLink from "next/link";

const DrawerWidth = 220;

type AppListButtonProps = {
    name: string,
    link: string,
    icon: ReactElement
};

const UnmemoizedAppListButton = (props: AppListButtonProps) => {
    const router = useRouter();

    return (
        <NextLink passHref href={props.link}>
            <Link underline='none'>
                <ListItemButton selected={router.route === props.link} key={props.name} sx={useMemo(() => ({ minHeight: 48, justifyContent: 'initial', px: 2.5 }), [])}>
                    <ListItemIcon sx={{ minWidth: 0, mr: 3, justifyContent: 'center' }}>
                        {props.icon}
                    </ListItemIcon>
                    <ListItemText primary={props.name} sx={{ opacity: 1 }} />
                </ListItemButton>
            </Link>
        </NextLink>
    );
};

const AppListButton = memo(UnmemoizedAppListButton);

const DrawerContent = () => {
    return (
        <>
            <Toolbar />
            <List>
                <AppListButton link='/' name='Home' icon={(<HomeIcon />)} />
                <AppListButton link='/clients' name='API Clients' icon={(<LaptopIcon />)} />
                <AppListButton link='/notifications' name='Notifications' icon={(<NotificationIcon />)} />
                <AppListButton link='/dns' name='DNS Records' icon={(<NetworkIcon />)} />
                <AppListButton link='/status' name='Status' icon={(<DashboardIcon />)} />
            </List>
            <Divider />
            <List>
                <AppListButton link='/settings' name='Settings' icon={(<SettingsIcon />)} />
            </List>
        </>
    );
};

const AuthenticatePageContent = ({ Component, pageProps: { ...pageProps } }: AppProps) => {
    const session = useSession();
    const router = useRouter();

    useEffect(() => {
        if (router.route !== "/user/login" && session.status === "unauthenticated") {
            signIn();
        }
    });

    if (router.route === "/user/login") {
        return (<Component {...pageProps}></Component>);
    }

    if (session.status === "loading") {
        return (<LoadingScreen />);
    } else if (session.status === "unauthenticated") {
        return (<Typography>Unauthenticated</Typography>);
    } else if (session.status === "authenticated") {
        return (<Component {...pageProps}></Component>);
    } else {
        return (<div>SERVER ERROR.</div>);
    }
};


const AppFrame = (appProps: AppProps) => {
    const [mobileDrawerOpen, setMobileDrawerOpen] = useState(false);
    const handleDrawerToggle = () => {
        setMobileDrawerOpen(!mobileDrawerOpen);
    };

    const session = useSession();
    const handleSignout = async () => { await signOut(); };


    return (<>
        <Head>
            {/* https://stackoverflow.com/a/19903063/16018952 */}
            <title>{[appProps.pageProps.title, "MoodyAPI Dashboard"].filter(Boolean).join(" — ")}</title>
            <link rel="icon" href="/favicon.ico" />
        </Head>
        <CssBaseline />
        <Box sx={{ display: 'flex' }}>
            {session.status === "authenticated" && <NextNProgress color="#29D" showOnShallow={true} />}
            <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
                <Toolbar>
                    <IconButton color="inherit" edge="start" onClick={handleDrawerToggle} sx={{ mr: 2, display: { sm: 'none' } }}>
                        <MenuIcon />
                    </IconButton>
                    <Typography variant="h6" noWrap sx={{ flexGrow: 1 }}>MoodyAPI Dashboard</Typography>
                    {session.status === "authenticated" && <Button color="inherit" startIcon={<LogoutIcon />} onClick={handleSignout}>Log out</Button>}
                </Toolbar>
            </AppBar>

            <Box component="nav" sx={{ width: { sm: DrawerWidth } }}>
                <Drawer
                    variant="temporary"
                    sx={{ display: { xs: 'block', sm: 'none' }, '& .MuiDrawer-paper': { boxSizing: 'border-box', width: DrawerWidth } }}
                    open={mobileDrawerOpen}
                    onClose={handleDrawerToggle}
                    ModalProps={{ keepMounted: true }}>
                    <DrawerContent />
                </Drawer>
                <Drawer
                    variant="permanent"
                    sx={{ display: { xs: 'none', sm: 'block' }, '& .MuiDrawer-paper': { boxSizing: 'border-box', width: DrawerWidth } }}>
                    <DrawerContent />
                </Drawer>
            </Box>

            <Box component="main" sx={{ flexGrow: 1, p: 3, width: { sm: `calc(100% - ${DrawerWidth}px)` } }}>
                <Toolbar />
                <AuthenticatePageContent {...appProps} />
            </Box>
        </Box>
    </>);
};

export default function DashboardApp(props: AppProps) {
    return (<>
        <SessionProvider>
            <AppFrame {...props} />
        </SessionProvider>
    </>);
}
