import { atom } from "jotai";
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

export const openAtom = atom(false);
export const dialogPropsAtom = atom<AlertDialogProps | null>(null);
export const openDialogAtom = atom<null, AlertDialogProps>(
    null,
    (_, set, by) => {
        set(dialogPropsAtom, by);
        set(openAtom, true);
    }
);

export const closeDialogAtom = atom<null, never>(null, (_, set) => {
    set(dialogPropsAtom, null);
    set(openAtom, false);
});
