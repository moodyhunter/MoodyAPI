import { Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from "@mui/material";
import { useAtom } from "jotai";
import { useCallback } from "react";
import { dialogPropsAtom, openAtom } from "./atoms/AlertDialogAtom";

export default function AlertDialog() {
    const [dprops] = useAtom(dialogPropsAtom);
    const [open, setOpen] = useAtom(openAtom);

    const handleClose1 = useCallback(() => {
        dprops?.onAction1();
        setOpen(false);
    }, [dprops, setOpen]);
    const handleClose2 = useCallback(() => {
        dprops?.onAction2();
        setOpen(false);
    }, [dprops, setOpen]);

    return (
        <Dialog
            open={open}
            onClose={() => { setOpen(false); }}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description">
            <DialogTitle id="alert-dialog-title">{dprops?.title ?? "Dialog Title"}</DialogTitle>
            <DialogContent>
                <DialogContentText id="alert-dialog-description">
                    {dprops?.title ?? "Message"}
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button color={dprops?.action1Color ?? "inherit"} onClick={handleClose1}>{dprops?.action1 ?? "Action 1"}</Button>
                <Button color={dprops?.action2Color ?? "inherit"} onClick={handleClose2}>{dprops?.action2 ?? "Action 2"}</Button>
            </DialogActions>
        </Dialog>
    );
}
