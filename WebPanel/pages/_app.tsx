import { Home as HomeIcon, Laptop as LaptopIcon, Menu as MenuIcon, NetworkCheck as NetworkIcon, NotificationsNone as NotificationIcon, Settings as SettingsIcon } from '@mui/icons-material';
import { AppBar, Box, CssBaseline, Divider, Drawer, IconButton, Link, List, ListItemButton, ListItemIcon, ListItemText, Toolbar, Typography } from '@mui/material';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { ReactElement, useState } from 'react';

const DrawerWidth = 220;

type AppListButtonProps = {
    name: string,
    link: string,
    icon: ReactElement
};

export default function DashboardApp({ Component, pageProps }: AppProps) {
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

    const [mobileDrawerOpen, setMobileDrawerOpen] = useState(false);
    const handleDrawerToggle = () => {
        setMobileDrawerOpen(!mobileDrawerOpen);
    };

    function DrawerContent() {
        return (
            <>
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
            </>
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
                        <IconButton
                            color="inherit"
                            aria-label="open drawer"
                            edge="start"
                            onClick={handleDrawerToggle}
                            sx={{ mr: 2, display: { sm: 'none' } }}>
                            <MenuIcon />
                        </IconButton>
                        <Toolbar>
                            <Typography variant="h6" noWrap>MoodyAPI Dashboard</Typography>
                        </Toolbar>
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
                    <Component {...pageProps}></Component>
                </Box>
            </Box>
        </>
    );
}
