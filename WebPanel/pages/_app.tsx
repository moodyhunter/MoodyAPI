import {
    Home as HomeIcon,
    Laptop as LaptopIcon,
    Menu as MenuIcon,
    Settings as SettingsIcon,
    NotificationsNone as NotificationIcon,
    NetworkCheck as NetworkIcon
} from '@mui/icons-material';
import { Box, Container, CssBaseline, Divider, IconButton, List, ListItemButton, ListItemIcon, ListItemText, Toolbar, Typography } from '@mui/material';
import { AppProps } from 'next/app';
import { ReactElement, useEffect, useState } from 'react';
import { AppDrawer, AppTopBar, DrawerHeader } from '../components';

type AppListButtonProps = {
    open: boolean,
    name: string,
    icon: ReactElement
};

function AppListButton(props: AppListButtonProps) {
    return (
        <ListItemButton
            key={props.name} sx={{ minHeight: 48, justifyContent: props.open ? 'initial' : 'center', px: 2.5 }}>
            <ListItemIcon sx={{ minWidth: 0, mr: props.open ? 3 : 'auto', justifyContent: 'center' }}>
                {props.icon}
            </ListItemIcon>
            <ListItemText primary={props.name} sx={{ opacity: props.open ? 1 : 0 }} />
        </ListItemButton>
    );
}

export default function MyApp({ Component, pageProps }: AppProps) {
    const [open, setOpen] = useState(false);

    useEffect(() => { setOpen(JSON.parse(window.localStorage.getItem('open') ?? "false")); }, []);
    useEffect(() => { window.localStorage.setItem('open', String(open)); }, [open]);

    const toggleDrawer = () => { setOpen(!open); };
    return (
        <Box sx={{ display: 'flex' }}>
            <CssBaseline />
            <AppTopBar position="fixed" open={open} sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
                <Toolbar>
                    <IconButton
                        color="inherit"
                        onClick={toggleDrawer}
                        edge="start"
                        sx={{ marginRight: 5 }}>
                        <MenuIcon />
                    </IconButton>
                    <Typography variant="h6" noWrap component="div">
                        MoodyAPI Dashboard
                    </Typography>
                </Toolbar>
            </AppTopBar>
            <AppDrawer variant="permanent" open={open}>
                <DrawerHeader />
                <List>
                    <AppListButton name='Home' open={open} icon={(<HomeIcon />)} />
                    <AppListButton name='API Clients' open={open} icon={(<LaptopIcon />)} />
                    <AppListButton name='Wireguard Clients' open={open} icon={(<NetworkIcon />)} />
                    <AppListButton name='Notifications' open={open} icon={(<NotificationIcon />)} />
                </List>
                <Divider />
                <List>
                    <AppListButton name='Settings' open={open} icon={(<SettingsIcon />)} />
                </List>
            </AppDrawer>

            <Container>
                <Component {...pageProps}></Component>
            </Container>
        </Box >
    );
}
