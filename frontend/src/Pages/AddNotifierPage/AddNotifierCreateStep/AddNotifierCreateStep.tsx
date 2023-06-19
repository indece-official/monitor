import React from 'react';
import { AddNotifierV1Request, NotifierService } from '../../../Services/NotifierService';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';


export interface AddNotifierCreateStepProps
{
    addNotifier:    AddNotifierV1Request;
    onFinish:       ( notifierUID: string ) => any;
}


export interface AddNotifierCreateStepState
{
    loading:    boolean;
    error:      Error | null;
}


export class AddNotifierCreateStep extends React.Component<AddNotifierCreateStepProps, AddNotifierCreateStepState>
{
    private readonly _notifierService:  NotifierService;


    constructor ( props: AddNotifierCreateStepProps )
    {
        super(props);

        this.state = {
            loading:    true,
            error:      null
        };

        this._notifierService = NotifierService.getInstance();
    }


    private async _create ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const notifierUID = await this._notifierService.addNotifier(this.props.addNotifier);

            this.setState({
                loading:    false
            });

            this.props.onFinish(notifierUID);
        }
        catch ( err )
        {
            console.error(`Error loading tags: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._create();
    }


    public render ( )
    {
        return (
            <div className='AddNotifierCreateStep'>
                <h1>Add a notifier</h1>

                <Spinner active={this.state.loading} />

                <ErrorBox error={this.state.error} />
            </div>
        );
    }
}
