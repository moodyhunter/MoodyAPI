import { Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from "@mui/material";
import { useState } from "react";

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

export default function AlertDialog() {
    const [open, setOpen] = useState(false);
    const [dprops, setDProps] = useState<AlertDialogProps | null>(null);

    const handleClose1 = () => { setOpen(false); if (dprops) dprops.onAction1(); };
    const handleClose2 = () => { setOpen(false); if (dprops) dprops.onAction2(); };

    const OpenDialog = (props: AlertDialogProps) => {
        setDProps(props);
        setOpen(true);
    };

    const CloseDialog = () => {
        setOpen(false);
        setDProps(null);
    };

    const AlertDialogComponent = (
        <Dialog
            open={open}
            onClose={() => { setOpen(false); }}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description">
            <DialogTitle id="alert-dialog-title">{dprops?.title ?? "Dialog Title"}</DialogTitle>
            <DialogContent>
                <DialogContentText id="alert-dialog-description">
                    {dprops?.message ?? "Content"}
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button color={dprops?.action1Color ?? "inherit"} onClick={handleClose1}>{dprops?.action1 ?? "Action 1"}</Button>
                <Button color={dprops?.action2Color ?? "inherit"} onClick={handleClose2}>{dprops?.action2 ?? "Action 2"}</Button>
            </DialogActions>
        </Dialog>
    );

    return { OpenDialog, CloseDialog, AlertDialogComponent };
}
