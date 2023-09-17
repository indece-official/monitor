import React from 'react';
import { SetupAddUserStep } from './SetupAddUserStep/SetupAddUserStep';
import { SetupSuccessStep } from './SetupSuccessStep/SetupSuccessStep';
import { SetupAddHostsStep } from './SetupAddHosts/SetupAddHosts';
import { SetupServerInfoStep } from './SetupServerInfoStep/SetupServerInfoStep';
import { SetupAddTagsStep } from './SetupAddTags/SetupAddTags';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface SetupPageProps
{
}


enum SetupStep
{
    AddUser     = 'add_user',
    ServerInfo  = 'server_info',
    AddTags     = 'add_tags',
    AddHosts    = 'add_hosts',
    Success     = 'success'
}


interface SetupPageState
{
    step:       SetupStep;
    userUID:    string | null;
}


export class SetupPage extends React.Component<SetupPageProps, SetupPageState>
{
    constructor ( props: SetupPageProps )
    {
        super(props);

        this.state = {
            step:       SetupStep.AddUser,
            userUID:    null
        };

        this._finishAddUser = this._finishAddUser.bind(this);
        this._finishServerInfo = this._finishServerInfo.bind(this);
        this._finishAddTags = this._finishAddTags.bind(this);
        this._finishAddHosts = this._finishAddHosts.bind(this);
    }


    private _finishAddUser ( userUID: string ): void
    {
        this.setState({
            step:   SetupStep.ServerInfo,
            userUID
        });
    }
   
   
    private _finishServerInfo ( ): void
    {
        this.setState({
            step:   SetupStep.AddTags
        });
    }


    private _finishAddTags ( ): void
    {
        this.setState({
            step:   SetupStep.AddHosts
        });
    }
   
   
    private _finishAddHosts ( ): void
    {
        this.setState({
            step:   SetupStep.Success
        });
    }


    public render ( )
    {
        return (
            <PageContent>
                {this.state.step === SetupStep.AddUser ?
                    <SetupAddUserStep
                        onFinish={this._finishAddUser}
                    />
                : null}
               
                {this.state.step === SetupStep.ServerInfo ?
                    <SetupServerInfoStep
                        onFinish={this._finishServerInfo}
                    />
                : null}
                
                {this.state.step === SetupStep.AddTags ?
                    <SetupAddTagsStep
                        onFinish={this._finishAddTags}
                    />
                : null}
                
                {this.state.step === SetupStep.AddHosts ?
                    <SetupAddHostsStep
                        onFinish={this._finishAddHosts}
                    />
                : null}
                
                {this.state.step === SetupStep.Success ?
                    <SetupSuccessStep />
                : null}
            </PageContent>
        );
    }
}
