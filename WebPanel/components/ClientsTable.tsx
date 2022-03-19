
import { Alert, AlertProps, Container, Snackbar } from '@mui/material';
import { DataGrid, GridCellEditCommitParams, GridColDef } from '@mui/x-data-grid';
import { useCallback, useState } from 'react';
import { APIClient } from '../common';
import { TestData } from './demo';

const columns: GridColDef[] = [
    {
        hideable: false, resizable: false, editable: false,
        field: 'id', headerName: 'ID', width: 50, sortingOrder: ['asc'],
    },
    {
        hideable: false, resizable: true, editable: true,
        field: 'clientName', headerName: 'Name', width: 200
    },
    {
        hideable: false, resizable: true, editable: true,
        field: 'clientUuid', headerName: 'UUID', width: 300
    },
    {
        hideable: false, resizable: true, editable: false,
        field: 'lastSeen', headerName: 'Last Seen', width: 300,
    },
];


export function ClientsTable() {
    const [snackbarState, setSnackbarState] = useState<Pick<AlertProps, 'children' | 'severity'> | null>(null);
    const [rows, setRows] = useState<APIClient[]>(TestData);
    const handleCloseSnackbar = () => setSnackbarState(null);

    const updateData = useCallback(
        async (params: GridCellEditCommitParams) => {
            try {
                const data = await (await fetch('/api/clients', {
                    method: 'PATCH',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        id: params.id as number,
                        [params.field]: params.value,
                    })
                })).json();

                setRows((prev) => prev.map((row) => (row.id === (params.id as number) ? { ...row, ...data } : row)));
                setSnackbarState({ children: 'Client updated', severity: 'success' });
            } catch (error) {
                console.log(error);
                setSnackbarState({ children: 'Error while updating client', severity: 'error' });
                setRows((prev) => [...prev]);
            }
        },
        [],
    );

    return (
        <Container style={{ height: '60vh', width: '100%' }}>
            <DataGrid rows={rows} columns={columns} onCellEditCommit={updateData} />
            {!!snackbarState && (
                <Snackbar open onClose={handleCloseSnackbar} autoHideDuration={6000}>
                    <Alert {...snackbarState} onClose={handleCloseSnackbar} />
                </Snackbar>
            )}
        </Container>
    )
}