import { Checkbox, Container, Grid, Slider, Switch, Typography } from '@mui/material';
import { GetServerSideProps } from 'next';
import { useState } from 'react';
import { ClientAPIResponse, createAuth, getServerConnection, UpdateLightAPIResponse } from '../../common';
import { LightState } from '../../common/protos/light/light';
import { ColorResult, HuePicker, RgbColor } from '@hello-pangea/color-picker';

type Props = {
    direction?: "horizontal" | "vertical";
};

function SliderPointer({ direction }: Props) {
    const styles: Record<string, React.CSSProperties> = {
        picker: {
            width: "18px",
            height: "18px",
            borderRadius: "50%",
            transform:
                direction === "vertical"
                    ? "translate(-3px, -9px)"
                    : "translate(-9px, -1px)",
            backgroundColor: "rgb(248, 248, 248)",
            boxShadow: "0 1px 4px 0 rgba(0, 0, 0, 0.37)",
        },
    };

    return <div style={styles.picker} />;
}

function ColorBox({ powered, warmwhite, color, brightness }: { powered: boolean, warmwhite: boolean, color: RgbColor, brightness: number }) {
    const r = powered ? (warmwhite ? 255 : color.r) : 0;
    const g = powered ? (warmwhite ? 229 : color.g) : 0;
    const b = powered ? (warmwhite ? 167 : color.b) : 0;

    const styles: Record<string, React.CSSProperties> = {
        colorBox: {
            width: 50,
            height: 50,
            backgroundColor: `rgb(${r * brightness / 255}, ${g * brightness / 255}, ${b * brightness / 255})`,
        },
    };

    return <div style={styles.colorBox} />;
}


export const getServerSideProps: GetServerSideProps = async () => {
    const client = getServerConnection();
    const AuthObject = createAuth();
    const resp = await client.getLightState({ auth: AuthObject });

    if (resp.state === undefined)
        resp.state = {} as LightState;

    if (resp.state.brightness === undefined)
        resp.state.brightness = 0;

    if (resp.state.on === undefined)
        resp.state.on = false;

    if (resp.state.colored === undefined)
        resp.state.colored = { red: 0, green: 0, blue: 0 };

    if (resp.state.warmwhite === undefined)
        resp.state.warmwhite = true;

    return {
        props: {
            title: "Light Control",
            lightState: resp.state,
        }
    };
};

export default function Content(props: { title: string, lightState: LightState }) {
    const state = props.lightState;

    const [power, setPower] = useState(state.on);
    const [brightness, setBrightness] = useState(state.brightness);
    const [warmwhite, setWarmwhite] = useState(state.warmwhite ?? false);

    if (state.colored === undefined)
        state.colored = { red: 0, green: 0, blue: 0 };

    const [color, setColor] = useState({ r: state.colored.red, g: state.colored.green, b: state.colored.blue } as RgbColor);

    function doUpdate(req: LightState) {
        fetch('/api/light', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ state: req }),
        }).then(res => {
            return res.json();
        }).then((data: ClientAPIResponse<UpdateLightAPIResponse>) => {
            if (data.success && data.data?.state) {
                const state = data.data.state;
                setPower(state.on);
                setBrightness(state.brightness);
                setWarmwhite(state.warmwhite ?? false);
                if (state.colored)
                    setColor({ r: state.colored.red, g: state.colored.green, b: state.colored.blue });
            }
        }).catch(err => {
            console.error(err);
        });
    }

    const handlePowerChange = () => {
        const state: LightState = { on: !power, brightness: brightness, colored: warmwhite ? undefined : { red: color.r, green: color.g, blue: color.b }, warmwhite: warmwhite };
        doUpdate(state);
    };

    const handleBrightnessChange = (event: Event, newValue: number | number[]) => {
        setBrightness(newValue as number);
    };

    const handleBrightnessSubmit = () => {
        const state: LightState = {
            on: power,
            brightness: brightness as number,
            colored: warmwhite ? undefined : { red: color.r, green: color.g, blue: color.b },
            warmwhite: warmwhite ? warmwhite : undefined
        };
        doUpdate(state);
    };

    const swapWarmWhite = () => {
        if (!warmwhite) {
            const state: LightState = { on: power, brightness: brightness, colored: undefined, warmwhite: true };
            doUpdate(state);
        } else {
            const state: LightState = { on: power, brightness: brightness, colored: { red: color.r, green: color.g, blue: color.b }, warmwhite: undefined };
            doUpdate(state);
        }
    };

    const handleColorChange = (color: ColorResult) => {
        setColor(color.rgb);
    };

    const handleColorChangeComplete = () => {
        const state: LightState = { on: power, brightness: brightness, colored: { red: color.r, green: color.g, blue: color.b }, warmwhite: undefined };
        doUpdate(state);
    };

    return (
        <Container>
            <Grid container sx={{ alignItems: "center" }}>
                <Grid item xs={3}>
                    Power
                </Grid>
                <Grid item xs={8}>
                    <Switch checked={power} onChange={handlePowerChange} name="Power" inputProps={{ 'aria-label': 'secondary checkbox' }} />
                </Grid>
                <Grid item xs={1}>
                    <ColorBox powered={power} warmwhite={warmwhite} color={color} brightness={brightness} />
                </Grid>

                <Grid item xs={3}>
                    <Typography id="discrete-slider" gutterBottom>Brightness</Typography>
                </Grid>
                <Grid item xs={9} columns={2}>
                    <Slider
                        value={brightness}
                        onChange={handleBrightnessChange}
                        onChangeCommitted={handleBrightnessSubmit}
                        min={0}
                        max={255}
                        disabled={!power}
                    />
                </Grid>

                <Grid item xs={3}>
                    <Typography id="discrete-slider" gutterBottom>Warm White</Typography>
                </Grid>
                <Grid item xs={9}>
                    <Checkbox
                        checked={warmwhite}
                        onChange={swapWarmWhite}
                        name="Warm White"
                        inputProps={{ 'aria-label': 'secondary checkbox' }}
                        disabled={!power}
                    />
                </Grid>

                <Grid item xs={3}>
                    <Typography id="discrete-slider" gutterBottom>Color</Typography>
                </Grid>
                <Grid item xs={9}>
                    <HuePicker
                        color={color}
                        width="100%"
                        onChange={handleColorChange}
                        onChangeComplete={handleColorChangeComplete}
                        pointer={SliderPointer}
                    />
                </Grid>
            </Grid>
        </Container>
    );
}
