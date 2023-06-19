import React from 'react';
import { FieldArray, Form, Formik } from 'formik';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { NotifierV1Filter } from '../../../Services/NotifierService';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { InputMultiSelect } from '../../../Components/Input/InputMultiSelect';
import { InputCheckbox } from '../../../Components/Input/InputCheckbox';
import { TagService, TagV1 } from '../../../Services/TagService';


export interface AddNotifierFiltersStepProps
{
    onFinish:   ( filters: Array<NotifierV1Filter> ) => any;
}


interface AddNotifierFiltersStepFormDataFilter
{
    tag_uids:       Array<string>;
    critical:       boolean;
    warning:        boolean;
    unknown:        boolean;
    decline:        boolean;
    min_duration:   string;
}


interface AddNotifierFiltersStepFormData
{
    filters:     Array<AddNotifierFiltersStepFormDataFilter>;
}


interface AddNotifierFiltersStepState
{
    initialFormData:    AddNotifierFiltersStepFormData;
    tags:               Array<TagV1>
    loading:            boolean;
    error:              Error | null;
}


export class AddNotifierFiltersStep extends React.Component<AddNotifierFiltersStepProps, AddNotifierFiltersStepState>
{
    private readonly _tagService: TagService;


    constructor ( props: AddNotifierFiltersStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                filters: []
            },
            tags:       [],
            loading:    false,
            error:      null
        };

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


    private async _submit ( values: AddNotifierFiltersStepFormData ): Promise<void>
    {
        const filters: Array<NotifierV1Filter> = values.filters.map( ( filterValues ) => ({
            tag_uids:       filterValues.tag_uids.filter( s => !!s ),
            critical:       filterValues.critical,
            warning:        filterValues.warning,
            unknown:        filterValues.unknown,
            decline:        filterValues.decline,
            min_duration:   filterValues.min_duration.trim()
        }));

        this.props.onFinish(filters);
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();
    }


    public render ( )
    {
        return (
            <div className='AddNotifierFiltersStep'>
                <h1>Add a notifier</h1>

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            <FieldArray name='filters'>
                                {( arrayHelpers ) => (
                                    <div className='SetupAddHostsStep-hosts'>
                                        <Button
                                            type='button'
                                            onClick={() => arrayHelpers.push({
                                                tag_uids:       [],
                                                critical:       true,
                                                warning:        true,
                                                unknown:        true,
                                                decline:        true,
                                                min_duration:   '5m'
                                            })}>
                                            <FontAwesomeIcon icon={faPlus} />
                                            Add a filter
                                        </Button>

                                        <div className='SetupAddHostsStep-hosts-list'>
                                            {values.filters.length === 0 ?
                                                <div className='SetupAddHostsStep-hosts-empty'>
                                                    No filters added yet
                                                </div>
                                            : null}
                                        </div>

                                        {values.filters.map( ( filter, index ) => (
                                            <div key={index} className='AdminAddEventPage-host'>
                                                <InputMultiSelect
                                                    name={`filters.${index}.tag_uids`}
                                                    label='Tags'
                                                    options={this.state.tags.map( ( tag ) => ({
                                                        label:  tag.name,
                                                        value:  tag.uid
                                                    }))}
                                                />

                                                <InputCheckbox
                                                    name={`filters.${index}.critical`}
                                                    label='Notify on critical'
                                                />

                                                <InputCheckbox
                                                    name={`filters.${index}.warning`}
                                                    label='Notify on warning'
                                                />

                                                <InputCheckbox
                                                    name={`filters.${index}.unknown`}
                                                    label='Notify on unknown'
                                                />

                                                <InputCheckbox
                                                    name={`filters.${index}.decline`}
                                                    label='Notify on decline'
                                                />

                                                <InputText
                                                    name={`filters.${index}.min_duration`}
                                                    label='Min. duration'
                                                    required={true}
                                                />

                                                <Button
                                                    title='Delete filter'
                                                    type='button'
                                                    onClick={() => arrayHelpers.remove(index)}>
                                                    <FontAwesomeIcon icon={faTimes} />
                                                </Button>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </FieldArray>

                            <Button type='submit'>
                                Continue
                            </Button>
                        </Form>
                    )}
                </Formik>
            </div>
        );
    }
}
