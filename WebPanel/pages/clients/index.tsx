import { Add as AddIcon, Delete as DeleteIcon, Refresh as RefreshIcon } from '@mui/icons-material';
import { Alert, AlertProps, Box, Button, Container, Divider, IconButton, LinearProgress, Snackbar, styled, Switch, Typography } from '@mui/material';
import { DataGrid, GridCellEditCommitParams, GridColDef, GridRenderCellParams, GridToolbarContainer, GridToolbarExport, useGridApiContext } from '@mui/x-data-grid';
import dayjs from 'dayjs';
import { GetServerSideProps } from 'next';
import { useCallback, useEffect, useState } from 'react';
import { APIClient, ClientAPIResponse, CreateClientAPIResponse, DeleteClientAPIResponse, ListClientsAPIResponse, UpdateClientAPIResponse } from '../../common';
import { EmptyFunction } from '../../components';
import useAlertDialog, { AlertDialogProps } from '../../components/AlertDialog';


export const getServerSideProps: GetServerSideProps = async () => {
    return {
        props: {
            title: "API Clients"
        }
    };
};


function CustomNoRowsOverlay() {
    const StyledGridOverlay = styled('div')(({ theme }) => ({
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100%',
        '& .ant-empty-img-1': {
            fill: theme.palette.mode === 'light' ? '#aeb8c2' : '#262626',
        },
        '& .ant-empty-img-2': {
            fill: theme.palette.mode === 'light' ? '#f5f5f7' : '#595959',
        },
        '& .ant-empty-img-3': {
            fill: theme.palette.mode === 'light' ? '#dce0e6' : '#434343',
        },
        '& .ant-empty-img-4': {
            fill: theme.palette.mode === 'light' ? '#fff' : '#1c1c1c',
        },
        '& .ant-empty-img-5': {
            fillOpacity: theme.palette.mode === 'light' ? '0.8' : '0.08',
            fill: theme.palette.mode === 'light' ? '#f5f5f5' : '#fff',
        },
    }));
    return (
        <StyledGridOverlay>
            <svg width="120" height="100" viewBox="0 0 184 152" aria-hidden focusable="false">
                <g fill="none" fillRule="evenodd">
                    <g transform="translate(24 31.67)">
                        <ellipse className="ant-empty-img-5" cx="67.797" cy="106.89" rx="67.797" ry="12.668" />
                        <path className="ant-empty-img-1" d="M122.034 69.674L98.109 40.229c-1.148-1.386-2.826-2.225-4.593-2.225h-51.44c-1.766 0-3.444.839-4.592 2.225L13.56 69.674v15.383h108.475V69.674z" />
                        <path className="ant-empty-img-2" d="M33.83 0h67.933a4 4 0 0 1 4 4v93.344a4 4 0 0 1-4 4H33.83a4 4 0 0 1-4-4V4a4 4 0 0 1 4-4z" />
                        <path className="ant-empty-img-3" d="M42.678 9.953h50.237a2 2 0 0 1 2 2V36.91a2 2 0 0 1-2 2H42.678a2 2 0 0 1-2-2V11.953a2 2 0 0 1 2-2zM42.94 49.767h49.713a2.262 2.262 0 1 1 0 4.524H42.94a2.262 2.262 0 0 1 0-4.524zM42.94 61.53h49.713a2.262 2.262 0 1 1 0 4.525H42.94a2.262 2.262 0 0 1 0-4.525zM121.813 105.032c-.775 3.071-3.497 5.36-6.735 5.36H20.515c-3.238 0-5.96-2.29-6.734-5.36a7.309 7.309 0 0 1-.222-1.79V69.675h26.318c2.907 0 5.25 2.448 5.25 5.42v.04c0 2.971 2.37 5.37 5.277 5.37h34.785c2.907 0 5.277-2.421 5.277-5.393V75.1c0-2.972 2.343-5.426 5.25-5.426h26.318v33.569c0 .617-.077 1.216-.221 1.789z" />
                    </g>
                    <path className="ant-empty-img-3" d="M149.121 33.292l-6.83 2.65a1 1 0 0 1-1.317-1.23l1.937-6.207c-2.589-2.944-4.109-6.534-4.109-10.408C138.802 8.102 148.92 0 161.402 0 173.881 0 184 8.102 184 18.097c0 9.995-10.118 18.097-22.599 18.097-4.528 0-8.744-1.066-12.28-2.902z" />
                    <g className="ant-empty-img-4" transform="translate(149.65 15.383)">
                        <ellipse cx="20.654" cy="3.167" rx="2.849" ry="2.815" />
                        <path d="M5.698 5.63H0L2.898.704zM9.259.704h4.985V5.63H9.259z" />
                    </g>
                </g>
            </svg>
            <Box sx={{ mt: 1 }}>No Rows</Box>
        </StyledGridOverlay>
    );
}

async function getClientsAsync(): Promise<ClientAPIResponse<ListClientsAPIResponse>> {
    const res = await fetch('/api/clients');
    return await res.json();
}

async function updateClientsAsync(params: GridCellEditCommitParams): Promise<ClientAPIResponse<UpdateClientAPIResponse>> {
    const resp = await fetch('/api/clients', {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            id: params.id as number,
            [params.field]: params.value,
        })
    });
    return await resp.json();
}

async function deleteClientsAsync(clientId: number): Promise<ClientAPIResponse<DeleteClientAPIResponse>> {
    const resp = await fetch('/api/clients', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: clientId })
    });
    return await resp.json();
}

async function createClientsAsync(client: APIClient): Promise<ClientAPIResponse<CreateClientAPIResponse>> {
    const resp = await fetch('/api/clients', {
        method: "POST",
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(client)
    });
    return await resp.json();
}

export default function Content() {
    const [snackbarState, setSnackbarState] = useState<AlertProps | null>(null);
    const [rows, setRows] = useState<APIClient[]>([]);
    const [loading, setLoading] = useState(true);

    const [rowDisabled, setRowDisabled] = useState<{ [a: string]: boolean }>({});
    const [newClientButtonDisabled, setNewClientButtonDisabled] = useState<boolean>(false);

    const handleCloseSnackbar = () => setSnackbarState(null);
    const handleSuccessMessage = (msg: string) => setSnackbarState({ children: msg, severity: 'success' });
    const handleErrorMessage = (msg: string) => setSnackbarState({ children: msg, severity: 'error' });

    const updateClient = useCallback(
        async (params: GridCellEditCommitParams) => {
            setRowDisabled((prev) => { prev[params.id] = true; return prev; });
            const { success, message, data } = await updateClientsAsync(params);
            setRowDisabled((prev) => { prev[params.id] = false; return prev; });
            if (success && data) {
                const client = data.client;
                setRows((prev) => prev.map((row) => (row.id == params.id ? { ...row, ...client } : row)));
                handleSuccessMessage(`Successfully updated client ${params.id}`);
            } else {
                setRows((prev) => [...prev]);
                handleErrorMessage(`Failed to update client ${params.id}: ${message}`);
            }
        }, []
    );

    const refreshClients = () => {
        setLoading(true);
        setRows([]);
        getClientsAsync().then(({ success, message, data }) => {
            if (success && data) {
                setRows(data.clients);
                handleSuccessMessage(`Successfully loaded ${data.clients.length} client(s)`);
                setLoading(false);
            } else {
                setRows([]);
                handleErrorMessage(`Failed to list clients: ${message}`);
                setLoading(false);
            }
        });
    };

    useEffect(refreshClients, []);

    const { OpenDialog, AlertDialog } = useAlertDialog();

    const columns: GridColDef[] = [
        { hideable: false, editable: false, width: 50, align: 'center', headerAlign: 'center', field: 'id', headerName: 'ID', sortable: true },
        { hideable: false, editable: true, width: 200, align: 'left', headerAlign: 'left', field: 'name', headerName: 'Client Name' },
        {
            hideable: true, editable: false, width: 350, align: 'center', headerAlign: 'center', field: 'uuid', headerName: 'Client ID',
            renderCell: (params) => (<code>{(params.row as APIClient).uuid}</code>)
        },
        {
            hideable: false, editable: false, width: 150, align: 'center', headerAlign: 'center',
            field: 'lastSeen', headerName: 'Last Seen', sortingOrder: ['asc', 'desc'],
            renderCell: (params) => (<div>{(params.row as APIClient).lastSeen ? (dayjs(params.row.lastSeen).format('YYYY/MM/DD HH:mm:ss')) : "N/A"}</div>)
        },
        {
            hideable: false, editable: false, width: 100, align: 'center', headerAlign: 'center', sortable: false, filterable: false,
            field: 'enabled', headerName: "Enabled",
            renderCell: function (params: GridRenderCellParams<boolean>) {
                const { id, field } = params;
                // eslint-disable-next-line react-hooks/rules-of-hooks
                const apiContext = useGridApiContext();

                const handleChange = async (event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => {
                    apiContext.current.setEditCellValue({ id, field, value: checked }, event);
                    const isValid = await apiContext.current.commitCellChange({ id, field });
                    if (isValid) {
                        apiContext.current.setCellMode(id, field, 'view');
                    }
                };

                return (<Switch disabled={rowDisabled[params.id]} color="success" checked={!!(params.row as APIClient).enabled} onChange={handleChange} />);
            }
        },
        {
            hideable: false, editable: false, width: 150, align: 'center', headerAlign: 'center', sortable: false, filterable: false,
            field: '_actions', headerName: 'Actions',
            renderCell: (params) => {
                // eslint-disable-next-line react-hooks/rules-of-hooks
                const apiContext = useGridApiContext();
                const { id } = params;
                const field = 'uuid';

                const resetUuidDialogProps: AlertDialogProps = {
                    title: `Reset UUID for '${(params.row as APIClient).name ?? "<unknown>"}'?`,
                    message: "This action will go into effect immediately! The action cannot be reverted!",
                    action1: "Cancel",
                    action2: `Reset UUID`,
                    action1Color: "inherit",
                    action2Color: "error",
                    onAction1: (EmptyFunction),
                    onAction2: (async () => {
                        const newuuid = crypto.randomUUID ? crypto.randomUUID() : "invalid";
                        apiContext.current.setEditCellValue({ id, field, value: newuuid });
                        const isValid = await apiContext.current.commitCellChange({ id, field });
                        if (isValid) {
                            apiContext.current.setCellMode(id, field, 'view');
                        }
                    }),
                };

                const deleteClientDialogProps: AlertDialogProps = {
                    title: `Delete Client '${(params.row as APIClient).name ?? "<unknown>"}'?`,
                    message: "This action will go into effect immediately! The action cannot be reverted!",
                    action1: "Cancel",
                    action2: `Delete Client`,
                    action1Color: "inherit",
                    action2Color: "error",
                    onAction1: (EmptyFunction),
                    onAction2: (async () => {
                        setRowDisabled((prev) => { prev[params.id] = true; return prev; });
                        const { success, message, data } = await deleteClientsAsync(params.id as number);
                        setRowDisabled((prev) => { prev[params.id] = false; return prev; });

                        if (success && data && data.deleted) {
                            handleSuccessMessage(`Successfully deleted client ${params.id}`);
                            setRows((r) => r.filter((a) => a.id !== params.id));
                        }
                        else {
                            handleErrorMessage(`Failed to delete client ${params.id}: ${message}`);
                        }
                    }),
                };

                return (<>
                    <IconButton disabled={rowDisabled[params.id]} color="success" onClick={() => { OpenDialog(resetUuidDialogProps); }}><RefreshIcon /></IconButton>
                    <IconButton disabled={rowDisabled[params.id]} color="error" onClick={() => { OpenDialog(deleteClientDialogProps); }}><DeleteIcon /></IconButton>
                </>);
            }
        }
    ];

    const CustomToolbar = () => {
        const createClient = async () => {
            setNewClientButtonDisabled(true);
            const newClient: APIClient = { id: 0, name: "Client " + rows.length.toString() };
            const { success, message, data } = await createClientsAsync(newClient);
            setNewClientButtonDisabled(false);

            if (success && data) {
                setRows((r) => r.concat([data.client]));
            } else {
                handleErrorMessage("Failed to create new client: " + message);
            }
        };

        return (
            <GridToolbarContainer>
                <Button disabled={newClientButtonDisabled} size="small" startIcon={<AddIcon />} onClick={createClient}>New Client</Button>
                <Divider orientation='vertical' flexItem sx={{ margin: '5px' }} />
                <Button size="small" startIcon={<RefreshIcon />} onClick={refreshClients}>Refresh</Button>
                <Divider orientation='vertical' flexItem sx={{ margin: '5px' }} />
                <GridToolbarExport />
            </GridToolbarContainer >
        );
    };

    return (
        <Container style={{ height: '62vh', width: '100%' }}>
            <Typography variant='h4'>API Clients</Typography>
            <br />
            <DataGrid
                components={{
                    LoadingOverlay: LinearProgress,
                    NoRowsOverlay: CustomNoRowsOverlay,
                    Toolbar: CustomToolbar,
                }}
                rows={rows}
                columns={columns}
                loading={loading}
                onCellEditCommit={updateClient} />
            {!!snackbarState && (
                <Snackbar open onClose={handleCloseSnackbar} autoHideDuration={6000}>
                    <Alert {...snackbarState} onClose={handleCloseSnackbar} />
                </Snackbar>
            )}
            {AlertDialog}
        </Container>
    );
}
