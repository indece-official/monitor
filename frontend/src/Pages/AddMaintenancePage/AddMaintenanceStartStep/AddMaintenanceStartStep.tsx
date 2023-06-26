import React from 'react';
import DayJS from 'dayjs';
import { Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { MaintenanceService } from '../../../Services/MaintenanceService';
import { HostService, HostV1 } from '../../../Services/HostService';
import { TagService, TagV1 } from '../../../Services/TagService';
import { InputSelect } from '../../../Components/Input/InputSelect';
import { InputTextarea } from '../../../Components/Input/InputTextarea';
import { InputDateTime } from '../../../Components/Input/InputDateTime';


export interface AddMaintenanceStartStepProps
{
    hostUID:    string;
    onFinish:   ( maintenanceUID: string ) => any;
}


interface AddMaintenanceStartStepFormData
{
    title:              string;
    message:            string;
    affected_host_uids: Array<string>;
    affected_tag_uids:  Array<string>;
    datetime_start:     string;
    datetime_finish:    string;
}


interface AddMaintenanceStartStepState
{
    initialFormData:    AddMaintenanceStartStepFormData;
    hosts:              Array<HostV1>;
    tags:               Array<TagV1>;
    loading:            boolean;
    error:              Error | null;
}


export class AddMaintenanceStartStep extends React.Component<AddMaintenanceStartStepProps, AddMaintenanceStartStepState>
{
    private readonly _hostService: HostService;
    private readonly _tagService: TagService;
    private readonly _maintenanceService: MaintenanceService;


    constructor ( props: AddMaintenanceStartStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                title:              '',
                message:            '',
                affected_host_uids: [],
                affected_tag_uids:  [],
                datetime_start:     DayJS().format(),
                datetime_finish:    ''
            },
            hosts:      [],
            tags:       [],
            loading:    false,
            error:      null
        };

        this._hostService = HostService.getInstance();
        this._tagService = TagService.getInstance();
        this._maintenanceService = MaintenanceService.getInstance();

        this._submit = this._submit.bind(this);
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const hosts = await this._hostService.getHosts();
            const tags = await this._tagService.getTags();

            this.setState({
                loading:    false,
                hosts,
                tags
            });
        }
        catch ( err )
        {
            console.error(`Error loading hosts and tags: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    private async _submit ( values: AddMaintenanceStartStepFormData ): Promise<void>
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

            const maintenanceUID = await this._maintenanceService.addMaintenance({
                title:              values.title.trim(),
                message:            values.message.trim(),
                affected: {
                    host_uids:          values.affected_host_uids || [],
                    check_uids:         [],
                    tag_uids:           values.affected_tag_uids || []
                },
                datetime_start:     DayJS(values.datetime_start.trim()).format(),
                datetime_finish:    values.datetime_finish.trim() ? DayJS(values.datetime_finish.trim()).format() : null
            });

            this.setState({
                loading:    false
            });

            this.props.onFinish(maintenanceUID);
        }
        catch ( err )
        {
            console.error(`Error adding maintenance: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();
    }


    public render ( )
    {
        return (
            <div className='AddMaintenanceStartStep'>
                <h1>Add a maintenance</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            <InputText
                                name='title'
                                label='Title'
                                required={true}
                            />
                            
                            <InputTextarea
                                name='message'
                                label='Message'
                            />

                            <InputSelect
                                name='affected_host_uids'
                                label='Affected hosts'
                                options={this.state.hosts.map( ( host ) => ({
                                    label:  host.name,
                                    value:  host.uid
                                }))}
                            />

                            <InputSelect
                                name='affected_tag_uids'
                                label='Affected tags'
                                options={this.state.tags.map( ( tag ) => ({
                                    label:  tag.name,
                                    value:  tag.uid
                                }))}
                            />

                            <InputDateTime
                                name='datetime_start'
                                label='Start at'
                                required={true}
                            />
                          
                            <InputDateTime
                                name='datetime_finish'
                                label='Finish at'
                            />

                            <Button
                                type='submit'
                                disabled={this.state.loading}>
                                Create
                            </Button>
                        </Form>
                    )}
                </Formik>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}
