import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';
import { ConfigService } from '../../../Services/ConfigService';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { Button } from '../../../Components/Button/Button';


export interface SetupSuccessStepProps
{
}


interface SetupSuccessStepState
{
    error:      Error | null;
    loading:    boolean;
}


export class SetupSuccessStep extends React.Component<SetupSuccessStepProps,SetupSuccessStepState>
{
    private readonly _configService: ConfigService;


    constructor ( props: SetupSuccessStepProps )
    {
        super(props);

        this.state = {
            loading:    false,
            error:      null
        };

        this._configService = ConfigService.getInstance();
    }


    private async _finish ( ): Promise<void>
    {
        try
        {
            if ( this.state.loading )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            await this._configService.finishSetup();

            this.setState({
                loading:    false
            });
        }
        catch ( err )
        {
            console.error(`Error finishing setup: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._finish();
    }


    public render ( )
    {
        return (
            <div className='SetupSuccessStep'>
                <h1>Setup abgeschlossen</h1>

                <ErrorBox error={this.state.error} />

                <SuccessBox
                    message={`Der Setup wurde erfolgreich angeschlossen`}
                />

                <Spinner active={this.state.loading} />

                <Button href='/login'>
                    Weiter zum Login
                </Button>
            </div>
        );
    }
}
