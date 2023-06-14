import React from 'react';
import { FieldArray, Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';
import { HostService } from '../../../Services/HostService';
import { InputMultiSelect } from '../../../Components/Input/InputMultiSelect';
import { TagService, TagV1 } from '../../../Services/TagService';


export interface SetupAddHostsStepProps
{
    onFinish:   ( ) => any;
}


interface SetupAddHostsStepFormDataHost
{
    name:       string;
    tag_uids:   Array<string>;
}


interface SetupAddHostsStepFormData
{
    hosts: Array<SetupAddHostsStepFormDataHost>;
}


interface SetupAddHostsStepState
{
    initialFormData:    SetupAddHostsStepFormData;
    tags:               Array<TagV1>
    loading:            boolean;
    error:              Error | null;
}


export class SetupAddHostsStep extends React.Component<SetupAddHostsStepProps, SetupAddHostsStepState>
{
    private readonly _hostService: HostService;
    private readonly _tagService: TagService;


    constructor ( props: SetupAddHostsStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                hosts: [
                    {
                        name:           '',
                        tag_uids:   []
                    }
                ]
            },
            tags:       [],
            loading:    false,
            error:      null
        };

        this._hostService = HostService.getInstance();
        this._tagService = TagService.getInstance();

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

            const tags = await this._tagService.getTags();

            this.setState({
                loading:    false,
                tags
            });
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
    
    
    private async _submit ( values: SetupAddHostsStepFormData ): Promise<void>
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

            for ( const host of values.hosts )
            {
                await this._hostService.addHost({
                    name:       host.name.trim(),
                    tag_uids:   host.tag_uids.map( s => s.trim() )
                });
            }

            this.setState({
                loading:    false
            });

            this.props.onFinish();
        }
        catch ( err )
        {
            console.error(`Error adding host: ${(err as Error).message}`, err);

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
            <div className='SetupAddHostsStep'>
                <h1>Setup - Add hosts</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            <FieldArray
                                name='hosts'>
                                {( arrayHelpers ) => (
                                    <div className='SetupAddHostsStep-hosts'>
                                        <Button
                                            type='button'
                                            onClick={() => arrayHelpers.push({
                                                labels: []
                                            })}>
                                            <FontAwesomeIcon icon={faPlus} />
                                            Add a host
                                        </Button>

                                        <div className='SetupAddHostsStep-hosts-list'>
                                            {values.hosts.length === 0 ?
                                                <div className='SetupAddHostsStep-hosts-empty'>
                                                    No hosts added yet
                                                </div>
                                            : null}
                                        </div>

                                        {values.hosts.map( ( host, index ) => (
                                            <div key={index} className='AdminAddEventPage-host'>
                                                <div className='AdminAddEventPage-host-form'>
                                                    <InputText
                                                        name={`hosts.${index}.name`}
                                                        label='Name'
                                                        required={true}
                                                    />

                                                    <InputMultiSelect
                                                        name={`hosts.${index}.tag_uids`}
                                                        label='Tags'
                                                        options={this.state.tags.map( ( tag ) => ({
                                                            label:  tag.name,
                                                            value:  tag.uid
                                                        }))}
                                                    />
                                                </div>

                                                <Button
                                                    title='Delete host'
                                                    type='button'
                                                    onClick={() => arrayHelpers.remove(index)}>
                                                    <FontAwesomeIcon icon={faTimes} />
                                                </Button>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </FieldArray>

                            <Button
                                type='submit'
                                disabled={this.state.loading}>
                                Save
                            </Button>
                        </Form>
                    )}
                </Formik>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}
