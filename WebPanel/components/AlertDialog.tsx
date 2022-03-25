import { Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from "@mui/material";
import { useCallback, useState } from "react";
import { EmptyFunction } from ".";

type _Color = 'inherit' | 'primary' | 'secondary' | 'success' | 'error' | 'info' | 'warning';

export type AlertDialogProps = {
    title: string,
    message: string,
    onAction1: () => void,
    onAction2: () => void,
    action1: string,
    action2: string,
    action1Color: _Color,
    action2Color: _Color
};

const DefaultAlertDialogProps: AlertDialogProps = {
    action1: "action 1",
    action2: "action 2",
    action1Color: "success",
    action2Color: "error",
    message: "Message",
    onAction1: EmptyFunction,
    onAction2: EmptyFunction,
    title: "Title"
};

export default function useAlertDialog() {
    const [open, setOpen] = useState(false);
    const [dprops, setDProps] = useState<AlertDialogProps>(DefaultAlertDialogProps);

    const OpenDialog = useCallback((props: AlertDialogProps) => { setDProps(props); setOpen(true); }, []);
    const CloseDialog = useCallback(() => { setOpen(false); }, []);

    const handleClose1 = useCallback(() => { setOpen(false); dprops.onAction1(); }, [dprops]);
    const handleClose2 = useCallback(() => { setOpen(false); dprops.onAction2(); }, [dprops]);

    const AlertDialog = (
        <Dialog
            open={open}
            onClose={() => { setOpen(false); }}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description">
            <DialogTitle id="alert-dialog-title">{dprops.title}</DialogTitle>
            <DialogContent>
                <DialogContentText id="alert-dialog-description">
                    {dprops.message}
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button color={dprops.action1Color} onClick={handleClose1}>{dprops.action1}</Button>
                <Button color={dprops.action2Color} onClick={handleClose2}>{dprops.action2}</Button>
            </DialogActions>
        </Dialog>
    );

    return { OpenDialog, CloseDialog, AlertDialog };
}
